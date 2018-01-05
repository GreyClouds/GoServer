package bean

import (
	"log"
	"time"

	"github.com/astaxie/beego/orm"
)

type Payment struct {
	ID          int64     `orm:"column(id);auto;pk"`  // 主键
	Uid         uint32    `orm:"column(uid);index"`   // 角色编号
	Happen      time.Time `orm:"column(happen)"`      // 发生时间
	Provider    int       `orm:"column(provider)"`    // 提供方[1:AppStore]
	Tracscation string    `orm:"column(tracscation)"` // 原始交易号
	SKU         string    `orm:"column(sku)"`         // SKU
}

func (self Payment) TableName() string {
	return "payment"
}

// ---------------------------------------------------------------------------

func ExistsAppStorePayment(uid uint32, transcation string) (error, bool) {
	o := orm.NewOrm()

	qs := o.QueryTable("payment")
	n, e := qs.Filter("uid", uid).Filter("tracscation", transcation).Count()
	if e != nil {
		return e, false
	}

	return nil, (n > 0)
}

func NewAndSaveAppStorePayment(happen time.Time, uid uint32, tracscation string, sku string) (error, *Payment) {
	o := orm.NewOrm()

	b := &Payment{
		Uid:         uid,
		Happen:      happen,
		Provider:    1,
		Tracscation: tracscation,
		SKU:         sku,
	}

	_, e := o.Insert(b)
	if e != nil {
		log.Printf("角色%d插入付费记录时出错: %v", uid, e)
		return e, nil
	}

	return nil, b
}
