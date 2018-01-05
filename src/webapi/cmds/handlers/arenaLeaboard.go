package handlers

import (
	"github.com/golang/protobuf/proto"

	pbd "crazyant.com/deadfat/pbd/hero"

	oaccount "webapi/account"
	obean "webapi/bean"

	osession "webapi/session"
	oskeleton "webapi/skeleton"
)

func HandleGetLeaboardRange(skeleton *oskeleton.Skeleton, session *osession.Session, role *oaccount.Role, packet proto.Message) (error, uint16, proto.Message) {
	payload, _ := packet.(*pbd.GetLearderboardRange)
	resp := &pbd.GetLearderboardRangeResp{}
	name := payload.Name
	if name == "ArenaLeaderboard"{
		resp.Infos = obean.GetArenaLeaderboard(int(payload.Start), int(payload.End))
	}else{
		resp.Infos = obean.GetChallLeaderboard(int(payload.Start), int(payload.End))
	}
	resp.Count = uint32(len(resp.Infos));

	return nil, uint16(pbd.GC_ID_GET_LEARDERBOARD_RANGE_RESP), resp
}

func HandleLearderBoardUpdae(skeleton *oskeleton.Skeleton, session *osession.Session, role *oaccount.Role, packet proto.Message) (error, uint16, proto.Message) {
	payload, _ := packet.(*pbd.UpdateLearderBoardScore)
	resp := &pbd.UpdateLearderBoardScoreResp{}

	name := payload.Name
	if name == "ArenaLeaderboard"{
		resp.Success, resp.NewRank = obean.DoUpdateArenaScore(role.Uid(), uint32(payload.Score))
	}else{
		resp.Success, resp.NewRank = obean.DoUpdateChallScore(role.Uid(), payload.Score)
	}


	return nil, uint16(pbd.GC_ID_UPDATE_LEARD_SCOREE_RESP), resp
}
