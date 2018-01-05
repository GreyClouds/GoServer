package bean

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/astaxie/beego/orm"
)

const (
	GUIDE_NICK_SET = int32(0) // 新手引导: 昵称设置
	GUIDE_CONTINUE = int32(1) // 新手引导：待续

	kInitialScore = int32(5000)
)

type Character struct {
	ID                uint32 `orm:"pk;column(id)"`                        // 角色编号
	Nick              string `orm:"column(nick)"`                         // 昵称(限制14个字符)
	Created           int64  `orm:"column(created)"`                      // 创建时间
	Updated           int64  `orm:"column(updated)"`                      // 上次更新时间
	Modified          int64  `orm:"column(modified)"`                     // 修改次数
	TotalTimes        int32  `orm:"column(total)"`                        // 游戏总次数
	WinedTimes        int32  `orm:"column(wined)"`                        // 游戏胜利次数
	LosedTimes        int32  `orm:"column(losed)"`                        // 游戏失败次数
	ContiLoses        int32  `orm:"column(conti_losed)"`                  // 连败场次
	ContiAITimes      int32  `orm:"column(conti_ai_times)"`               // 连续ai场次
	TestBattles       int32  `orm:"column(test_battles)"`                 // 新手测试场次
	TestBattleID      int32  `orm:"column(test_battle_id)"`               // 新手测试比赛id
	GuideID           int32  `orm:"column(guide)"`                        // 引导编号
	LastSigninTime    int64  `orm:"column(last_signin_t)"`                // 上次签到时间
	LastSigninID      int32  `orm:"column(last_signin_id)"`               // 上次签到编号
	Score             int32  `orm:"column(score)"`                        // 分数值
	RemovedAD         bool   `orm:"column(remove_ad)"`                    // 是否去除广告
	LastBattleSkin    int32  `orm:"column(last_battle_skin)"`             // 最近战斗使用的皮肤
	LastLoginTime     int64  `orm:"column(last_login_time)"`              // 最近一次登录时间
	LastChannelID     string `orm:"column(last_channel);size(6)"`         // 最近一次登录来源渠道
	LastClientVersion string `orm:"column(last_client_version);size(32)"` // 最近一次登录来源客户端版本
	FirstRecharge     bool   `orm:"column(first_recharge)"`               // 是否首充状态
	InviteCode        string `orm:"column(invite_code);size(10)"`         // 邀请码
	dirty             bool   `orm:"-"`                                    // 是否脏数据
}

func (self Character) TableName() string {
	return "character"
}

func NewCharacter(uid uint32) (error, *Character) {
	character := &Character{
		ID:         uid,
		Nick:       fmt.Sprintf("Eater%d", 1000+rand.Intn(9000)),
		Created:    time.Now().Unix(),
		Updated:    time.Now().Unix(),
		Modified:   0,
		TotalTimes: 0,
		WinedTimes: 0,
		LosedTimes: 0,
		ContiLoses: 0,
		Score:      kInitialScore,
		GuideID:    GUIDE_NICK_SET,
	}

	o := orm.NewOrm()
	_, err := o.Insert(character)
	if err != nil {
		return err, nil
	}

	return nil, character
}

func (self *Character) SetDirty() {
	self.dirty = true
}

func (self *Character) SetNick(nick string) {
	self.Nick = nick
	self.dirty = true
}

func (self *Character) SetGuideID(id int32) {
	self.GuideID = id
	self.dirty = true
}

func (self *Character) Signin(id int32) {
	self.LastSigninID = id
	self.LastSigninTime = time.Now().Unix()
	self.dirty = true
}

func (self *Character) ScoreChange(scoreDelta int32) {
	score := self.Score + scoreDelta
	if score < 0 {
		score = 0
	}
	self.Score = score
	self.dirty = true
}

func (self *Character) SetScore(score int32) {
	self.Score = score
	self.dirty = true
}

func (self *Character) OneAIMatch() {
	self.ContiAITimes += 1
	self.dirty = true
}

func (self *Character) SetLastBattleSkin(lastBattleSkin int32) {
	self.LastBattleSkin = lastBattleSkin
	self.dirty = true
}

func (self *Character) OneNonAIMatch() {
	self.ContiAITimes = 0
	self.dirty = true
}

func (self *Character) WinMatch() {
	self.TotalTimes += 1
	self.WinedTimes += 1
	self.ContiLoses = 0
	self.dirty = true
}

func (self *Character) LoseMatch() {
	self.TotalTimes += 1
	self.LosedTimes += 1
	self.ContiLoses += 1
	self.dirty = true
}

func (self *Character) SetTestID(id int32) {
	self.TestBattleID = id
	self.dirty = true
}

func (self *Character) TestBattleAddOne() {
	self.TestBattles += 1
	self.dirty = true
}

func (self *Character) DrawMatch() {
	self.TotalTimes += 1
	self.dirty = true
}

func (self *Character) ResetDirty() bool {
	if !self.dirty {
		return false
	}

	self.dirty = false

	return true
}

// ---------------------------------------------------------------------------

func LoadCharacter(uid uint32) (error, *Character) {
	o := orm.NewOrm()
	character := Character{
		ID: uid,
	}

	err := o.Read(&character)
	if err != nil {
		// 未找到
		if err == orm.ErrNoRows {
			return err, nil
		}

		return err, nil
	}

	return nil, &character
}

func UpdateCharacter(bean *Character) error {
	o := orm.NewOrm()
	_, e := o.Update(bean)
	if e != nil {
		return e
	}
	return nil
}
