package handlers

import (
	"log"


	"github.com/golang/protobuf/proto"

	pbd "crazyant.com/deadfat/pbd/hero"

	oaccount "webapi/account"
	obean "webapi/bean"

	osession "webapi/session"
	oskeleton "webapi/skeleton"

)

//func HandleGuestLogin(skeleton *oskeleton.Skeleton, session *osession.Session, _ *oaccount.Role, packet proto.Message) (error, uint16, proto.Message) {
//	payload, _ := packet.(*pbd.Login)
//
//	resp := &pbd.LoginResp{}
//
//	err, account := obean.FindSimAccount(payload.Account)
//	if err == nil {
//		resp.Uid = account.Uid
//		resp.ArenaScore =obean.GetScore(account.Uid)
//		resp.ArenaRank = obean.GetRank(account.Uid)
//		resp.ChallengeRank = obean.GetChallRank(account.Uid)
//		resp.ChallengeScore = obean.GetChallScore(account.Uid)
//		resp.Token = skeleton.CreateUserSecret(account.Uid)
//	} else{
//		//没有注册过的账号，注册
//		resp.Err = doHeroError(kRoleNotExists)
//	}
//
//	return nil, uint16(pbd.GC_ID_GUEST_LOGIN_RESP), resp
//}

func HandleRegister(skeleton *oskeleton.Skeleton, session *osession.Session, _ *oaccount.Role, packet proto.Message) (error, uint16, proto.Message) {
	payload, _ := packet.(*pbd.GuestRegister)

	resp := &pbd.LoginResp{}


	pCache := skeleton.CacheManager()
	pAccountManager := skeleton.AccountManager()

	err, _ := obean.FindSimAccount(payload.Account)
	if err == nil {
		resp.Err = doHeroError(kAlreadyRegistered)
	} else{
		var uid uint32
		//没有注册过的账号，注册
		if uid = pCache.GenID("uid", 4); uid == 0 {
			log.Printf("生成角色编号失败: %d", payload.Account)
			resp.Err = doHeroError(kInternelServerError)
			return nil, uint16(pbd.GC_ID_GUEST_LOGIN_RESP), resp
		} else {
			if err := obean.RegisterAccount(payload.Account, payload.Name, uid); err != nil {
				log.Printf("插入新角色%d账号数据时出错: %v", uid, err)
				resp.Err = doHeroError(kInternelServerError)
			} else {
				resp.Uid = uid
				resp.ArenaRank = 0
				resp.ArenaScore = 0
				resp.ChallengeRank = 0
				resp.ChallengeScore = 0
				resp.Token = skeleton.CreateUserSecret(uid)
				pAccountManager.LoadRole(uid)
			}
		}
	}

	return nil, uint16(pbd.GC_ID_GUEST_LOGIN_RESP), resp
}
//func HandleNickSet(_ *oskeleton.Skeleton, _ *osession.Session, roleob *oaccount.Role, packet proto.Message) (error, uint16, proto.Message) {
//	payload, _ := packet.(*pbd.NickSet)
//	resp := &pbd.NickSetResp{}
//
//	if roleob.GetGuideID() != obean.GUIDE_NICK_SET {
//		resp.Err = doError(kNickNameHaveSet)
//		return nil, uint16(pbd.GC_ID_NICK_SET_RESP), respHero
//	}
//
//	nick := strings.TrimSpace(payload.GetNick())
//	if nick == "" {
//		resp.Err = doError(kNickIsSensitivity)
//		return nil, uint16(pbd.GC_ID_NICK_SET_RESP), resp
//	}
//
//	lang := payload.GetLang()
//
//	switch lang {
//	case 1: // 简体中文
//		if valid := wordstock.ValidWord(nick); !valid {
//			resp.Err = doError(kNickIsSensitivity)
//			return nil, uint16(pbd.GC_ID_NICK_SET_RESP), resp
//		}
//
//	default:
//	}
//
//	roleob.SetNickAndGuide(nick)
//
//	return nil, uint16(pbd.GC_ID_NICK_SET_RESP), resp
//}

