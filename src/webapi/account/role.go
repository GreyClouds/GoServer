package account

import (
	"log"
	"math"
	"sync"
	"time"

	pkgActivity "webapi/activity"
	pkgBean "webapi/bean"
	pkgConfig "webapi/config"

	pkgProto "crazyant.com/deadfat/pbd/go"
	"github.com/golang/protobuf/proto"
)

const winRewardPerNum = 3

// 角色
type Role struct {
	character      *pkgBean.Character
	skins          map[int32]*pkgBean.Skin
	skinTasks      map[int32]*pkgBean.SkinTask
	zone           int32 // 区域编号
	rank           *pkgBean.Rank
	money          *pkgBean.Money
	winReward      *pkgBean.WinReward
	adReward       *pkgBean.AdReward
	notifies       []*pkgBean.Notify
	l              sync.Mutex
	lastActionTime int64
	ID                uint32
}

func NewRole(character *pkgBean.Character, skins []*pkgBean.Skin, skinTasks []*pkgBean.SkinTask, notifies []*pkgBean.Notify) *Role {
	self := &Role{
		character: character,
		skins:     make(map[int32]*pkgBean.Skin),
		skinTasks: make(map[int32]*pkgBean.SkinTask),
		notifies:  []*pkgBean.Notify{},
	}

	if skins != nil && len(skins) > 0 {
		for _, skin := range skins {
			self.skins[skin.SkinID] = skin
		}
	}

	if skinTasks != nil && len(skinTasks) > 0 {
		for _, skinTask := range skinTasks {
			self.skinTasks[skinTask.SkinID] = skinTask
		}
	}

	if notifies != nil && len(notifies) > 0 {
		self.notifies = append(self.notifies, notifies...)
	}

	return self
}

func (self *Role) GetLastActionTime() int64 {
	return self.lastActionTime
}

func (self *Role) UpdateLastActionTime() *Role {
	self.lastActionTime = time.Now().UnixNano()
	return self
}

func (self *Role) Lock() {
	self.l.Lock()
}

func (self *Role) Unlock() {
	self.l.Unlock()
}

// 角色编号
func (self Role) Uid() uint32 {
	return self.ID
}

func (self *Role) GetGuideID() int32 {
	return self.character.GuideID
}

func (self Role) Score() int32 {
	return self.character.Score
}

func (self *Role) GetRank() *pkgBean.Rank {
	return self.rank
}

func (self *Role) SetRank(r *pkgBean.Rank) *Role {
	self.rank = r
	return self
}

func (self *Role) GetRankTiers() int32 {
	return self.rank.Tiers
}

func (self *Role) GetRankScore() int32 {
	t, ok := pkgConfig.Conf().GetRankTiers(self.rank.Tiers)
	if ok {
		return t.Score
	} else {
		return 0
	}
}

func (self *Role) RankElo() float64 {
	//rankScore := float64(self.rank.score)
	rankScore := 0.0
	normalScore := float64(self.character.Score)
	score := rankScore*rankScore*1.5 + normalScore*normalScore*0.005
	return math.Sqrt(score)
}

func (self Role) WinedTimes() int32 {
	return self.character.WinedTimes
}

func (self Role) TotalBattleTimes() int32 {
	return self.character.TotalTimes
}

func (self Role) TestBattleTimesAndID() (int32, int32) {
	return self.character.TestBattles, self.character.TestBattleID
}

func (self Role) TestBattleID() int32 {
	return self.character.TestBattleID
}

func (self *Role) SetTestID(level int32) {
	self.character.SetTestID(level)
}

func (self Role) ContiLoseTimes() int32 {
	return self.character.ContiLoses
}

func (self Role) ContiAITimes() int32 {
	return self.character.ContiAITimes
}

func (self Role) NickName() string {
	return self.character.Nick
}

func (self *Role) AddOneAIDispatchRecord() {
	self.character.OneAIMatch()
}

func (self *Role) SetLastBattleSkin(skinId int32) {
	self.character.SetLastBattleSkin(skinId)
}

func (self *Role) LastBattleSkin() int32 {
	return self.character.LastBattleSkin
}

func (self *Role) ClearAIDispatchRecord() {
	self.character.OneNonAIMatch()
}

func (self *Role) TestBattleAddOne() {
	self.character.TestBattleAddOne()
}

// 设置昵称并设置引导编号
func (self *Role) SetNickAndGuide(nick string) {
	self.character.SetNick(nick)
	self.character.SetGuideID(pkgBean.GUIDE_CONTINUE)
}

func (self *Role) WinMatch(scoreDelta int32) {
	self.character.WinMatch()
	self.ClearWinReward()
	self.character.ScoreChange(scoreDelta)
}

func (self *Role) LoseMatch(scoreDelta int32) {
	self.character.LoseMatch()
	self.character.ScoreChange(scoreDelta)
}

func (self *Role) SetScore(score int32) {
	self.character.SetScore(score)
}

func (self *Role) DrawMatch() {
	self.character.DrawMatch()
}

func (self *Role) GetLastSigninTask() (int64, int32) {

	return self.character.LastSigninTime, self.character.LastSigninID
}

func (self *Role) Signin(id int32) {
	self.character.Signin(id)
}

func (self *Role) AddSkin(id int32, src int) bool {
	if _, exists := self.skins[id]; exists {
		return false
	}

	err, skins := pkgBean.NewSkins(self.character.ID, []int32{id}, src)
	if err != nil {
		log.Printf("角色%d获得皮肤%d时出错: %v", self.character.ID, id, err)
		return false
	}

	if skins == nil || len(skins) == 0 {
		log.Printf("角色%d获得皮肤%d时失败", self.character.ID, id)
		return false
	}

	for _, skin := range skins {
		self.skins[skin.SkinID] = skin
	}

	return true
}

func (self *Role) SetZone(zone int32) {
	self.zone = zone
}

func (self Role) Zone() int32 {
	return self.zone
}

func (self *Role) RemoveAD() {
	if self.character.RemovedAD {
		return
	}

	self.character.RemovedAD = true
	self.character.SetDirty()
}

func (self Role) IsRemoveAD() bool {
	return self.character.RemovedAD
}

func (self Role) DispatchData() (uint32, int32, int32, int32, string, int32) {
	return self.character.ID, self.character.WinedTimes, self.character.Score, self.zone, self.character.Nick, self.character.ContiLoses
}

func (self Role) SkinOwned(skin int32) bool {
	if _, exist := self.skins[skin]; !exist {
		return false
	}

	return true
}

func (self *Role) GetSkinTask(id int32) int32 {
	v, exists := self.skinTasks[id]
	if exists {
		return v.WinNum
	}

	return 0
}

// 皮肤任务进度累计
func (self *Role) GainSkinTask(id int32) int32 {
	_, exists := self.skinTasks[id]
	if exists {
		self.skinTasks[id].WinNum += 1
		self.skinTasks[id].SetDirty()
		return self.skinTasks[id].WinNum
	}

	err, skinTask := pkgBean.NewSkinTask(self.character.ID, id)
	if err != nil {
		log.Printf("角色%d活动任务皮肤%d时出错: %v", self.character.ID, id, err)
		return 0
	}

	if skinTask == nil {
		log.Printf("角色%d活动任务皮肤%d时出错: %v", self.character.ID, id)
		return 0
	}

	self.skinTasks[skinTask.SkinID] = skinTask

	return skinTask.WinNum
}

// 设置激活码
func (self *Role) SetInviteCode(code string) {
	self.character.InviteCode = code
	self.character.SetDirty()
}

func (self *Role) ExistsInviteCode() bool {
	return self.character.InviteCode != ""
}

func (self *Role) AddNotify(notify *pkgBean.Notify) {
	self.notifies = append(self.notifies, notify)
}

// 数据回存
func (self *Role) Save() {
	if dirty := self.character.ResetDirty(); dirty {
		pkgBean.UpdateCharacter(self.character)
	}

	for _, skinTask := range self.skinTasks {
		if dirty := skinTask.ResetDirty(); dirty {
			pkgBean.UpdateSkinTask(skinTask)
		}
	}

	for i := 0; i < len(self.notifies); i++ {
		notify := self.notifies[i]
		if dirty := notify.ResetDirty(); dirty {
			pkgBean.UpdateNotify(notify)
		}
	}
}

// 序列化皮肤任务进度
func (self *Role) SerializeSkinTask(skins []int32) []*pkgProto.SkinWinData {
	skinWinDatas := []*pkgProto.SkinWinData{}
	for _, id := range skins {
		skin := &pkgProto.SkinWinData{
			SkinId: proto.Int32(id),
			WinNum: proto.Int32(0),
		}
		v, exists := self.skinTasks[id]
		if exists {
			skin.WinNum = proto.Int32(v.WinNum)
		}
		skinWinDatas = append(skinWinDatas, skin)
	}
	return skinWinDatas
}

// 序列化财富
func (self *Role) SerializeMoney() *pkgProto.Money {
	return &pkgProto.Money{
		Gold:          proto.Int32(self.money.Gold),
		Diamond:       proto.Int32(self.money.Diamond),
		FirstRecharge: proto.Bool(!self.character.FirstRecharge),
	}
}

func (self *Role) Serialize(conf *pkgConfig.Manager) *pkgProto.Profile {
	skins := []int32{}
	for k := range self.skins {
		skins = append(skins, k)
	}

	winRewardNum, winRewardCount := self.Get3WinRewardNum()

	profile := &pkgProto.Profile{
		Nick:              proto.String(self.character.Nick),
		Skins:             skins,
		Newbie:            proto.Bool(self.character.GuideID == pkgBean.GUIDE_NICK_SET),
		Win:               proto.Int32(self.character.WinedTimes),
		RemoveAd:          proto.Bool(self.character.RemovedAD),
		AdReward:          proto.Bool(self.adReward.Num < pkgConfig.Conf().GetMaxRewardedAdsPerDay()),
		RemainAdRewardSec: proto.Int32(int32(self.adReward.CD - time.Now().Unix())),
		WinRewardNum:      proto.Int32(winRewardNum),
		WinRewardCount:    proto.Int32(winRewardCount),
	}

	if activities := pkgActivity.DefaultActivityManager.Serialize(self); activities != nil {
		profile.Activity = activities
	}

	return profile
}

func (self *Role) SerializeNotifies() *pkgProto.HeartBeat {
	resp := &pkgProto.HeartBeat{}

	if n := len(self.notifies); n > 0 {
		for i := 0; i < n; i++ {
			v := self.notifies[i]

			// 过滤已经删除的通知
			if v.Removed {
				continue
			}

			switch v.Category {
			case pkgBean.NOTIFY_PAYMENT:
				item := &pkgProto.PaymentNotify{}
				if err := v.GetRaw(item); err == nil {
					v.Removed = true
					v.SetDirty()

					resp.Payments = append(resp.Payments, item)
				}
			}
		}
	}

	return resp
}

func (self *Role) SetLoginInfo(channel, version string) {
	self.character.LastLoginTime = time.Now().Unix()
	self.character.LastClientVersion = version
	self.character.LastChannelID = channel
	self.character.SetDirty()
}

// 未触发Login流程的，全部默认当成旧版本处理
// 仅限于当重启服务器的期间，客户端未关闭过的情况
func (self *Role) IsAiClientVersion() bool {

	return self.character.LastClientVersion != ""
}

func (self *Role) IsQiXiVersion() bool {
	v := self.character.LastClientVersion

	return v == "" || v == "1.5.0" || v == "1.5.1"
}

func (self *Role) GoldAdd(i int32) *Role {
	self.money.ChangeGold(i)
	return self
}

func (self *Role) DiamondAdd(i int32, src int) *Role {
	self.money.ChangeDiamond(i, src)
	return self
}

func (self *Role) GoldAndDiamondAdd(g int32, d int32, src int) *Role {
	if g != 0 {
		self.money.ChangeGold(g)
	}
	if d != 0 {
		self.money.ChangeDiamond(d, src)
	}
	return self
}

func (self *Role) GetMoney() *pkgBean.Money {
	return self.money
}

func (self *Role) ClearWinReward() {
	max := pkgConfig.Conf().GetMaxTriWinReward() * winRewardPerNum
	if self.winReward.Win < max {
		self.winReward.Win++
		if self.winReward.Win%winRewardPerNum == 0 {
			gold := pkgConfig.Conf().GetGoldEveryTriWin()
			self.GoldAdd(gold)
		}
		self.winReward.Update()
	}
}

// 当前三胜回合内的胜场数和剩余次数
func (self *Role) Get3WinRewardNum() (int32, int32) {
	max := pkgConfig.Conf().GetMaxTriWinReward() * winRewardPerNum
	if self.winReward.Win < max {
		return self.winReward.Win, int32((max - self.winReward.Win + winRewardPerNum - 1) / winRewardPerNum)
	}

	return 0, 0
}

func (self *Role) GetAdReward() *pkgBean.AdReward {
	return self.adReward
}

func (self *Role) GetWinReward() *pkgBean.WinReward {
	return self.winReward
}

// 获取上次登录时使用的渠道编号
func (self *Role) GetChannelID() string {
	return self.character.LastChannelID
}

// 是否满足首充条件
func (self *Role) IsFirstRecharge() bool {
	return !self.character.FirstRecharge
}

// 清空首充福利标记
func (self *Role) CleanFirstRechargeFlag() {
	self.character.FirstRecharge = true
	self.character.SetDirty()
}
