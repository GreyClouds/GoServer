package handlers

import (
	"github.com/golang/protobuf/proto"

	pbd "crazyant.com/deadfat/pbd/hero"

	oaccount "webapi/account"
	obean "webapi/bean"

	osession "webapi/session"
	oskeleton "webapi/skeleton"
)

func HandleSetAchievement(skeleton *oskeleton.Skeleton, session *osession.Session, role *oaccount.Role, packet proto.Message) (error, uint16, proto.Message) {
	payload, _ := packet.(*pbd.SetAchievement)
	obean.SetAchievement(role.Uid(), payload.Name, payload.NowValue)
	resp := &pbd.Empty{}
	return nil, uint16(pbd.GC_ID_SET_ACHIEVEMENT_RESP), resp
}


func HandleGetAchievement(skeleton *oskeleton.Skeleton, session *osession.Session, role *oaccount.Role, packet proto.Message) (error, uint16, proto.Message) {
	payload, _ := packet.(*pbd.GetAchievement)
	resp := &pbd.GetAchievementResp{}
	resp.Value = obean.GetAchievement(role.Uid(), payload.Name)
	return nil, uint16(pbd.GC_ID_SET_ACHIEVEMENT_RESP), resp
}

func HandleUnLockAchievement(skeleton *oskeleton.Skeleton, session *osession.Session, role *oaccount.Role, packet proto.Message) (error, uint16, proto.Message) {
	payload, _ := packet.(*pbd.UnLockAchievement)
	resp := &pbd.UnLockAchievementResp{}
	resp.Date = obean.UnLockAchievement(role.Uid(), payload.Name)
	return nil, uint16(pbd.GC_ID_UNLOCK_ACHIEVEMENT_RESP), resp
}

func GetUnLockAchievement(skeleton *oskeleton.Skeleton, session *osession.Session, role *oaccount.Role, packet proto.Message) (error, uint16, proto.Message) {
	payload, _ := packet.(*pbd.GetUnLockAchieveDate)
	resp := &pbd.GetUnLockAchieveDateResp{}
	resp.Date = obean.GetUnLockAchieveDate(role.Uid(), payload.Name)
	return nil, uint16(pbd.GC_ID_GET_UNLOCK_ACHIEVEDATE_RESP), resp
}
