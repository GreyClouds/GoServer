package storage

import (
	// "log"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

type DBManager struct {
	isVerbose bool
}

func NewDBManager(encrypted bool, uri string) *DBManager {
	var str string
	if encrypted {
		str = Decrypt(uri)
	} else {
		str = uri
	}

	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", str, 5, 5)

	// client, err := gorm.Open("mysql", str)
	// if err != nil {
	// 	log.Printf("[Phoenix]连接MySQL时出错：%v", err)
	// 	return nil
	// }

	return &DBManager{}
}

// 启用日志
func (self *DBManager) EnableLog() {
	self.isVerbose = true

	orm.Debug = true
}

func (self *DBManager) Run() error {
	return orm.RunSyncdb("default", false, self.isVerbose)
}
