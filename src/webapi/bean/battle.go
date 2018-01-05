package bean

import (
	"fmt"
	"log"
	"time"

	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/orm"
)

const (
	kTBLBattleRoom = "battle_room"
)

type BattleRoom struct {
	RoomID   uint32 `orm:"column(roomid);pk"`  // 房间编号
	BeginTS  int64  `orm:"column(begin_ts)"`   // 匹配结束的时间
	P1Uid    uint32 `orm:"column(uid)"`        // 玩家1编号
	P1SkinID int32  `orm:"column(p1_skin_id)"` // 玩家1皮肤编号
	P2Uid    uint32 `orm:"column(p2_uid)"`     // 玩家2uid
	P2SkinID int32  `orm:"column(p2_skin_id)"` // 玩家2皮肤编号
	PkTypes  int32  `orm:"column(pk_types)"`   // 战斗类型 0 常规 1 好友
}

func (b *BattleRoom) TableName() string {
	return kTBLBattleRoom
}

// ----------------------------------------------------------------------------------------------------------------------

var gBattleRoomCache cache.Cache

func init() {
	gBattleRoomCache = cache.NewMemoryCache()
}

func NewBattleRecord(roomid uint32, ts int64, p1Uid, p2Uid uint32, p1Skin, p2Skin, pkTypes int32) (error, *BattleRoom) {
	record := &BattleRoom{
		RoomID:   roomid,
		BeginTS:  ts,
		P1Uid:    p1Uid,
		P1SkinID: p1Skin,
		P2Uid:    p2Uid,
		P2SkinID: p2Skin,
		PkTypes:  pkTypes,
	}

	gBattleRoomCache.Put(fmt.Sprintf("%d", roomid), record, 120*time.Second)
	DefaultKeeper().SendBattleRoom(record)

	// log.Printf("玩家%d,玩家%d匹配成功至房间%d", p1Uid, p2Uid, roomid)
	return nil, record
}

func BattleRoomExists(roomid uint32) (error, *BattleRoom) {
	if bean := gBattleRoomCache.Get(fmt.Sprintf("%d", roomid)); bean != nil {

		return nil, bean.(*BattleRoom)
	} else {
		o := orm.NewOrm()

		var bean BattleRoom
		qs := o.QueryTable(kTBLBattleRoom)
		if err := qs.Filter("roomid", roomid).One(&bean); err != nil {
			return err, nil
		}

		return nil, &bean
	}
}

func BattleRoomCleanTask() {
	t := time.Now().Add(-90 * 24 * time.Hour)

	o := orm.NewOrm()
	result, err := o.Raw("delete from battle_room where begin_ts < ?", t.Unix()).Exec()
	if err != nil {
		log.Printf("清理历史房间记录时出错1: %v", err)
		return
	}

	num, err := result.RowsAffected()
	if err != nil {
		log.Printf("清理历史房间记录时出错2: %v", err)
		return
	}

	if num > 0 {
		log.Printf("清理历史房间记录共%d条", num)
	}
}
