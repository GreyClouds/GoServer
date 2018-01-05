package main

import (
	"math/rand"
	"time"

	pbd "crazyant.com/deadfat/pbd/hero"

	ohandlers "webapi/cmds/handlers"
	oskeleton "webapi/skeleton"

	"webapi/fysdk"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	skeleton := oskeleton.New()

	// APP ID 	10040
	// AppKey 	ae35c607096d1b35a4ff66111085d078
	// 充值Key 	7ac57d8178fe13ed07d42f37bf747adc
	fysdk.Initialize("10040", "ae35c607096d1b35a4ff66111085d078", "7ac57d8178fe13ed07d42f37bf747adc")

	skeleton.Register(false, uint16(pbd.CG_ID_GUEST_LOGIN), &pbd.Login{}, ohandlers.HandleGuestLogin)
	skeleton.Register(false, uint16(pbd.CG_ID_GUEST_REGISTER), &pbd.GuestRegister{}, ohandlers.HandleRegister)
	skeleton.Register(true, uint16(pbd.CG_ID_GET_ACHIEVEMENT), &pbd.GetAchievement{}, ohandlers.HandleGetAchievement)
	skeleton.Register(true, uint16(pbd.CG_ID_SET_ACHIEVEMENT), &pbd.SetAchievement{}, ohandlers.HandleSetAchievement)
	skeleton.Register(true, uint16(pbd.CG_ID_GET_LEARDERBOARD_RANGE), &pbd.GetLearderboardRange{}, ohandlers.HandleGetLeaboardRange)
	skeleton.Register(true, uint16(pbd.CG_ID_UPDATE_LEARD_SCORE), &pbd.UpdateLearderBoardScore{}, ohandlers.HandleLearderBoardUpdae)
	skeleton.Register(true, uint16(pbd.CG_ID_UNLOCK_ACHIEVEMENT), &pbd.UnLockAchievement{}, ohandlers.HandleUnLockAchievement)
	skeleton.Register(true, uint16(pbd.CG_ID_GET_UNLOCK_ACHIEVEDATE), &pbd.GetUnLockAchieveDate{}, ohandlers.GetUnLockAchievement)
	//skeleton.Register(true, uint16(pbd.CG_ID_CDKEY_EXCHANGE), &pbd.CDKey{}, ohandlers.HandleCDKeyExchange)
	//skeleton.Register(true, uint16(pbd.CG_ID_SELF_MONEY), nil, ohandlers.HandleSelfMoney)
	//skeleton.Register(true, uint16(pbd.CG_ID_SHOPPING), &pbd.Shopping{}, ohandlers.HandleShopping)
	//skeleton.Register(true, uint16(pbd.CG_ID_AD_REWARD), nil, ohandlers.HandleAdReward)
	skeleton.RegisterDBModel()
	skeleton.Serve()

}
