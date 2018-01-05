package phoenix

import (
	"yunjing.me/phoenix/cli"
)

type Options struct {
	name        string
	description string
	version     string
	flags       []cli.Flag
	action      func(*cli.Context)
	category    uint8
	id          uint8

	// Before and After funcs
	BeforeStart []func() error
	BeforeStop  []func() error
	AfterStart  []func() error
	AfterStop   []func() error

	// 其它
	userAuthSeed []byte

	// 内存池
	isMemPoolEnable bool
	memPoolType     string
	memPoolFactor   int
	memPoolMinChunk int
	memPoolMaxChunk int
	memPoolPageSize int

	// 数据缓存
	isCacheEnable bool

	// 数据持久
	isDBEnable  bool
	isDBSQLShow bool

	// 数据包
	maxPacketSize int

	// Web服务
	isWebEnable bool

	// 是否启用数据采集服务
	isCollectorEnable bool
}

func newOptions(opts ...Option) Options {
	opt := Options{
		name:        "go-phoenix",
		description: "Phoenix",
		version:     "0.0.0",
		flags:       []cli.Flag{},
		action:      func(c *cli.Context) {},
		category:    0,
		id:          0,

		maxPacketSize: 2 * 1024 * 1024,

		memPoolType:     "atom",
		memPoolFactor:   2,
		memPoolMinChunk: 64,
		memPoolMaxChunk: 64 * 1024,
		memPoolPageSize: 1024 * 1024,
	}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

// Name of the service
func Name(n string) Option {
	return func(o *Options) {
		o.name = n
	}
}

func Version(v string) Option {
	return func(o *Options) {
		o.version = v
	}
}

func Category(v uint8) Option {
	return func(o *Options) {
		o.category = v
	}
}

func Flags(flags ...cli.Flag) Option {
	return func(o *Options) {
		o.flags = append(o.flags, flags...)
	}
}

func Action(a func(*cli.Context)) Option {
	return func(o *Options) {
		o.action = a
	}
}

func BeforeStart(fn func() error) Option {
	return func(o *Options) {
		o.BeforeStart = append(o.BeforeStart, fn)
	}
}

func BeforeStop(fn func() error) Option {
	return func(o *Options) {
		o.BeforeStop = append(o.BeforeStop, fn)
	}
}

func AfterStart(fn func() error) Option {
	return func(o *Options) {
		o.AfterStart = append(o.AfterStart, fn)
	}
}

func AfterStop(fn func() error) Option {
	return func(o *Options) {
		o.AfterStop = append(o.AfterStop, fn)
	}
}

func UserAuthSeed(seed string) Option {
	return func(o *Options) {
		o.userAuthSeed = []byte(seed)
	}
}

// 启用数据缓存
func EnableCache() Option {
	return func(o *Options) {
		o.isCacheEnable = true
	}
}

// 启用数据落地
func EnableDB() Option {
	return func(o *Options) {
		o.isDBEnable = true
	}
}

// 启用数据存储SQL语句调试
func EnableDBLog() Option {
	return func(o *Options) {
		o.isDBSQLShow = true
	}
}

// 启用内存池模块
func EnableMemPool() Option {
	return func(o *Options) {
		o.isMemPoolEnable = true
	}
}

// 启用HTTP Web服务
func EnableWeb() Option {
	return func(o *Options) {
		o.isWebEnable = true
	}
}

// 启用数据收集
func EnableCollector() Option {
	return func(o *Options) {
		o.isCollectorEnable = true
	}
}
