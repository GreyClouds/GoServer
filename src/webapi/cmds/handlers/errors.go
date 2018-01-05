package handlers

import (

	pbd "crazyant.com/deadfat/pbd/go"
	pbdHero "crazyant.com/deadfat/pbd/hero"
)

type error_t int32

const (
	kInternelServerError = error_t(1) // 服务器内部错误
	kRoleNotExists       = error_t(2) // 角色不存在
	kBadRequest          = error_t(4) // 客户端请求参数不合理
	kFriendMatchIgnore   = error_t(5) // 好友对战忽略比赛
	kAlreadyRegistered   = error_t(6) // 账户已经注册过了


	kNickNameHaveSet        = error_t(100)  // 昵称已经设置过
	kIMEIIsEmpty            = error_t(101)  // 设备号为空
	kRoleTodayAlreadySignin = error_t(102)  // 角色今日已签到
	kNickIsSensitivity      = error_t(103)  // 昵称包含敏感词
	kPaymentAlreadyExists   = error_t(104)  // 支付订单已存在
	kSkinNotOwn             = error_t(1000) // 未获得皮肤

	kCDKeyInvalid         = error_t(1100) // 兑换码不合法
	kCDKeyNotExists       = error_t(1101) // 兑换码已过期或者不存在
	kCDKeyIsAchieved      = error_t(1102) // 兑换码已经被使用过
	kCDKeyChannelNotMatch = error_t(1103) // 兑换码限定渠道不匹配
	kCDKeyIsTimeout       = error_t(1104) // 兑换码已过期[1:过期时间]
	kCDKeyGiftAlreadyUse  = error_t(1105) // 兑换码同一礼包已兑换过

	kInviteCodeInvalid = error_t(1120) // 邀请码不合法
	kInviteCodeUsed    = error_t(1121) // 邀请码已被使用
	kInviteCodeTimeout = error_t(1122) // 邀请码已过期[1:过期时间]
)

func doError(code error_t, args ...string) *pbd.Error {
	err := &pbd.Error{
		Code: int32(code),
	}

	if args != nil && len(args) > 0 {
		err.Args = args
	}

	return err
}

func doHeroError(code error_t, args ...string) *pbdHero.Error {
	err := &pbdHero.Error{
		Code: int32(code),
	}

	if args != nil && len(args) > 0 {
		err.Args = args
	}

	return err
}
