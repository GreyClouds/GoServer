package activity

import (
	"time"

	pkgPBD "crazyant.com/deadfat/pbd/go"

	"github.com/golang/protobuf/proto"
)

type Activity struct {
	id         int32
	start      int64
	end        int64
	skins      map[int32]int32
	skinIDList []int32
}

func newActivity(id int32, start int64, end int64, skins map[int32]int32) *Activity {
	arr := make([]int32, len(skins))

	var i int
	for k, _ := range skins {
		arr[i] = k
		i++
	}

	return &Activity{
		id:         id,
		start:      start,
		end:        end,
		skins:      skins,
		skinIDList: arr,
	}
}

// 是否开放
func (self *Activity) IsOpened() bool {
	current := time.Now().Unix()

	return self.start <= current && current <= self.end
}

// 是否允许试用皮肤
func (self *Activity) IsSkinTry(id int32) bool {
	_, exists := self.skins[id]
	if exists {
		return true
	}

	return false
}

// 剩余天数
func (self *Activity) RemainDays() int32 {
	dt := self.end - time.Now().Unix()
	if dt <= 0 {
		return 0
	}

	return int32((dt + 86399) / 86400)
}

// 判断任务是否达成
func (self *Activity) IsTaskEnough(id, current int32) bool {
	need, exists := self.skins[id]
	if !exists {
		return false
	}

	if current >= need {
		return true
	}

	return false
}

func (self *Activity) Serialize(participant IParticipant) *pkgPBD.Activity {
	payload := &pkgPBD.Activity{
		ActId:   proto.Int32(self.id),
		EndTime: proto.Int32(self.RemainDays()),
	}

	if tasks := participant.SerializeSkinTask(self.skinIDList); tasks != nil && len(tasks) > 0 {
		payload.SkinWinData = tasks
	}

	return payload
}
