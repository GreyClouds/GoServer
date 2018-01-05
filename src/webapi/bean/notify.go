package bean

import (
	"encoding/base64"
	"log"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/golang/protobuf/proto"

	pkgProto "crazyant.com/deadfat/pbd/go"
)

// ---------------------------------------------------------------------------

const (
	NOTIFY_PAYMENT = uint32(1)
)

type Notify struct {
	ID       int64     `orm:"column(id);auto;pk"`     // 主键
	Uid      uint32    `orm:"column(uid);index"`      // 角色编号
	Happen   time.Time `orm:"column(happen)"`         // 发生时间
	Category uint32    `orm:"column(category)"`       // 类别
	Raw      string    `orm:"column(raw);size(1024)"` // 内容
	Removed  bool      `orm:"column(removed)"`        // 是否标记删除
	dirty    bool      `orm:"-"`                      // 是否脏数据
}

func (self *Notify) TableName() string {
	return "notify"
}

func (self *Notify) SetDirty() {
	self.dirty = true
}

func (self *Notify) ResetDirty() bool {
	if !self.dirty {
		return false
	}

	self.dirty = false

	return true
}

func (self *Notify) SetRaw(payload proto.Message) error {
	raw, err := proto.Marshal(payload)
	if err != nil {
		return err
	}

	self.Raw = base64.StdEncoding.EncodeToString(raw)

	return nil
}

func (self *Notify) GetRaw(payload proto.Message) error {
	raw, err := base64.StdEncoding.DecodeString(self.Raw)
	if err != nil {
		return err
	}

	return proto.Unmarshal(raw, payload)
}

// ---------------------------------------------------------------------------

func LoadNotifies(uid uint32) (error, []*Notify) {
	o := orm.NewOrm()

	var beans []*Notify
	qs := o.QueryTable("notify")
	_, e := qs.Filter("uid", uid).Filter("removed", false).All(&beans)
	if e != nil {
		return e, nil
	}

	return nil, beans
}

func UpdateNotify(bean *Notify) error {
	o := orm.NewOrm()
	_, e := o.Update(bean)
	if e != nil {
		return e
	}

	return nil
}

func AddAndroidPaymentNotify(orderid string, uid uint32, isRemovedAD bool, resources map[int32]int32, money *pkgProto.Money) (*Notify, error) {
	notify := &Notify{
		Uid:      uid,
		Happen:   time.Now(),
		Category: NOTIFY_PAYMENT,
		Removed:  false,
	}

	payload := &pkgProto.PaymentNotify{
		Id:       proto.String(orderid),
		RemoveAd: proto.Bool(isRemovedAD),
		Money:    money,
	}

	for id, count := range resources {
		payload.Resources = append(payload.Resources, &pkgProto.Resource{
			Id:    proto.Int32(id),
			Value: proto.Int32(count),
		})
	}

	notify.SetRaw(payload)

	o := orm.NewOrm()
	if _, err := o.Insert(notify); err != nil {
		log.Printf("新增角色%d通知记录时出错: %v", uid, err)
		return nil, err
	}

	return notify, nil
}
