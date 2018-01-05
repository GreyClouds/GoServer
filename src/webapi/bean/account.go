package bean

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/astaxie/beego/orm"

	. "webapi/common"
)

type Account struct {
	Channel   string    `orm:"column(channel);size(6)"`  // 渠道编号
	Way       uint32    `orm:"column(way)"`              // 登录方式[1:设备号,2:FYSDK联网,3:FYSDK单机]
	Version   string    `orm:"column(version);size(32)"` // 客户端版本
	IMEI      string    `orm:"column(imei);size(255)"`   // 机器设备号
	UserID    string    `orm:"column(userid);size(255)"` // 账号标识
	Uid       uint32    `orm:"pk;column(uid)"`           // 角色编号
	CreatedAt time.Time `orm:"column(created_at)"`       // 注册时间
}

func (self *Account) TableName() string {
	return "account"
}

// ---------------------------------------------------------------------------

var (
	pIMEIAccountCache         sync.Map
	pFYSDKOnlineAccountCache  sync.Map
	pFYSDKOfflineAccountCache sync.Map
)

func cacheAccount(way uint32, userid string, account *Account) {
	switch way {
	case GUEST_LOGIN: // 游客登录
		pIMEIAccountCache.Store(userid, account)

	case FYSDK_ONLINE: // FYSDK联网版本
		pFYSDKOnlineAccountCache.Store(userid, account)

	case FYSDK_OFFLINE: // FYSDK单机版本
		pFYSDKOfflineAccountCache.Store(userid, account)

	default:

	}
}

// 预载入所有的账号信息
func PreloadAccount() {
	o := orm.NewOrm()

	var accounts []*Account
	num, err := o.QueryTable(new(Account)).Limit(-1).All(&accounts)
	if err != nil {
		panic(fmt.Sprintf("载入账号信息缓存数据时出错: %v", err))
	}

	for i := int64(0); i < num; i++ {
		item := accounts[i]
		cacheAccount(item.Way, item.UserID, item)
	}

	log.Printf("预载入%d条账号记录", num)
}

func FindAccount(way uint32, userid string) (error, *Account) {
	switch way {
	case GUEST_LOGIN:
		account, ok := pIMEIAccountCache.Load(userid)
		if ok {
			return nil, account.(*Account)
		}
	case FYSDK_ONLINE: // FYSDK联网版本
		account, ok := pFYSDKOnlineAccountCache.Load(userid)
		if ok {
			return nil, account.(*Account)
		}
	case FYSDK_OFFLINE: // FYSDK单机版本
		account, ok := pFYSDKOfflineAccountCache.Load(userid)
		if ok {
			return nil, account.(*Account)
		}
	default:
		return errors.New("unknown login way"), nil
	}

	return nil, nil
}

// 保存账号信息
func InsertAccount(channel string, way uint32, userid string, uid uint32, imei, version string) error {
	account := Account{
		Channel:   channel,
		Way:       way,
		UserID:    userid,
		Uid:       uid,
		IMEI:      imei,
		Version:   version,
		CreatedAt: time.Now(),
	}

	o := orm.NewOrm()
	_, err := o.Insert(&account)
	if err != nil {
		log.Printf("新增账号记录%d时出错: %v", uid, err)
		return err
	}

	cacheAccount(account.Way, account.UserID, &account)

	return nil
}
