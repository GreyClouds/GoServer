package handlers

import (
	"log"
	"time"

	"webapi/bean"

	"github.com/golang/protobuf/proto"

	pbd "crazyant.com/deadfat/pbd/go"

	oaccount "webapi/account"
	osession "webapi/session"
	oskeleton "webapi/skeleton"
)

// 获取今日剩余秒数
func getTodayRemainSeconds() int32 {
	v := time.Now().Unix()
	return 86400 - int32((v+28800)%86400)
}

func HandleSigninView(skeleton *oskeleton.Skeleton, session *osession.Session, roleob *oaccount.Role, packet proto.Message) (error, uint16, proto.Message) {
	resp := &pbd.SigninViewResp{}

	pConfManager := skeleton.ConfigManager()

	ts, id := roleob.GetLastSigninTask()
	if ts > 0 {
		last := time.Unix(ts, 0)
		y1, m1, d1 := last.Date()
		y2, m2, d2 := time.Now().Date()
		if y1 == y2 && m1 == m2 && d1 == d2 {
			if cc2 := pConfManager.GetNextSigninTask(id); cc2 != nil {
				resp.Id = proto.Int32(cc2.ID)
				resp.Sec = proto.Int32(getTodayRemainSeconds())
			}

			return nil, uint16(pbd.GC_ID_SIGNIN_VIEW_RESP), resp
		}
	}

	if cc2 := pConfManager.GetNextSigninTask(id); cc2 != nil {
		resp.Id = proto.Int32(cc2.ID)
		resp.Sec = proto.Int32(0)
	}

	return nil, uint16(pbd.GC_ID_SIGNIN_VIEW_RESP), resp
}

func HandleSignin(skeleton *oskeleton.Skeleton, session *osession.Session, roleob *oaccount.Role, packet proto.Message) (error, uint16, proto.Message) {
	p, _ := packet.(*pbd.SigninForm)
	_ = p

	pConfManager := skeleton.ConfigManager()
	taskId := p.Id

	resp := &pbd.SigninResp{}

	ts, id := roleob.GetLastSigninTask()
	if ts > 0 {
		last := time.Unix(ts, 0)
		y1, m1, d1 := last.Date()
		y2, m2, d2 := time.Now().Date()
		if y1 == y2 && m1 == m2 && d1 == d2 {
			// log.Printf("角色%d今日已签到: %d", roleob.uid(), taskId)
			resp.Err = doError(kRoleTodayAlreadySignin)

			if cc2 := pConfManager.GetNextSigninTask(id); cc2 != nil {
				resp.NextId = proto.Int32(cc2.ID)
				resp.NextSec = proto.Int32(getTodayRemainSeconds())
			}

			return nil, uint16(pbd.GC_ID_SIGNIN_RESP), resp
		}
	}

	uid := roleob.Uid()

	if id != 0 && id == taskId {
		log.Printf("角色%d重复签到: %d", uid, id)

		resp.Err = doError(kRoleTodayAlreadySignin)
		return nil, uint16(pbd.GC_ID_SIGNIN_RESP), resp
	}

	cc := pConfManager.GetNextSigninTask(id)
	if cc == nil {
		log.Printf("角色%d的签到任务%d下次有效任务配置数据不存在: %d", uid, id, taskId)

		resp.Err = doError(kInternelServerError)
		return nil, uint16(pbd.GC_ID_SIGNIN_RESP), resp
	}

	if cc.ID != taskId {
		log.Printf("角色%d下次可签到任务%d != %d", uid, taskId, cc.ID)
	}

	roleob.Signin(cc.ID)

	res := make(map[int32]int32)
	res[cc.Reward] = 1
	roleob.AddSkin(cc.Reward, bean.SkinSrcTaskReward)

	resources := []*pbd.Resource{}
	for k, v := range res {
		resources = append(resources, &pbd.Resource{
			Id:    proto.Int32(k),
			Value: proto.Int32(v),
		})
	}

	resp.Id = proto.Int32(cc.ID)

	if cc2 := pConfManager.GetNextSigninTask(cc.ID); cc2 != nil {
		resp.NextId = proto.Int32(cc2.ID)
		resp.NextSec = proto.Int32(getTodayRemainSeconds())
	}

	resp.Resources = resources

	return nil, uint16(pbd.GC_ID_SIGNIN_RESP), resp
}
