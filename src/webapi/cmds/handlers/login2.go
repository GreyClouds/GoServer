package handlers

import (
	"fmt"
	"log"
	"time"
	"math/rand"

	"github.com/golang/protobuf/proto"

	pbd "crazyant.com/deadfat/pbd/hero"

	oaccount "webapi/account"
	obean "webapi/bean"
	. "webapi/common"
	FYSDK "webapi/fysdk"
	osession "webapi/session"
	oskeleton "webapi/skeleton"
	pkgTrace "webapi/trace"
)

func HandleGuestLogin(skeleton *oskeleton.Skeleton, session *osession.Session, _ *oaccount.Role, packet proto.Message) (error, uint16, proto.Message) {
	payload, _ := packet.(*pbd.Login)

	resp := &pbd.LoginResp{}

	imei := payload.Imei
	clientVersion := payload.ClientVersion
	clientChannel := payload.Channel
	nickName := payload.NickName

	if imei == "" {

		log.Printf("登录时设备号为空: [%s][%s][%s]", clientChannel, clientVersion, session.IP)
		resp.Err = doHeroError(kIMEIIsEmpty)

		return nil, uint16(pbd.GC_ID_GUEST_LOGIN_RESP), resp
	}

	// log.Printf("设备登录: %s %s", imei, session.IP)

	pCache := skeleton.CacheManager()
	pAccountManager := skeleton.AccountManager()
	pConfManager := skeleton.ConfigManager()

	// 新版本触发版本检查
	if clientChannel != "" {
		haveNewVersion, isForceUpdate, newVersion, downloadLink := pConfManager.CheckNewVersion(clientChannel, clientVersion)
		if haveNewVersion {
			resp.Update = &pbd.VersionUpdateAlert{
				Force:   proto.Bool(isForceUpdate),
				Version: proto.String(newVersion),
				Link:    proto.String(downloadLink),
			}
		}
	}


	channelConf := pConfManager.GetChannel(clientChannel)
	if channelConf == nil {
		log.Printf("查询登录渠道配置时失败: %s", clientChannel)
		resp.Err = doHeroError(kInternelServerError)
		return nil, uint16(pbd.GC_ID_GUEST_LOGIN_RESP), resp
	}
	loginWay := uint32(channelConf.LoginMode)
	loginUUID := imei
	switch loginWay {
	case GUEST_LOGIN:
		loginUUID = imei
	case FYSDK_ONLINE:
		// 校验token
		platform := channelConf.SignName
		token := payload.FysdkToken
		if platform == "" || token == "" {
			log.Printf("FYSDK网游版本登录时platform或token为空: %s, %s", platform, token)
			resp.Err = doHeroError(kIMEIIsEmpty)
			return nil, uint16(pbd.GC_ID_GUEST_LOGIN_RESP), resp
		}

		uuid, err := FYSDK.Login(platform, token)
		if err != nil {
			log.Printf("FYSDK网游版本登录时校验出错: %v", err)
			resp.Err = doHeroError(kInternelServerError)
			return nil, uint16(pbd.GC_ID_GUEST_LOGIN_RESP), resp
		}

		if uuid == "" {
			log.Printf("FYSDK网游版本uuid为空[渠道:%s][IMEI:%s][版本:%s]", clientChannel, imei, clientVersion)
			resp.Err = doHeroError(kIMEIIsEmpty)
			return nil, uint16(pbd.GC_ID_GUEST_LOGIN_RESP), resp
		}

		loginUUID = uuid

	case FYSDK_OFFLINE:
		loginUUID = payload.FysdkUuid
		if loginUUID == "" {
			log.Printf("FYSDK单机版本登录时UUID为空[渠道:%s][IMEI:%s][版本:%s]", clientChannel, imei, clientVersion)
			resp.Err = doHeroError(kIMEIIsEmpty)
			return nil, uint16(pbd.GC_ID_GUEST_LOGIN_RESP), resp
		}
	}

	accountStr := fmt.Sprintf("%s:%s", clientChannel, loginUUID);

	err, account := obean.FindSimAccount(accountStr)
	var uid uint32
	if err != nil {
		//没有注册过的账号，注册
		if uid = pCache.GenID("uid", 1); uid == 0 {
			log.Printf("生成角色编号失败: %d", accountStr)
			resp.Err = doHeroError(kInternelServerError)
			return nil, uint16(pbd.GC_ID_GUEST_LOGIN_RESP), resp
		} else {
			if nickName == ""{
				nickName = fmt.Sprintf("%s%05d","英雄", rand.Intn(99999));
			}
			if err := obean.RegisterAccount(accountStr, nickName, uid); err != nil {
				log.Printf("插入新角色%d账号数据时出错: %v", uid, err)
				resp.Err = doHeroError(kInternelServerError)
			} else {
				resp.Uid = uid
				resp.ArenaRank = 0
				resp.ArenaScore = 0
				resp.ChallengeRank = 0
				resp.ChallengeScore = 0
			}
		}
	}else{
		uid = account.Uid
		if nickName != "" {
			obean.SetNickName(nickName, uid)
		}
		resp.Uid = uid
		resp.ArenaScore =obean.GetScore(uid)
		resp.ArenaRank = obean.GetRank(uid)
		resp.ChallengeRank = obean.GetChallRank(uid)
		resp.ChallengeScore = obean.GetChallScore(uid)
	}
	resp.Token = skeleton.CreateUserSecret(uid)
	pAccountManager.LoadRole(uid)

	// 登陆事件
	skeleton.Collect(pkgTrace.UserLogin,
		clientChannel, imei, loginWay, loginUUID, uid, session.IP, clientVersion,
		time.Now().Unix())

	return nil, uint16(pbd.GC_ID_GUEST_LOGIN_RESP), resp
}

//func HandleNickSet(_ *oskeleton.Skeleton, _ *osession.Session, roleob *oaccount.Role, packet proto.Message) (error, uint16, proto.Message) {
//	payload, _ := packet.(*pbd.NickSet)
//	resp := &pbd.NickSetResp{}
//
//	if roleob.GetGuideID() != obean.GUIDE_NICK_SET {
//		resp.Err = doError(kNickNameHaveSet)
//		return nil, uint16(pbd.GC_ID_NICK_SET_RESP), resp
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

//func HandleInviteCodeInput(skeleton *oskeleton.Skeleton, session *osession.Session, roleob *oaccount.Role, packet proto.Message) (error, uint16, proto.Message) {
//	payload, _ := packet.(*pbd.InviteCodeReq)
//	resp := &pbd.InviteCodeResp{}
//
//	// 判断邀请码是否符合规范
//	code := strings.ToUpper(payload.GetInviteCode())
//	if len(code) != 8 {
//		resp.Err = doError(kInviteCodeInvalid)
//		return nil, uint16(pbd.GC_ID_INVITE_CODE_RESP), resp
//	}
//
//	// 判断邀请码是否存在以及是否未被使用
//	err, bean := obean.QueryInviteCode(code)
//	if err != nil {
//		log.Printf("查询邀请码%s时出错: %v", code, err)
//		resp.Err = doError(kInviteCodeInvalid)
//		return nil, uint16(pbd.GC_ID_INVITE_CODE_RESP), resp
//	}
//
//	// 判断是否存在
//	if bean == nil {
//		resp.Err = doError(kInviteCodeInvalid)
//		return nil, uint16(pbd.GC_ID_INVITE_CODE_RESP), resp
//	}
//
//	// 判断是否满足限定渠道
//	if roleob.GetChannelID() != bean.ChannelID {
//		resp.Err = doError(kInviteCodeInvalid)
//		return nil, uint16(pbd.GC_ID_INVITE_CODE_RESP), resp
//	}
//
//	// 判断是否已被使用
//	if bean.Uid != 0 {
//		resp.Err = doError(kInviteCodeUsed)
//		return nil, uint16(pbd.GC_ID_INVITE_CODE_RESP), resp
//	}
//
//	// 判断是否过期
//	if bean.Deadline > 0 && time.Now().Unix() > bean.Deadline {
//		resp.Err = doError(kInviteCodeTimeout, fmt.Sprintf("%d", bean.Deadline))
//		return nil, uint16(pbd.GC_ID_INVITE_CODE_RESP), resp
//	}
//
//	bean.SetAchieved(roleob.Uid())
//	if err := obean.UpdateInviteCode(bean); err != nil {
//		log.Printf("更新邀请码[%d][%s]使用记录时出错: %v", roleob.Uid(), code, err)
//		resp.Err = doError(kInternelServerError)
//		return nil, uint16(pbd.GC_ID_INVITE_CODE_RESP), resp
//	}
//
//	roleob.SetInviteCode(code)
//
//	return nil, uint16(pbd.GC_ID_INVITE_CODE_RESP), resp
//}
//
//func HandleHeartBeat(skeleton *oskeleton.Skeleton, session *osession.Session, roleob *oaccount.Role, packet proto.Message) (error, uint16, proto.Message) {
//
//	return nil, uint16(pbd.GC_ID_HEART_BEAT_RESP), roleob.SerializeNotifies()
//}
