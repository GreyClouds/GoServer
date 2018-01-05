package bean

import (
	"log"
	"time"

	"github.com/astaxie/beego/orm"
)

const (
	kTBLBattleResult = "battle_result"
)

type BattleResult struct {
	ID     uint32 `orm:"column(id);auto;pk"`   // 编号
	RoomID uint32 `orm:"column(roomid);index"` // 房间编号
	Uid    uint32 `orm:"column(uid);index"`    // 玩家编号
	TS     int64  `orm:"column(ts)"`           // 提交时间点
	Result int32  `orm:"column(result)"`       // 结果
	Score  int32  `orm:"column(score)"`        // 分数变动
}

func (b BattleResult) TableName() string {
	return kTBLBattleResult
}

// ----------------------------------------------------------------------------------------------------------------------

func NewBattleResultRecord(roomid, uid uint32, ts int64, result, score int32) (error, *BattleResult) {
	record := &BattleResult{
		RoomID: roomid,
		TS:     ts,
		Uid:    uid,
		Result: result,
		Score:  score,
	}

	o := orm.NewOrm()
	_, err := o.Insert(record)
	if err != nil {
		log.Printf("角色%d新增结算记录时出错: %v", uid, err)
		return err, nil
	}

	// log.Printf("玩家%d提交结算信息", uid)
	return nil, record
}

func BattleResultCleanTask() {
	t := time.Now().Add(-90 * 24 * time.Hour)

	o := orm.NewOrm()
	result, err := o.Raw("delete from battle_result where ts < ?", t.Unix()).Exec()
	if err != nil {
		log.Printf("清理历史房间结算时出错1: %v", err)
		return
	}

	num, err := result.RowsAffected()
	if err != nil {
		log.Printf("清理历史房间结算时出错2: %v", err)
		return
	}

	if num > 0 {
		log.Printf("清理历史房间结算共%d条", num)
	}
}
