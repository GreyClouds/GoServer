package skeleton

import (
	"log"
	"net"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/golang/protobuf/proto"
	"yunjing.me/phoenix/cli"
	"yunjing.me/phoenix/go-phoenix"
	odisconf "yunjing.me/phoenix/go-phoenix/disconf"
	ostorage "yunjing.me/phoenix/go-phoenix/storage"
	oweb "yunjing.me/phoenix/go-phoenix/web"

	oaccount "webapi/account"
	//pkgActivity "webapi/activity"
	obean "webapi/bean"
	pkgGMHandler "webapi/cmds/gm"
	pkgOtherHandler "webapi/cmds/other"
	//. "webapi/common"
	oconf "webapi/config"
	FYSDK "webapi/fysdk"
	osession "webapi/session"
	pkgTrace "webapi/trace"
	//"webapi/wordstock"
	"webapi/wordstock"
)

type Skeleton struct {
	qps             int32 // 全部的QPS
	loginQPS        int32 // 登陆QPS
	mux             sync.Mutex
	pService        phoenix.Service
	audits          map[uint16]bool
	prototypes      map[uint16]proto.Message
	handlers        map[uint16]MessageHandler
	pAccountManager *oaccount.Manager
	qpsNum          *int64
	closeFlag       int32
}

func New() *Skeleton {
	var i int64 = 0
	self := &Skeleton{
		audits:          make(map[uint16]bool),
		prototypes:      make(map[uint16]proto.Message),
		handlers:        make(map[uint16]MessageHandler),
		pAccountManager: oaccount.New(),
		qpsNum:          &i,
	}

	self.pService = phoenix.NewService(
		phoenix.Name("deadfat.webapi"),
		phoenix.Flags(
			cli.StringFlag{
				Name:  "static_data",
				Value: "./data",
				Usage: "静态配置文件存放目录",
			},
			cli.StringFlag{
				Name:  "dirty_word",
				Value: "./words.txt",
				Usage: "敏感词过滤文件存放路径",
			},
		),
	)

	self.pService.Init(
		phoenix.EnableDB(),
		phoenix.EnableCache(),
		phoenix.EnableMemPool(),
		phoenix.EnableWeb(),
		phoenix.EnableCollector(),
		phoenix.Action(func(c *cli.Context) {
			oconf.Singleton().Preload(odisconf.New(c.String("static_data")))
			wordstock.Configurate(c.String("dirty_word"))
		}),
		phoenix.BeforeStart(self.OnPreload),
		phoenix.AfterStop(self.OnServerClose),
	)

	return self
}

func (self *Skeleton) Register(audit bool, id uint16, prototype proto.Message, handler MessageHandler) {
	if audit {
		self.audits[id] = true
	}

	if prototype != nil {
		self.prototypes[id] = prototype
	}

	if handler != nil {
		self.handlers[id] = handler
	}
}

func (self *Skeleton) doBeanPreload() {
	//obean.PreloadAccount()
	//obean.PreloadAnnouncement()

	obean.DefaultORM()
	//obean.DefaultKeeper().Run()
	//oaccount.DefaultRankList().Load()
	obean.ArenaLeaBoardLoadAndSort();
	obean.ChallLeaBoardLoadAndSort();

}

func (self *Skeleton) OnPreload() error {
	self.doBeanPreload()

	//pkgActivity.DefaultActivityManager.Preload(oconf.Singleton().GetActivityList())

	// 处理真实秒心跳
	//go func() {
	//	t := time.NewTicker(time.Second)
	//	for {
	//		select {
	//		case <-t.C:
	//			self.OnSecondTick()
	//		}
	//	}
	//}()

	// 处理空闲协程
	//go func() {
	//	t1 := ZeroAclock().Add(3 * time.Hour)
	//	time.AfterFunc(t1.Sub(time.Now()), func() {
	//		for {
	//			self.OnIdle()
	//			time.Sleep(24 * time.Hour)
	//		}
	//	})
	//}()

	return nil
}

func (self *Skeleton) RegisterDBModel() *Skeleton {
	orm.RegisterModel(&obean.AchievementAttr{},&obean.AchievementUnLock{},&obean.ArenaLeaderboard{},
		&obean.ChallLeaderboard{}, &obean.SimAccount{}, &obean.SimRole{},
	)

	return self
}

func (self *Skeleton) Serve() {
	self.pService.RegisterWebHandleFunc("/fysdk/pay", self.HandleFYSDKPay)
	self.pService.RegisterWebHandleFunc("/o", self.HandleAPI)
	self.pService.RegisterWebHandleFunc("/gm", self.HandleGM)
	self.pService.RegisterWebHandleFunc("/homepage", self.HandleHomePage)
	self.pService.Serve(self)
}


// ---------------------------------------------------------------------------

func (self *Skeleton) CacheManager() *ostorage.CacheManager {
	return self.pService.CacheManager()
}

func (self *Skeleton) AccountManager() *oaccount.Manager {
	return self.pAccountManager
}

func (self *Skeleton) DBManager() *ostorage.DBManager {
	return self.pService.DBManager()
}

func (self *Skeleton) ConfigManager() *oconf.Manager {
	return oconf.Singleton()
}

func (self *Skeleton) Web() oweb.Service {
	return self.pService.Web()
}

func (self *Skeleton) Collect(id uint32, args ...interface{}) {
	self.pService.Collect(time.Now(), id, args...)
}

// ---------------------------------------------------------------------------

func (self *Skeleton) OnHeartBeat() {
	self.pAccountManager.OnHeartBeat()
}

// 真实环境的秒心跳
func (self *Skeleton) OnSecondTick() {
	self.pService.Collect(time.Now(),
		pkgTrace.ProcessStatusReport,
		self.pAccountManager.GetOnlineNum(),
		atomic.SwapInt32(&self.qps, 0),
		self.pAccountManager.GetLoginNum(),
	)
}

// 空闲处理
func (self *Skeleton) OnIdle() {
	start := time.Now()

	obean.BattleRoomCleanTask()
	obean.BattleResultCleanTask()

	log.Printf("空闲处理耗时: %v", time.Now().Sub(start))
}

func (self *Skeleton) isClosed() bool {
	return atomic.LoadInt32(&self.closeFlag) == 1
}

// 响应服务关闭
func (self *Skeleton) OnServerClose() error {
	log.Println("服务关闭准备...")

	if atomic.CompareAndSwapInt32(&self.closeFlag, 0, 1) {
		//self.pAccountManager.BGSave()
		//obean.DefaultKeeper().Close()
	}

	log.Println("服务准备优雅关闭")

	return nil
}

// ---------------------------------------------------------------------------

func (self *Skeleton) RemoteIP(req *http.Request) string {
	forward := req.Header.Get("X-Forwarded-For")
	if forward != "" {
		return forward
	}

	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return ""
	}

	userIP := net.ParseIP(ip)
	if userIP == nil {
		return ""
	}

	return ip
}

// 处理API逻辑
func (self *Skeleton) HandleAPI(w http.ResponseWriter, r *http.Request) {
	if self.isClosed() {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	// 解析消息包
	err, packet := recv(r)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	if ok := packet.CheckSign(); !ok {
		log.Printf("消息包签名校验失败: %s", r.RemoteAddr)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	var session *osession.Session
	var roleob *oaccount.Role

	if _, exists := self.audits[packet.pid]; exists {
		if packet.uid == 0 {
			log.Printf("收到协议%d时请求未授权", packet.pid)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		if ok := self.pService.CheckUserSecret(packet.uid, packet.token); !ok {
			log.Printf("收到协议%d时请求认证失败: %d %s", packet.pid, packet.uid, packet.token)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		session = osession.NewSession(packet.rn, packet.uid, self.RemoteIP(r))

		roleob = self.pAccountManager.LoadRole(packet.uid)
		if roleob == nil {
			log.Printf("角色对象为空: %v", packet.uid)
			http.Error(w, "", http.StatusBadRequest)
			return
		}
	} else {
		session = osession.NewSession(packet.rn, packet.uid, self.RemoteIP(r))
	}

	var instance proto.Message
	if prototype := self.prototypes[packet.pid]; prototype != nil {
		instance = proto.Clone(prototype)
		if err = proto.Unmarshal(packet.payload, instance); err != nil {
			doServerInternalErrorSend(w, err)
			return
		}
	}

	if handler := self.handlers[packet.pid]; handler != nil {
		err, id, packet := self.doMSGHandle(session, roleob, instance, handler)

		// 累计QPS数据
		atomic.AddInt32(&self.qps, 1)

		if err != nil {
			doServerInternalErrorSend(w, err)
			return
		}

		send(w, id, packet)
	} else {
		log.Printf("协议%d未处理", packet.pid)
		doServerInternalErrorSend(w, nil)
		return
	}
}

func (self *Skeleton) HandleGM(w http.ResponseWriter, r *http.Request) {

	pkgGMHandler.Handle(self, w, r)
}

func (self *Skeleton) HandleFYSDKPay(w http.ResponseWriter, r *http.Request) {

	FYSDK.PayNotify(self, w, r)
}

func (self *Skeleton) HandleHomePage(w http.ResponseWriter, r *http.Request) {

	pkgOtherHandler.HandleHomePage(self, w, r)
}

func (self *Skeleton) doMSGHandle(session *osession.Session, roleob *oaccount.Role, instance proto.Message, callback MessageHandler) (error, uint16, proto.Message) {
	if roleob != nil {
		roleob.Lock()
		defer func() {
			roleob.UpdateLastActionTime()
			roleob.Unlock()
		}()
	}

	return callback(self, session, roleob, instance)
}

func (self *Skeleton) CreateUserSecret(uid uint32) []byte {
	return self.pService.CreateUserSecret(uid)
}

func (self *Skeleton) OnPaid(bean *obean.AndroidPayment) {
	// log.Printf("处理FYSDK支付订单: %v", bean)

	uid := bean.UID
	roleob := self.pAccountManager.GetRole(uid)
	if roleob == nil {
		return
	}

	if err := bean.UpdateAsAchieved(); err != nil {
		log.Printf("标记订单为已领取状态时出错: %v", err)
		return
	}

	roleob.Lock()
	defer roleob.Unlock()

	self.pAccountManager.OnPaid(roleob, bean)
}
