package bean

import (
	// "log"
	"time"

	"github.com/astaxie/beego/orm"
)

type WinReward struct {
	UID        uint32 `orm:"column(uid);pk"`
	Win        int32  `orm:"column(win)"`
	UpdateTime int64  `orm:"column(update_time)"`
}

func (w *WinReward) Update() {
	_, err := defaultOrm.Update(w)
	checkError("更新玩家胜场奖励数据,错误:", err)
}

func (w *WinReward) Reset() *WinReward {
	w.UpdateTime = nextDayZero().Unix()
	w.Win = 0
	w.Update()
	return w
}

func nextDayZero() time.Time {
	return time.Now().Truncate(time.Hour * 24).Add(16 * time.Hour)
}

func CreateWinReward(uid uint32) *WinReward {
	w := &WinReward{
		UID:        uid,
		UpdateTime: nextDayZero().Unix(),
	}
	defaultOrm.Insert(w)
	return w
}

func LoadWinReward(uid uint32) *WinReward {
	w := &WinReward{
		UID: uid,
	}
	err := defaultOrm.Read(w)

	switch err {
	case nil:
		if time.Now().Unix() >= w.UpdateTime {
			w.Reset()
		}
		return w
	case orm.ErrNoRows:
		// log.Println("旧版本玩家, 新建胜场奖励数据")
		return CreateWinReward(uid)
	default:
		checkError("加载玩家胜场奖励数据,错误:", err)
		return w
	}
}
