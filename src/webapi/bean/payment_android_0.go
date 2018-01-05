package bean

import (
	"errors"
	"log"

	"github.com/astaxie/beego/orm"
)

var (
	ErrPaymentAlreadyAchieved = errors.New("payment already achieved")
)

// 安卓支付
type AndroidPayment struct {
	OrderID  string `orm:"column(order_id);pk"`   // SDK订单号
	UUID     string `orm:"column(uuid)"`          // SDK唯一用户编号
	ZoneID   int32  `orm:"column(zone_id)"`       // 游戏大区编号
	UID      uint32 `orm:"column(uid)"`           // 游戏角色编号
	SKU      string `orm:"column(sku);size(128)"` // 商品编号
	Amount   int32  `orm:"column(amount)"`        // 商品价格
	PayTime  int64  `orm:"column(pay_time)"`      // 订单支付时间
	Sandbox  bool   `orm:"column(sandbox)"`       // 是否测试订单
	Happen   int64  `orm:"column(happen)"`        // 收到记录时间
	Achieved bool   `orm:"column(achieved)"`
}

func (self *AndroidPayment) TableName() string {
	return "android_payment"
}

func (bean *AndroidPayment) UpdateAsAchieved() error {
	// 判断是否已领取
	if bean.Achieved {
		return ErrPaymentAlreadyAchieved
	}

	bean.Achieved = true

	o := orm.NewOrm()
	_, e := o.Update(bean)
	if e != nil {
		return e
	}

	return nil
}

func InsertAndroidPayment(bean *AndroidPayment) error {
	o := orm.NewOrm()

	_, e := o.Insert(bean)
	if e != nil {
		log.Printf("新增FYSDK支付订单记录%v时出错: %v", bean, e)
		return e
	}

	return nil
}

func LoadAndroidPayments(uid uint32) (error, []*AndroidPayment) {
	o := orm.NewOrm()

	var beans []*AndroidPayment
	qs := o.QueryTable("android_payment")
	_, e := qs.Filter("uid", uid).Filter("achieved", false).All(&beans)
	if e != nil {
		return e, nil
	}

	return nil, beans
}
