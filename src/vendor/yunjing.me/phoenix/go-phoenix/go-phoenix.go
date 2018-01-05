package phoenix

import (
	"errors"
	"net/http"
	"time"

	oemitter "yunjing.me/phoenix/go-phoenix/emitter"
	ostorage "yunjing.me/phoenix/go-phoenix/storage"
	oweb "yunjing.me/phoenix/go-phoenix/web"
)

var (
	ErrSlabMemPoolNotEnable = errors.New("Phoenix: slab mempool not enable")
)

type IListener interface {
	OnHeartBeat()
}

type Service interface {
	// -----------------------------------------------------------------------

	// 服务类型
	Category() uint8

	// 服务编号
	ID() uint8

	// 生命周期：初始化
	Init(...Option)

	// 生命周期：运行
	Serve(IListener)

	CheckUserSecret(uid uint32, secret []byte) bool

	// 生成访问令牌
	CreateUserSecret(uid uint32) []byte

	DBManager() *ostorage.DBManager
	CacheManager() *ostorage.CacheManager

	// -----------------------------------------------------------------------

	Web() oweb.Service
	RegisterWebHandle(pattern string, handler http.Handler)
	RegisterWebHandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))

	// -----------------------------------------------------------------------

	// 事件模块
	EventEmitter() *oemitter.Emitter

	// -----------------------------------------------------------------------

	// 注册定时器
	AddTimer(holder interface{}, delay int, interval int, callback func()) uint32

	// 移除定时器
	RemoveTimer(uint32)

	// -----------------------------------------------------------------------

	// 数据收集
	Collect(happen time.Time, id uint32, args ...interface{})

	// -----------------------------------------------------------------------
}

type Option func(*Options)

func NewService(opts ...Option) Service {
	return newService(opts...)
}
