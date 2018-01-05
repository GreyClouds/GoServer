package bean

import (
	"time"
)

type SimAccount struct {
	Account    string    `orm:"pk;column(userid);size(255)"` // 账号标识
	Uid       uint32    `orm:"column(uid)"`           // 角色编号
	CreatedAt time.Time `orm:"column(created_at)"`       // 注册时间
}

type SimRole struct {
	Name 	   string 	`orm:"column(name);size(255)"`	  //角色名称
	Uid       uint32    `orm:"pk;column(uid)"`           // 角色编号
}

func (self *SimAccount) TableName() string {
	return "simAccount"
}



func FindSimAccount(acc string) (error, *SimAccount) {
	account := SimAccount{Account:acc}
	err := defaultOrm.Read(&account)
	return  err, &account
}

func RegisterAccount(acc string, name string, uid uint32) error {
	account := SimAccount{acc, uid, time.Now()}
	role := SimRole{name, uid}
	_, err := defaultOrm.Insert(&account)
	defaultOrm.Insert(&role)
	return err
}

func SetNickName(name string, uid uint32) error {
	role := SimRole{name, uid}
	_, err := defaultOrm.InsertOrUpdate(&role)
	return err
}