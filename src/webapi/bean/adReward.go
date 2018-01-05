package bean

import (
	// "log"
	"time"

	"github.com/astaxie/beego/orm"
)

type AdReward struct {
	UID        uint32 `orm:"column(uid);pk"`
	Num        int32  `orm:"column(number)"`
	CD         int64  `orm:"column(cd)"`
	UpdateTime int64  `orm:"column(update_time)"`
}

func (a *AdReward) Update() {
	_, err := defaultOrm.Update(a)
	checkError("更新广告奖励数据,错误:", err)
}

func (a *AdReward) Reset() *AdReward {
	a.UpdateTime = nextDayZero().Unix()
	a.Num = 0
	a.CD = 0
	a.Update()
	return a
}

func CreateAdReward(uid uint32) *AdReward {
	adReward := &AdReward{
		UID:        uid,
		UpdateTime: nextDayZero().Unix(),
	}
	defaultOrm.Insert(adReward)
	return adReward
}

func LoadAdReward(uid uint32) *AdReward {
	adReward := &AdReward{
		UID: uid,
	}
	err := defaultOrm.Read(adReward)
	switch err {
	case nil:
		if time.Now().Unix() >= adReward.UpdateTime {
			adReward.Reset()
		}
		return adReward
	case orm.ErrNoRows:
		// log.Println("旧版本玩家, 新建广告奖励数据")
		return CreateAdReward(uid)
	default:
		checkError("加载玩家广告奖励数据,错误:", err)
		return adReward
	}
}
