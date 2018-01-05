package bean

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"

	. "webapi/common"
)

const (
	kInsertCountPerLine = 100 // 并行插入数目

	NormalCDKey = uint32(0) // 普通礼包
	DailyCDKey  = uint32(1) // 每日礼包
)

type CDKey struct {
	ID         string `orm:"pk;column(id);size(9)"`   // 兑换码
	Category   uint32 `orm:"column(category)"`        // 兑换码类型[0:普通,1:每日礼包]
	Gift       int    `orm:"column(gift)"`            // 礼品号
	ChannelID  string `orm:"column(channel);size(6)"` // 限定渠道(空字符串代表不限定)
	CreatedAt  int64  `orm:"column(created_at)"`      // 创建时间
	Deadline   int64  `orm:"column(deadline)"`        // 有效截止时间
	Uid        uint32 `orm:"column(uid)"`             // 领取者
	AchievedAt int64  `orm:"column(achieved_at)"`     // 领取时间
}

func (self *CDKey) TableName() string {
	return "cdkey"
}

func (self *CDKey) SetAchieved(uid uint32) {
	self.Uid = uid
	self.AchievedAt = time.Now().Unix()
}

type CDKeyGift struct {
	Gift      int    `orm:"pk;auto;column(gift)"` // 礼品号
	Resources string `orm:"column(resources)"`    // 资源
}

func (self *CDKeyGift) TableName() string {
	return "cdkey_gift"
}

func (self *CDKeyGift) GetResources() map[int32]int32 {
	results := make(map[int32]int32)

	arr1 := strings.Split(self.Resources, ",")
	if n1 := len(arr1); n1 > 0 {
		for i := 0; i < n1; i++ {
			v := arr1[i]
			arr2 := strings.Split(v, ":")
			if n2 := len(arr2); n2 == 2 {
				v1, _ := strconv.ParseInt(arr2[0], 10, 32)
				v2, _ := strconv.ParseInt(arr2[1], 10, 32)
				if v1 > 0 && v2 > 0 {
					results[int32(v1)] = int32(v2)
				}
			}
		}
	}

	return results
}

func ValidCDKeyCategory(category uint32) bool {
	switch category {
	case NormalCDKey, DailyCDKey:
		return true
	}

	return false
}

func AddCDKeyGift(resources map[int32]int32) (error, *CDKeyGift) {
	i := 0
	arr := make([]string, len(resources))
	for id, count := range resources {
		arr[i] = fmt.Sprintf("%d:%d", id, count)
		i++
	}

	o := orm.NewOrm()

	b := &CDKeyGift{
		Resources: strings.Join(arr, ","),
	}

	_, e := o.Insert(b)
	if e != nil {
		log.Printf("新增兑换码礼包记录时出错: %v", e)
		return e, nil
	}

	return nil, b
}

func queryCDKey(id string) (error, *CDKey) {
	o := orm.NewOrm()
	bean := CDKey{
		ID: id,
	}

	err := o.Read(&bean)
	if err != nil {
		// 未找到
		if err == orm.ErrNoRows {
			return nil, nil
		}

		return err, nil
	}

	return nil, &bean
}

func queryCDKeyGift(id int) (error, *CDKeyGift) {
	o := orm.NewOrm()
	bean := CDKeyGift{
		Gift: id,
	}

	err := o.Read(&bean)
	if err != nil {
		// 未找到
		if err == orm.ErrNoRows {
			return nil, nil
		}

		return err, nil
	}

	return nil, &bean
}

// 判断角色是否已领取过同一礼包的兑换码
func IsExchangeGiftCDKEY(uid uint32, gift int) (error, bool) {
	o := orm.NewOrm()

	qs := o.QueryTable("cdkey")
	n, err := qs.Filter("uid", uid).Filter("gift", gift).Count()
	if err != nil {
		return err, false
	}

	return nil, n > 0
}

// 判断角色今日是否已领取过同一礼包的兑换码
func IsExchangeDailyGiftCDKey(uid uint32, gift int) (error, bool) {
	t := ZeroAclock().AddDate(0, 0, -1).Unix()

	o := orm.NewOrm()

	qs := o.QueryTable("cdkey")
	n, err := qs.Filter("uid", uid).Filter("gift", gift).Filter("achieved_at__gte", t).Count()
	if err != nil {
		println(err.Error())
		return err, false
	}

	return nil, n > 0
}

func IsGiftIDExists(id int) (error, bool) {
	o := orm.NewOrm()

	qs := o.QueryTable("cdkey_gift")
	n, err := qs.Filter("gift", id).Count()
	if err != nil {
		return err, false
	}

	return nil, n > 0
}

// 查询兑换码记录
func QueryCDKey(id string) (error, *CDKey, *CDKeyGift) {
	err, bean1 := queryCDKey(id)
	if err != nil {
		return err, nil, nil
	}

	if bean1 == nil {
		return nil, nil, nil
	}

	err, bean2 := queryCDKeyGift(bean1.Gift)
	if err != nil {
		return err, nil, nil
	}

	if bean2 == nil {
		return nil, nil, nil
	}

	return nil, bean1, bean2
}

// 批量新增兑换码
func BatchAddCDKeyWithGift(category uint32, channel string, gift int, deadline int64, ids []string) (error, []string) {
	cdkeys := make([]CDKey, len(ids))

	for i := 0; i < len(ids); i++ {
		cdkeys[i] = CDKey{
			ID:         ids[i],
			Category:   category,
			Gift:       gift,
			ChannelID:  channel,
			CreatedAt:  time.Now().Unix(),
			Deadline:   deadline,
			Uid:        0,
			AchievedAt: 0,
		}
	}

	o := orm.NewOrm()

	success, err := o.InsertMulti(kInsertCountPerLine, cdkeys)
	if err != nil {
		if strings.HasPrefix(err.Error(), "Error 1062: Duplicate entry ") {
			if success == 0 {
				return nil, []string{}
			}

			return nil, ids[:success]
		}

		log.Printf("并行插入兑换码记录时出错: %v %d", err, success)
		return err, nil
	}

	return nil, ids[:success]
}

func UpdateCDKey(bean *CDKey) error {
	o := orm.NewOrm()
	_, e := o.Update(bean)
	if e != nil {
		return e
	}
	return nil
}
