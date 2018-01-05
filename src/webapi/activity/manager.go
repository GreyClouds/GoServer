package activity

import (
	"log"
	"time"

	"webapi/bean"

	// pkgSD "crazyant.com/deadfat/data/go"
	pkgPBD "crazyant.com/deadfat/pbd/go"

	pkgConfig "webapi/config"
)

type ActivityManager struct {
	activitys map[int32]*Activity
}

func (self *ActivityManager) Preload(confs []*pkgConfig.ActivityConf) {
	currentTS := time.Now().Unix()

	for i := 0; i < len(confs); i++ {
		conf := confs[i]

		// 判断活动是否已过期
		if conf.EndTS < currentTS {
			continue
		}

		// 判断活动是否重复
		if _, exists := self.activitys[conf.ID]; exists {
			log.Printf("活动%d重复配置", conf.ID)
			continue
		}

		self.activitys[conf.ID] = newActivity(conf.ID, conf.StartTS, conf.EndTS, conf.Skins)
	}
}

func (self *ActivityManager) Serialize(participant IParticipant) []*pkgPBD.Activity {
	results := []*pkgPBD.Activity{}

	for _, obj := range self.activitys {
		// 活动已关闭
		if ok := obj.IsOpened(); !ok {
			continue
		}

		results = append(results, obj.Serialize(participant))
	}

	if len(results) > 0 {
		return results
	}

	return nil
}

// 判断是否允许试用皮肤
func (self *ActivityManager) IsSkinTry(skin int32) bool {
	for _, obj := range self.activitys {
		// 活动已关闭
		if ok := obj.IsOpened(); !ok {
			continue
		}

		if ok := obj.IsSkinTry(skin); ok {
			return true
		}
	}

	return false
}

// 皮肤任务结算
func (self *ActivityManager) SkinTaskReward(participant IParticipant, skin int32) []int32 {
	results := []int32{}

	for _, obj := range self.activitys {
		// 活动已关闭
		if ok := obj.IsOpened(); !ok {
			continue
		}

		// 是否拥有这个皮肤
		if exists := participant.SkinOwned(skin); exists {
			continue
		}

		// 是否试用皮肤
		if ok := obj.IsSkinTry(skin); !ok {
			continue
		}

		current := participant.GainSkinTask(skin)
		if enough := obj.IsTaskEnough(skin, current); enough {
			ok := participant.AddSkin(skin, bean.SkinSrcTaskReward)
			if ok {
				results = append(results, skin)
			}
		}
	}

	return results
}
