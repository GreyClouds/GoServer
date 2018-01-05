package activity

import (
	pkgPBD "crazyant.com/deadfat/pbd/go"
)

type IParticipant interface {
	// 累计皮肤任务胜场数
	GainSkinTask(id int32) int32

	// 获得皮肤
	AddSkin(id int32, src int) bool

	// 是否拥有皮肤
	SkinOwned(skin int32) bool

	// 序列化
	SerializeSkinTask(skins []int32) []*pkgPBD.SkinWinData
}
