package phoenix

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/funny/slab"

	"yunjing.me/phoenix/cli"
	ocollector "yunjing.me/phoenix/go-phoenix/collector"
	oemitter "yunjing.me/phoenix/go-phoenix/emitter"
	ostorage "yunjing.me/phoenix/go-phoenix/storage"
	oweb "yunjing.me/phoenix/go-phoenix/web"
)

const (
	kMaxPointChanSize   = 52 * 1024   // 统计数据队列最大尺寸
	kMaxNetMSGQueueSize = 1024 * 1024 // 网络消息队列最大尺寸
	kMaxEventQueueSize  = 1024        // 内部事件队列最大尺寸
)

type service struct {
	opts Options
	app  *cli.App

	pCacheManager *ostorage.CacheManager
	pDBManager    *ostorage.DBManager
	pool          slab.Pool
	web           oweb.Service
	pCollector2   *ocollector.Collector
	events        chan *oemitter.Event
	pEventEmitter *oemitter.Emitter
	listener      IListener

	oid    int32
	timers map[uint32]*internalTimer
}

func init() {
	rand.Seed(time.Now().UnixNano())
	help := cli.HelpPrinter
	cli.HelpPrinter = func(writer io.Writer, tpl string, data interface{}) {
		help(writer, tpl, data)
		os.Exit(0)
	}
}

func newService(opts ...Option) Service {
	self := &service{
		opts:   newOptions(opts...),
		app:    cli.NewApp(),
		events: make(chan *oemitter.Event, kMaxEventQueueSize),
		timers: make(map[uint32]*internalTimer),
	}

	self.pEventEmitter = oemitter.NewEventEmitter(self.events)

	return self
}

func (self *service) GenOID() uint32 {
	self.oid++
	return uint32(self.oid)
}

func (self *service) doFlagsInit() []cli.Flag {
	flags := []cli.Flag{
		cli.IntFlag{
			Name:   "id",
			EnvVar: "PHOENIX_ID",
			Value:  1,
			Usage:  "服务结点编号",
		},
	}

	if self.opts.isCacheEnable {
		flags = append(flags, []cli.Flag{
			cli.StringSliceFlag{
				Name:   "cache_addrs",
				EnvVar: "PHOENIX_CACHE_ADDRS",
				Usage:  "数据缓存连接地址",
			},
			cli.IntFlag{
				Name:   "cache_idx",
				EnvVar: "PHOENIX_CACHE_IDX",
				Value:  0,
				Usage:  "数据缓存连接编号",
			},
			cli.StringFlag{
				Name:   "cache_passwd",
				EnvVar: "PHOENIX_CACHE_PASSWD",
				Value:  "",
				Usage:  "数据缓存连接密码",
			},
		}...)
	}

	if self.opts.isDBEnable {
		flags = append(flags, []cli.Flag{
			cli.BoolFlag{
				Name:   "db_uri_encrypt",
				EnvVar: "PHOENIX_DB_URI_ENCRYPT",
				Usage:  "数据存储连接串是否加密",
			},
			cli.StringFlag{
				Name:   "db_uri",
				EnvVar: "PHOENIX_DB_URI",
				Value:  "",
				Usage:  "数据存储连接串",
			},
		}...)
	}

	if self.opts.isWebEnable {
		flags = append(flags, []cli.Flag{
			cli.StringFlag{
				Name:   "web_addr",
				EnvVar: "PHOENIX_WEB_ADDR",
				Value:  ":0",
				Usage:  "Web服务监听地址",
			},
		}...)
	}

	if self.opts.isCollectorEnable {
		flags = append(flags, []cli.Flag{
			cli.StringFlag{
				Name:   "collector_addr",
				EnvVar: "PHOENIX_COLLECTOR_ADDR",
				Value:  ":12178",
				Usage:  "数据收集服务监听地址",
			},
		}...)
	}

	return flags
}

func (self *service) Init(opts ...Option) {
	for _, o := range opts {
		o(&self.opts)
	}

	self.app.Name = self.opts.name
	self.app.Version = self.opts.version
	self.app.Usage = self.opts.description
	self.app.Flags = self.doFlagsInit()
	self.app.Action = self.opts.action
	self.app.Before = self.onLoaded

	if len(self.opts.flags) > 0 {
		self.app.Flags = append(self.app.Flags, self.opts.flags...)
	}

	if len(self.opts.version) == 0 {
		self.app.HideVersion = true
	}

	self.app.RunAndExitOnError()
}

// 响应初始化完成
func (self *service) onLoaded(ctx *cli.Context) error {
	self.opts.id = uint8(ctx.Int("id"))

	opts := self.opts

	// 是否启用数据缓存
	if enable := opts.isCacheEnable; enable {
		addrs := ctx.StringSlice("cache_addrs")
		idx := ctx.Int("cache_idx")
		passwd := ctx.String("cache_passwd")
		self.pCacheManager = ostorage.NewCacheManager(addrs, idx, passwd)
	}

	// 是否启用数据存储
	if enable := opts.isDBEnable; enable {
		encrypted := ctx.Bool("db_uri_encrypt")
		uri := ctx.String("db_uri")

		self.pDBManager = ostorage.NewDBManager(encrypted, uri)

		if opts.isDBSQLShow {
			self.pDBManager.EnableLog()
		}
	}

	// 是否启用内存池
	if enable := opts.isMemPoolEnable; enable {
		self.pool = newMemPool(opts.memPoolType, opts.memPoolFactor, opts.memPoolMinChunk, opts.memPoolMaxChunk, opts.memPoolPageSize)
	}

	// 是否启用Web服务
	if enable := opts.isWebEnable; enable {
		self.web = oweb.NewService(
			oweb.Address(ctx.String("web_addr")),
		)
	}

	// 是否启用数据采集服务
	if enable := opts.isCollectorEnable; enable {
		self.pCollector2 = ocollector.New(
			ctx.String("collector_addr"),
		)
	}

	return nil
}

func (self *service) onStart() error {
	// 处理DB启动
	if pDB := self.pDBManager; pDB != nil {
		if err := pDB.Run(); err != nil {
			return err
		}
	}

	for _, fn := range self.opts.BeforeStart {
		if err := fn(); err != nil {
			return err
		}
	}

	// 处理Web启动
	if svc := self.web; svc != nil {
		if err := svc.Start(); err != nil {
			return err
		}
	}

	for _, fn := range self.opts.AfterStart {
		if err := fn(); err != nil {
			return err
		}
	}

	return nil
}

func (self *service) onStop() error {
	var err error

	for _, fn := range self.opts.BeforeStop {
		if e := fn(); e != nil {
			err = e
		}
	}

	// Web关闭
	if svc := self.web; svc != nil {
		if e := svc.Stop(); e != nil {
			return e
		}
	}

	for _, fn := range self.opts.AfterStop {
		if e := fn(); e != nil {
			err = e
		}
	}

	return err
}

func (self *service) Serve(listener IListener) {
	self.listener = listener

	if err := self.onStart(); err != nil {
		log.Printf("[Phoenix] 服务启动出错: %v", err)
		return
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)

	timer := time.NewTicker(time.Second)

	for {
		select {
		case <-ch: // 外部信号量
			if err := self.onStop(); err != nil {
				log.Printf("[Phoenix] 服务停止出错: %v", err)
			}

			return
		case <-timer.C: // 定时心跳
			self.doHeartBeat()
		case evt := <-self.events: // 外部定义事件
			self.pEventEmitter.EmitSync(evt.ID, evt.Args)
		}
	}
}

func newMemPool(category string, factor, minChunk, maxChunk, pageSize int) slab.Pool {
	switch category {
	case "sync":
		return slab.NewSyncPool(minChunk, maxChunk, factor)
	case "atom":
		return slab.NewAtomPool(minChunk, maxChunk, factor, pageSize)
	case "chan":
		return slab.NewChanPool(minChunk, maxChunk, factor, pageSize)
	default:
		return nil
	}
}

// ---------------------------------------------------------------------------

func (self *service) DBManager() *ostorage.DBManager {
	return self.pDBManager
}

func (self *service) CacheManager() *ostorage.CacheManager {
	return self.pCacheManager
}

// ---------------------------------------------------------------------------

func (self *service) Web() oweb.Service {
	return self.web
}

func (self *service) RegisterWebHandle(pattern string, handler http.Handler) {
	if self.web == nil {
		log.Printf("[Phoenix]注册消息处理时失败: Web组件未初始化")
		return
	}

	self.web.Handle(pattern, handler)
}

func (self *service) RegisterWebHandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	if self.web == nil {
		log.Printf("[Phoenix]注册消息处理时失败: Web组件未初始化")
		return
	}

	self.web.HandleFunc(pattern, handler)
}

func (self *service) doHeartBeat() {
	// 处理定时器逻辑
	self.doTimerTick()

	// 处理外部心跳逻辑
	self.listener.OnHeartBeat()
}

// ---------------------------------------------------------------------------

// 事件模块
func (self *service) EventEmitter() *oemitter.Emitter {
	return self.pEventEmitter
}

// ---------------------------------------------------------------------------

// 数据收集
func (self *service) Collect(happen time.Time, id uint32, args ...interface{}) {
	if self.pCollector2 == nil {
		log.Printf("Collector组件未启用: %d", id)
		return
	}

	// self.ClusterID()

	if args != nil && len(args) > 0 {
		self.pCollector2.Do(happen, id, args...)
	} else {
		self.pCollector2.Do(happen, id)
	}
}

// ---------------------------------------------------------------------------

func (self *service) AddTimer(holder interface{}, delay, interval int, callback func()) uint32 {
	if delay <= 0 && interval <= 0 {
		log.Printf("注册定时器时出错: delay - %d, interval - %d", delay, interval)
		return 0
	}

	nDelay, nInterval := delay, interval
	if delay < 0 {
		nDelay = 0
	}
	if interval < 0 {
		nInterval = 0
	}

	oid := self.GenOID()
	timer := newTimer(oid, holder, nDelay, nInterval, callback)
	self.timers[oid] = timer
	return oid
}

func (self *service) RemoveTimer(id uint32) {
	delete(self.timers, id)
}

func (self *service) doTimerTick() {
	deletes := []uint32{}

	for id, v := range self.timers {
		v.OnTick()

		if end := v.IsEnd(); end {
			deletes = append(deletes, id)
		}
	}

	for _, v := range deletes {
		delete(self.timers, v)
	}
}

// ---------------------------------------------------------------------------

// 服务类型
func (self *service) Category() uint8 {
	return self.opts.category
}

// 服务编号
func (self *service) ID() uint8 {
	return self.opts.id
}

func (self *service) CreateUserSecret(uid uint32) []byte {
	raw := make([]byte, 4)
	binary.LittleEndian.PutUint32(raw, uid)

	var buf [md5.Size + 8]byte
	rand.Read(buf[md5.Size : md5.Size+8])

	hash := md5.New()
	hash.Write(buf[md5.Size : md5.Size+8])
	hash.Write(self.opts.userAuthSeed)
	hash.Write(raw)
	verify := hash.Sum(nil)
	copy(buf[:md5.Size], verify)

	return buf[:]
}

func (self *service) CheckUserSecret(uid uint32, token []byte) bool {
	if len(token) != (md5.Size + 8) {
		return false
	}

	raw := make([]byte, 4)
	binary.LittleEndian.PutUint32(raw, uid)

	hash := md5.New()
	hash.Write(token[md5.Size:])
	hash.Write(self.opts.userAuthSeed)
	hash.Write(raw)
	verify := hash.Sum(nil)

	return bytes.Equal(verify, token[:md5.Size])
}

// ---------------------------------------------------------------------------
