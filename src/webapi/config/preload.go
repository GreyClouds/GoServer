package config

import (
	"encoding/json"
	"log"
	"sort"
	"time"

	sd "crazyant.com/deadfat/data/go"
	"github.com/blang/semver"
	odisconf "yunjing.me/phoenix/go-phoenix/disconf"
)

func (m *Manager) Preload(manager *odisconf.Manager) {
	m.loadGlobals(manager)
	m.loadCharacter(manager)
	m.loadF1(manager)
	m.loadF3(manager)
	m.loadF2(manager)
	m.loadSigninTask(manager)
	m.loadMatchTask(manager)
	m.loadZoneList(manager)
	m.loadSKU(manager)

	m.loadInitRobots(manager)
	m.loadAItoScore(manager)
	m.loadScoretoAI(manager)
	m.loadRobotNames(manager)
	m.loadScoretoAISkin(manager)
	m.loadVersionUpdate(manager)
	m.loadActivity(manager)

	m.loadRank(manager)
	m.loadRankTiers(manager)
	m.loadRankGroup(manager)
	m.loadRankExtend(manager)
	m.loadRankSeasonConfig(manager)
	m.loadRankForbid(manager)
	m.loadRankCheatJudge(manager)
	m.loadCommodity(manager)
	m.loadChannel(manager)
	m.loadFYSDKIAP(manager)
}

// 加载全局杂项配置
func (m *Manager) loadGlobals(manager *odisconf.Manager) {
	payload := []*sd.GlobalOthers{}
	if err := manager.LoadJSONFile("GlobalOthers.json", &payload); err != nil {
		log.Panicln("加载杂项数据时出错", err)
	}

	for _, ele := range payload {
		m.others[ele.Key] = ele
	}
}

// 加载角色列表
func (self *Manager) loadCharacter(manager *odisconf.Manager) {
	payload := []*sd.Prop{}
	if err := manager.LoadJSONFile("Prop.json", &payload); err != nil {
		log.Panicln("加载角色数据时出错", err)
	}

	for _, ele := range payload {
		self.characters = append(self.characters, ele.ID)
	}
}

func (m *Manager) loadF1(manager *odisconf.Manager) {
	payload := []*sd.FunctionF1{}
	if err := manager.LoadJSONFile("FunctionF1.json", &payload); err != nil {
		log.Panicln("加载f1数据时出错", err)
		return
	}

	j := -1
	for i, v := range payload {
		if i == 0 || v.Score != payload[i-1].Score {
			m.functionF1 = append(m.functionF1, []*sd.FunctionF1{})
			j++
		}
		m.functionF1[j] = append(m.functionF1[j], v)
	}
}

// 加载活动配置
func (m *Manager) loadActivity(manager *odisconf.Manager) {
	payload := []*sd.Activity{}
	if err := manager.LoadJSONFile("Activity.json", &payload); err != nil {
		log.Panicln("加载活动数据时出错", err)
	}

	m.activitys = []*ActivityConf{}

	for _, ele := range payload {
		t1, err1 := time.Parse("20060102", ele.StartTime)
		if err1 != nil {
			log.Printf("解析活动启动时间%s时出错:%v", ele.StartTime, err1)
			continue
		}

		t2, err2 := time.Parse("20060102", ele.EndTime)
		if err2 != nil {
			log.Printf("解析活动结束时间%s时出错:%v", ele.EndTime, err2)
			continue
		}

		// 增加23小时59分59秒
		t1 = t1.Add(-8 * time.Hour)
		t2 = t2.Add(16*time.Hour - time.Second)

		// log.Printf("活动%d启动: %d ~ %d", ele.ID, t1.Unix(), t2.Unix())

		m.activitys = append(m.activitys, &ActivityConf{
			ID:      ele.ID,
			StartTS: t1.Unix(),
			EndTS:   t2.Unix(),
			Skins:   ele.DictValue,
		})
	}
}

func (m *Manager) loadF3(manager *odisconf.Manager) {
	payload := []*sd.FunctionF3{}
	if err := manager.LoadJSONFile("FunctionF3.json", &payload); err != nil {
		log.Panicln("加载f3数据时出错", err)
		return
	}

	j := -1
	for i, v := range payload {
		if i == 0 || v.PlayerN != payload[i-1].PlayerN {
			m.functionF3 = append(m.functionF3, []*sd.FunctionF3{})
			j++
		}
		m.functionF3[j] = append(m.functionF3[j], v)
	}
}

func (m *Manager) loadF2(manager *odisconf.Manager) {
	payload := []*sd.FunctionF2{}
	if err := manager.LoadJSONFile("FunctionF2.json", &payload); err != nil {
		log.Panicln("加载f2数据时出错", err)
		return
	}

	j := -1
	for i, v := range payload {
		if i == 0 || v.Score != payload[i-1].Score {
			m.functionF2 = append(m.functionF2, []*sd.FunctionF2{})
			j++
		}
		m.functionF2[j] = append(m.functionF2[j], v)
	}
}

func (m *Manager) loadSigninTask(manager *odisconf.Manager) {
	payload := []*sd.SigninTask{}
	if err := manager.LoadJSONFile("SigninTask.json", &payload); err != nil {
		log.Panicln("加载签到任务数据时出错", err)
		return
	}

	m.signins = payload
}

func (m *Manager) loadMatchTask(manager *odisconf.Manager) {
	payload := []*sd.MatchTask{}
	if err := manager.LoadJSONFile("MatchTask.json", &payload); err != nil {
		log.Panicln("加载胜场任务数据时出错", err)
		return
	}

	m.matchTasks = payload
}

func (m *Manager) loadZoneList(manager *odisconf.Manager) {
	payload := []*sd.GlobalZones{}
	if err := manager.LoadJSONFile("GlobalZones.json", &payload); err != nil {
		log.Panicln("加载区域列表数据时出错", err)
		return
	}

	m.zones = payload
}

func (m *Manager) loadSKU(manager *odisconf.Manager) {
	payload := []*sd.SKU{}
	if err := manager.LoadJSONFile("SKU.json", &payload); err != nil {
		log.Panicln("加载付费编号配置时出错", err)
		return
	}

	m.skus = make([]map[string]*sd.SKU, 1)

	appstores := make(map[string]*sd.SKU)
	for _, v := range payload {
		appstores[v.AppStore] = v
	}
	m.skus[0] = appstores
}

func (m *Manager) loadInitRobots(manager *odisconf.Manager) {
	payload := []*sd.InitRobotLvs{}
	if err := manager.LoadJSONFile("InitRobotLvs.json", &payload); err != nil {
		log.Panicln("加载新人机器人等级列表数据时出错", err)
		return
	}

	m.initRobotLvs = payload
}

func (m *Manager) loadAItoScore(manager *odisconf.Manager) {
	payload := []*sd.AItoScore{}
	if err := manager.LoadJSONFile("AItoScore.json", &payload); err != nil {
		log.Panicln("加载ai查询分数列表数据时出错", err)
		return
	}

	m.aiToScore = payload
}

func (m *Manager) loadScoretoAI(manager *odisconf.Manager) {
	payload := []*sd.ScoretoAI{}
	if err := manager.LoadJSONFile("ScoretoAI.json", &payload); err != nil {
		log.Panicln("加载分数查询ai列表数据时出错", err)
		return
	}

	m.scoreToAI = payload
}

func (m *Manager) loadRobotNames(manager *odisconf.Manager) {
	payload := []*sd.RobotName{}
	if err := manager.LoadJSONFile("RobotName.json", &payload); err != nil {
		log.Panicln("加载ai名字数据时出错", err)
		return
	}

	m.robotNames = payload
}

func (m *Manager) loadScoretoAISkin(manager *odisconf.Manager) {
	payload := []*sd.ScoretoAISkin{}
	if err := manager.LoadJSONFile("ScoretoAISkin.json", &payload); err != nil {
		log.Panicln("加载分数查询ai皮肤列表数据时出错", err)
		return
	}

	for skin, account := range payload[0].RobotSkinSolu {
		for i := int32(0); i < account; i++ {
			m.scoreToAISkin = append(m.scoreToAISkin, skin)
		}
	}
}

func (m *Manager) loadVersionUpdate(manager *odisconf.Manager) {
	payload := []*sd.VersionUpdateAlert{}
	if err := manager.LoadJSONFile("VersionUpdateAlert.json", &payload); err != nil {
		log.Panicln("加载版本更新提醒列表数据时出错", err)
		return
	}

	m.versions = payload

	for _, v := range payload {
		version, err := semver.Make(v.Version)
		if err != nil {
			log.Panicln("渠道%s的版本号%s不合理", v.Channel, v.Version, err)
			continue
		}

		if oldVersion, exists := m.newVersions[v.Channel]; exists {
			if oldVersion.LT(version) {
				m.newVersions[v.Channel] = version
			}
		} else {
			m.newVersions[v.Channel] = version
		}
	}
}

func (m *Manager) loadRank(conf *odisconf.Manager) {
	rank := []*sd.Rank{}
	if err := conf.LoadJSONFile("Rank.json", &rank); err != nil {
		log.Panicln("加载Rank.json数据时出错", err)
	}
	data, err := json.Marshal(rank)
	if nil != err {
		log.Panicln("还原Rank.json数据时出错", err)
	}
	m.tiersTable = data
	l := len(rank)
	for i := 0; i < l; i++ {
		key := (rank[i].ArenaLevel << 16) | (rank[i].SubLevel << 8)
		m.rank[key] = rank[i]
	}
}

func (m *Manager) loadRankTiers(c *odisconf.Manager) {
	tiers := []*sd.RankTiers{}
	err := c.LoadJSONFile("RankTiers.json", &tiers)
	if nil != err {
		log.Panicln("RankTiers.json数据时出错", err)
	}
	sort.Slice(tiers, func(i int, j int) bool {
		return getConfigTiersLv(tiers[i]) < getConfigTiersLv(tiers[j])
	})
	m.createRankTiers(tiers)
}

func getConfigTiersLv(t *sd.RankTiers) int32 {
	return t.ArenaLevel<<16 | t.SubLevel<<8 | t.LevelStars
}

func (m *Manager) createRankTiers(t []*sd.RankTiers) {
	l := len(t)
	for i := 0; i < l; i++ {
		node := getNode(m.rankTiers, t[i])
		if i == 0 {
			node.Down = node
		} else if m.IsRankTiersGuard(getConfigTiersLv(t[i])) {
			node.Down = node
		} else {
			if t[i].LevelStars == 0 {
				node.Down = getNode(m.rankTiers, t[i-2])
			} else {
				node.Down = getNode(m.rankTiers, t[i-1])
			}
		}
		if i == l-1 {
			node.Up = node
		} else if t[i+1].LevelStars == 0 {
			node.Up = getNode(m.rankTiers, t[i+2])
		} else {
			node.Up = getNode(m.rankTiers, t[i+1])
		}
	}
}

func makeTiers(t *sd.RankTiers) *Tiers {
	return &Tiers{
		Lv:    getConfigTiersLv(t),
		Score: t.Score,
	}
}

func getNode(n map[int32]*Tiers, t *sd.RankTiers) *Tiers {
	key := getConfigTiersLv(t)
	v, ok := n[key]
	if ok {
		return v
	} else {
		node := makeTiers(t)
		n[key] = node
		return node
	}
}

func (m *Manager) loadRankGroup(c *odisconf.Manager) {
	group := []*sd.RankGroup{}
	if err := c.LoadJSONFile("RankGroup.json", &group); err != nil {
		log.Panicln("RankGroup.json数据时出错", err)
	}

	sort.Slice(group, func(i int, j int) bool {
		if group[i].PlayerN == group[j].PlayerN {
			return group[i].ScoreRange < group[j].ScoreRange
		} else {
			return group[i].PlayerN < group[j].PlayerN
		}
	})

	l := len(group)
	for i := 0; i < l; i++ {
		n := len(m.rankGroup)
		if n == 0 {
			m.rankGroup = append(m.rankGroup, []*sd.RankGroup{group[i]})
		} else {
			if group[i].PlayerN == m.rankGroup[n-1][0].PlayerN {
				m.rankGroup[n-1] = append(m.rankGroup[n-1], group[i])
			} else {
				m.rankGroup = append(m.rankGroup, []*sd.RankGroup{group[i]})
			}
		}
	}
}

func (m *Manager) GetRankGroup(num int32, diff int32) int32 {
	l := len(m.rankGroup)
	if l == 0 {
		return 1
	}
	for i := 0; i < l; i++ {
		if m.rankGroup[i][0].PlayerN == num {
			return findRankGroup(m.rankGroup[i], diff)
		} else if m.rankGroup[i][0].PlayerN > num {
			if i == 0 {
				return findRankGroup(m.rankGroup[i], diff)
			} else {
				return findRankGroup(m.rankGroup[i-1], diff)
			}
		} else if i == l-1 {
			return findRankGroup(m.rankGroup[i], diff)
		}

	}
	return 6
}

func findRankGroup(group []*sd.RankGroup, diff int32) int32 {
	l := len(group)
	for i := 0; i < l; i++ {
		if group[i].ScoreRange > diff {
			if i == 0 {
				return group[i].GroupK
			} else {
				return group[i-1].GroupK
			}
		} else if group[i].ScoreRange == diff {
			return group[i].GroupK
		}
	}
	return group[l-1].GroupK
}

func (m *Manager) loadRankExtend(c *odisconf.Manager) {
	stdExtend := []*sd.RankMatchExtend{}
	if err := c.LoadJSONFile("RankMatchExtend.json", &stdExtend); err != nil {
		log.Panicln("RankMatchExtend.json数据时出错", err)
	}
	sort.Slice(stdExtend, func(i int, j int) bool {
		if stdExtend[i].Score == stdExtend[j].Score {
			return stdExtend[i].WaitTime < stdExtend[j].WaitTime
		} else {
			return stdExtend[i].Score < stdExtend[j].Score
		}
	})
	l := len(stdExtend)
	var node *RankExtend
	for i := 0; i < l; i++ {
		temp := newRankExtend(stdExtend[i])
		if node != nil && node.score == stdExtend[i].Score {
			node.next = temp
		} else {
			m.rankExtend = append(m.rankExtend, temp)
		}
		node = temp
	}
}

func newRankExtend(r *sd.RankMatchExtend) *RankExtend {
	return &RankExtend{
		score:      r.Score,
		waitTimes:  r.WaitTime,
		scoreRange: r.ScoreRange,
	}
}

func (m *Manager) loadRankSeasonConfig(c *odisconf.Manager) {
	seasonConf := []*sd.RankSeasonConfig{}
	err := c.LoadJSONFile("RankSeasonConfig.json", &seasonConf)
	if nil != err {
		log.Panicln("RankSeasonConfig.json数据时出错", err)
	}
	l := len(seasonConf)
	if l == 0 {
		log.Panicln("警告! 未获取赛季配置,请检查配置表")
		return
	}
	m.rankSeasonConf = make([]*RankSeasonConfig, l)
	for i := 0; i < l; i++ {
		m.rankSeasonConf[i] = newRankSeasonConf(seasonConf[i])
	}
	checkSeasonConfig(m.rankSeasonConf)
}

func newRankSeasonConf(r *sd.RankSeasonConfig) *RankSeasonConfig {
	start, err := time.Parse(bTimeFmt, r.StartTime)
	if nil != err {
		log.Panicln("错误 解析赛季配置表.开始时间", r)
	}
	end, err2 := time.Parse(bTimeFmt, r.EndTime)
	if nil != err2 {
		log.Panicln("错误 解析赛季配置表.结算时间", r)
	}
	return &RankSeasonConfig{
		r.SeasonID,
		start.Unix(),
		end.Unix(),
		int64(r.RankModeLockTime),
		newRewardList(r.Rank1Reward),
		newRewardList(r.Rank2Reward),
		newRewardList(r.Rank3Reward),
		newRewardList(r.Rank4Reward),
		newRewardList(r.Rank5Reward),
		newRewardList(r.Rank6Reward),
	}
}

func newRewardList(r map[int32]int32) []SeasonReward {
	l := len(r)
	reward := make([]SeasonReward, 0, l)
	for k, v := range r {
		reward = append(reward, SeasonReward{k, v})
	}
	return reward
}

func checkSeasonConfig(c []*RankSeasonConfig) {
	l := len(c)
	for i := 0; i < l; i++ {
		if c[i].Start <= 0 {
			log.Panicln("错误 解析赛季配置开始时间非法", c[i])
			return
		}
		if c[i].End <= 0 {
			log.Panicln("错误 解析赛季配置结算时间非法", c[i])
			return
		}
		if c[i].Start >= c[i].End-c[i].Lock {
			log.Panicln("错误 解析赛季配置开始时间超过结算时间", c[i])
			return
		} else {
			if i > 0 {
				if c[i].Start < c[i-1].End {
					log.Panicln("错误 解析赛季配置开始时间超过结算时间", c[i])
				}
			}
		}
	}
}

func (m *Manager) loadRankForbid(c *odisconf.Manager) {
	forbid := []*sd.RankForbid{}
	err := c.LoadJSONFile("RankForbid.json", &forbid)
	if nil != err {
		log.Panicln("RankForbid.json数据时出错", err)
	}
	sort.Slice(forbid, func(i, j int) bool {
		return forbid[i].ForbidID < forbid[j].ForbidID
	})
	m.rankForbid = forbid
}

func (m *Manager) loadRankCheatJudge(c *odisconf.Manager) {
	cheat := []*sd.RankCheatJudge{}
	if err := c.LoadJSONFile("RankCheatJudge.json", &cheat); err != nil {
		log.Panicln("RankCheatJudge.json数据时出错", err)
	}

	sort.Slice(cheat, func(i int, j int) bool {
		if cheat[i].All == cheat[j].All {
			return cheat[i].Suspicious < cheat[j].Suspicious
		} else {
			return cheat[i].All < cheat[j].All
		}
	})

	l := len(cheat)
	for i := 0; i < l; i++ {
		n := len(m.rankCheatJudge)
		if n == 0 {
			m.rankCheatJudge = append(m.rankCheatJudge, []*sd.RankCheatJudge{cheat[i]})
		} else {
			if cheat[i].All == m.rankCheatJudge[n-1][0].All {
				m.rankCheatJudge[n-1] = append(m.rankCheatJudge[n-1], cheat[i])
			} else {
				m.rankCheatJudge = append(m.rankCheatJudge, []*sd.RankCheatJudge{cheat[i]})
			}
		}
	}
}

func (m *Manager) loadCommodity(c *odisconf.Manager) {
	commodity := []*sd.Commodity{}
	err := c.LoadJSONFile("Commodity.json", &commodity)
	if nil != err {
		log.Panicln("Commodity.json数据时出错", err)
	}
	for _, v := range commodity {
		m.commodity[v.SellID] = v
	}
}

func (self *Manager) loadChannel(manager *odisconf.Manager) {
	payload := []*sd.Channel{}
	if err := manager.LoadJSONFile("Channel.json", &payload); err != nil {
		log.Panicln("加载渠道配置时出错", err)
		return
	}

	self.channels = make(map[string]*sd.Channel)

	for _, v := range payload {
		self.channels[v.ChannelID] = v
	}
}

func (self *Manager) loadFYSDKIAP(manager *odisconf.Manager) {
	payload := []*sd.AndroidIAP{}
	if err := manager.LoadJSONFile("AndroidIAP.json", &payload); err != nil {
		log.Panicln("加载FYSDK内购配置时出错", err)
		return
	}

	self.fysdkIAPs = make(map[string]*sd.AndroidIAP)

	for _, v := range payload {
		self.fysdkIAPs[v.ProductID] = v
	}
}
