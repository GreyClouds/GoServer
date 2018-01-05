package bean

import (
	"log"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

const (
	kInviteCodeInsertCountPerLine = 100 // 并行插入数目
)

type InviteCode struct {
	ID         string `orm:"pk;column(id);size(8)"`   // 兑换码
	ChannelID  string `orm:"column(channel);size(6)"` // 限定渠道
	CreatedAt  int64  `orm:"column(created_at)"`      // 创建时间
	Deadline   int64  `orm:"column(deadline)"`        // 有效截止时间
	Uid        uint32 `orm:"column(uid)"`             // 领取者
	AchievedAt int64  `orm:"column(achieved_at)"`     // 领取时间
}

func (self *InviteCode) TableName() string {
	return "invite_code"
}

func (self *InviteCode) SetAchieved(uid uint32) {
	self.Uid = uid
	self.AchievedAt = time.Now().Unix()
}

func QueryInviteCode(id string) (error, *InviteCode) {
	o := orm.NewOrm()
	bean := InviteCode{
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

// 批量新增兑换码
func BatchAddInviteCode(channel string, deadline int64, ids []string) (error, []string) {
	cdkeys := make([]InviteCode, len(ids))

	for i := 0; i < len(ids); i++ {
		cdkeys[i] = InviteCode{
			ID:         ids[i],
			ChannelID:  channel,
			CreatedAt:  time.Now().Unix(),
			Deadline:   deadline,
			Uid:        0,
			AchievedAt: 0,
		}
	}

	o := orm.NewOrm()

	success, err := o.InsertMulti(kInviteCodeInsertCountPerLine, cdkeys)
	if err != nil {
		if strings.HasPrefix(err.Error(), "Error 1062: Duplicate entry ") {
			if success == 0 {
				return nil, []string{}
			}

			return nil, ids[:success]
		}

		log.Printf("并行插入邀请码记录时出错: %v %d", err, success)
		return err, nil
	}

	return nil, ids[:success]
}

func UpdateInviteCode(bean *InviteCode) error {
	o := orm.NewOrm()
	_, e := o.Update(bean)
	if e != nil {
		return e
	}
	return nil
}
