package config

import (
	"log"
	"math/rand"
	"time"

	"github.com/blang/semver"

	sd "crazyant.com/deadfat/data/go"
)

const (
	bTimeFmt                = "20060102T15:04:05-07:00"
	defaultRankMatchWaitSec = 15 * time.Second
	defaultRankScan         = 3 * time.Second
)

type Tiers struct {
	Lv    int32
	Score int32
	Up    *Tiers
	Down  *Tiers
}

type RankExtend struct {
	score      int32
	waitTimes  int32
	scoreRange int32
	next       *RankExtend
}

func (r *RankExtend) Next() {
	if nil != r.next {
		*r = *r.next
	}
}

func (r *RankExtend) Jump(now int32) *RankExtend {
	if r.next == nil {
		return r
	} else {
		temp := r
		for nil != temp.next {
			if temp.next.waitTimes == now {
				return temp.next
			} else if r.next.waitTimes > now {
				return temp
			} else {
				temp = temp.next
			}
		}
		return temp
	}
}

func (r *RankExtend) Range(a int32, b int32) bool {
	if a > b {
		return a-b <= r.scoreRange
	} else {
		return b-a <= r.scoreRange
	}
}

type SeasonReward struct {
	ID    int32
	Value int32
}

type RankSeasonConfig struct {
	SeasonID int32
	Start    int64
	End      int64
	Lock     int64
	Bronze   []SeasonReward
	Silver   []SeasonReward
	Gold     []SeasonReward
	Platinum []SeasonReward
	Diamond  []SeasonReward
	King     []SeasonReward
}

type Manager struct {
	others         map[string]*sd.GlobalOthers // 通用数据
	characters     []int32                     // 角色编号
	functionF1     [][]*sd.FunctionF1
	functionF2     [][]*sd.FunctionF2
	functionF3     [][]*sd.FunctionF3
	signins        []*sd.SigninTask          // 签到任务数据
	matchTasks     []*sd.MatchTask           // 胜场任务数据
	zones          []*sd.GlobalZones         // 区域列表
	skus           []map[string]*sd.SKU      // 付费商品
	newVersions    map[string]semver.Version // 最新版本
	versions       []*sd.VersionUpdateAlert  // 版本更新信息
	activitys      []*ActivityConf           // 活动配置数据
	initRobotLvs   []*sd.InitRobotLvs        // 初始分数测试对战机器人
	aiToScore      []*sd.AItoScore           // ai等级查询分数
	scoreToAI      []*sd.ScoretoAI           // 分数查询ai等级
	scoreToAISkin  []int32                   // 分数查询ai皮肤
	robotNames     []*sd.RobotName           // 机器人名字
	rank           map[int32]*sd.Rank
	rankTiers      map[int32]*Tiers
	rankGroup      [][]*sd.RankGroup
	tiersTable     []byte
	rankExtend     []*RankExtend
	rankSeasonConf []*RankSeasonConfig
	rankForbid     []*sd.RankForbid
	rankCheatJudge [][]*sd.RankCheatJudge
	commodity      map[int32]*sd.Commodity
	channels       map[string]*sd.Channel
	fysdkIAPs      map[string]*sd.AndroidIAP
}

func newDisconfManager() *Manager {
	return &Manager{
		others:      make(map[string]*sd.GlobalOthers),
		characters:  []int32{},
		newVersions: make(map[string]semver.Version),
		rank:        map[int32]*sd.Rank{},
		rankTiers:   map[int32]*Tiers{},
		commodity:   map[int32]*sd.Commodity{},
	}
}

// 获取通用数据的整形值
func (m Manager) GetInt(id string) int32 {
	v, exists := m.others[id]
	if !exists {
		log.Printf("[配置模块]获取整数%s时出错: 字段不存在", id)
		return 0
	}

	return v.IntValue
}

// 获取通用数据的浮点数
func (m Manager) GetFloat(id string) float32 {
	v, exists := m.others[id]
	if !exists {
		log.Printf("[配置模块]获取浮点数%s时出错: 字段不存在", id)
		return 0
	}

	return v.FloatValue
}

// 获取通用数据的字符串
func (m Manager) GetString(id string) string {
	v, exists := m.others[id]
	if !exists {
		log.Printf("[配置模块]获取字符串%s时出错: 字段不存在", id)
		return ""
	}

	return v.StrValue
}

// 获取通用数据的列表
func (m Manager) GetList(id string) []int32 {
	v, exists := m.others[id]
	if !exists {
		log.Printf("[配置模块]获取列表%s时出错: 字段不存在", id)
		return nil
	}

	return v.ListValue
}

// 获取通用数据的字典
func (m Manager) GetDict(id string) map[int32]int32 {
	v, exists := m.others[id]
	if !exists {
		log.Printf("[配置模块]获取字典%s时出错: 字段不存在", id)
		return nil
	}

	return v.DictValue
}

func (m Manager) GetF1AndF3() ([][]*sd.FunctionF1, [][]*sd.FunctionF3) {
	return m.functionF1, m.functionF3
}

// 获取下次有效签到任务配置
func (m Manager) GetNextSigninTask(id int32) *sd.SigninTask {
	if id == 0 {
		return m.signins[0]
	}

	for _, v := range m.signins {
		if v.ID > id {
			return v
		}
	}

	return nil
}

// 获取有效胜场任务配置
func (m Manager) GetWinMatchTask(win int32) *sd.MatchTask {
	for _, v := range m.matchTasks {
		if win == v.Number {
			return v
		}
		if win < v.Number {
			break
		}
	}

	return nil
}

// 根据(自己的分数,对手分数)查出(自己赢了加多少分,输了扣多少分)
func (m Manager) GetScoreDelta(score, opponentScore int32) (int32, int32) {
	var i, j int
	for i < len(m.functionF2) && score >= m.functionF2[i][0].Score {
		i++
	}
	if i > 0 {
		i--
	}

	for j < len(m.functionF2[i]) && opponentScore >= m.functionF2[i][j].OpponentScore {
		j++
	}
	if j > 0 {
		j--
	}
	data := m.functionF2[i][j]

	return data.WinDelta, data.LoseDelta
}

func (m Manager) GetZoneList() []*sd.GlobalZones {
	return m.zones
}

// 获取全角色
func (self *Manager) GetCharacterIDList() []int32 {
	return self.characters
}

func (m Manager) GetAppStoreSKU(id string) (*sd.SKU, bool) {
	if len(m.skus) > 0 {
		app := m.skus[0]
		if app != nil {
			sku, ok := app[id]
			return sku, ok && sku.Num > 0 && len(sku.ID) > 0
		}
	}
	return nil, false
}

func (m Manager) GetNextRobotTestID(currentId int32, scoreDelta int32) int32 {
	for _, v := range m.initRobotLvs {
		if v.ID == currentId {
			return v.ScoreDeltaToNext[scoreDelta+10]
		}
	}

	return 0
}

func (m Manager) GetAItoScore(aiLevel int32) int32 {
	for _, v := range m.aiToScore {
		if v.AILevel == aiLevel {
			return v.Score
		}
	}

	return 0
}

func (m Manager) GetScoretoAI(score int32) int32 {
	var i int
	for _, v := range m.scoreToAI {
		if v.Score > score {
			break
		}
		i++
	}
	i--
	return m.scoreToAI[i].AILevel
}

func (m Manager) GetOneRobotName() string {
	return m.robotNames[rand.Intn(len(m.robotNames))].RobotName
}

func (m Manager) GetOneRobotSkin() int32 {
	return m.scoreToAISkin[rand.Intn(len(m.scoreToAISkin))]
}

// 版本检查
func (m *Manager) CheckNewVersion(channel, version string) (bool, bool, string, string) {
	var haveNewVersion, isForceUpdate bool
	var newVersion, downloadLink string

	currentVersion, err := semver.Make(version)
	if err != nil {
		log.Printf("解析客户端版本号时出错: %s - %s", channel, version)
		return haveNewVersion, isForceUpdate, newVersion, downloadLink
	}

	latestVersion, exists := m.newVersions[channel]
	if !exists {
		// log.Printf("获取最新版本号时未配置: %s - %s", channel, version)
		return haveNewVersion, isForceUpdate, newVersion, downloadLink
	}

	// 判断是否已经为最新版本
	if ok := currentVersion.GE(latestVersion); ok {
		return haveNewVersion, isForceUpdate, newVersion, downloadLink
	}

	haveNewVersion = true

	for _, v := range m.versions {
		if v.Channel != channel {
			continue
		}

		if v.Version == latestVersion.String() {
			newVersion = v.Version
			downloadLink = v.Link
		}

		if v.Force {
			isForceUpdate = true
		}
	}

	return haveNewVersion, isForceUpdate, newVersion, downloadLink
}

// 获取活动配置列表
func (m *Manager) GetActivityList() []*ActivityConf {
	return m.activitys
}

func (m Manager) GetInitRobotTestLevel() int32 {
	return m.initRobotLvs[0].RobotLevel
}

func (m Manager) GetRobotLevelByTestID(testId int32) int32 {
	for _, v := range m.initRobotLvs {
		if v.ID == testId {
			return v.RobotLevel
		}
	}
	return int32(1)
}

func (m *Manager) GetRank(id int32) (*sd.Rank, bool) {
	v, ok := m.rank[id&0xFFFF00]
	return v, ok
}

func (m *Manager) GetRankTiers(id int32) (*Tiers, bool) {
	v, ok := m.rankTiers[id&0xFFFFFF]
	return v, ok
}

func (m *Manager) IsRankTiersGuard(id int32) bool {
	v, ok := m.GetRank(id)
	if ok {
		return v.ScoreProtect > 0
	} else {
		return false
	}
}

func (m *Manager) GetRankExtend(score int32) *RankExtend {
	l := len(m.rankExtend)
	for i := 0; i < l; i++ {
		if m.rankExtend[i].score > score {
			if i == 0 {
				i++
			}
			r := *(m.rankExtend[i-1])
			return &r
		}
	}
	r := *(m.rankExtend[l-1])
	return &r
}

func (m *Manager) GetRankForbid(i int) *sd.RankForbid {
	l := len(m.rankForbid)
	if l <= i {
		i = l - 1
	}
	return m.rankForbid[i]
}

func (m *Manager) GetRankCheatJudge(sum int32, errNum int32) bool {
	l := len(m.rankCheatJudge)
	if l == 0 {
		return false
	}
	for k, v := range m.rankCheatJudge {
		if v[0].All == sum {
			return findRankCheatJudge(v, errNum)
		} else if v[0].All > sum {
			if k == 0 {
				return findRankCheatJudge(v, errNum)
			} else {
				return findRankCheatJudge(m.rankCheatJudge[k-1], errNum)
			}
		} else if k == l-1 {
			return findRankCheatJudge(v, errNum)
		}
	}
	return false
}

func findRankCheatJudge(c []*sd.RankCheatJudge, errNum int32) bool {
	l := len(c)
	for i := 0; i < l; i++ {
		if c[i].Suspicious > errNum {
			if i == 0 {
				return c[i].IsCheating == 1
			} else {
				return c[i-1].IsCheating == 1
			}
		} else if c[i].Suspicious == errNum {
			return c[i].IsCheating == 1
		}
	}
	return c[l-1].IsCheating == 1
}

func (m *Manager) GetRankMaxMatchTime() time.Duration {
	v, ok := m.others[KeyRankMaxMatchTime]
	if ok && v.IntValue > 0 {
		return time.Duration(v.IntValue) * time.Second
	}
	return defaultRankMatchWaitSec
}

func (m *Manager) GetRankGroupInterval() time.Duration {
	v, ok := m.others[KeyRankGroupInterval]
	if ok && v.IntValue > 0 {
		return time.Duration(v.IntValue) * time.Second
	}
	return defaultRankScan
}

func (m *Manager) GetSeasonConfig(now int64) (*RankSeasonConfig, bool) {
	var rankConfig *RankSeasonConfig
	for _, v := range m.rankSeasonConf {
		if now < v.End {
			if v.Start <= now {
				return v, true
			} else if rankConfig == nil || rankConfig.Start-now > v.Start-now {
				rankConfig = v
			}
		}
	}
	return rankConfig, rankConfig != nil
}

func (m *Manager) GetCommodity(id int32) (*sd.Commodity, bool) {
	v, ok := m.commodity[id]
	if ok {
		c := *v
		return &c, true
	} else {
		return nil, false
	}
}

func (m *Manager) GetGoldEveryTriWin() int32 {
	v, ok := m.others[KeyGoldEveryTriWin]
	if ok {
		return v.IntValue
	} else {
		return 0
	}
}

func (m *Manager) GetMaxTriWinReward() int32 {
	v, ok := m.others[KeyMaxTriWinReward]
	if ok {
		return v.IntValue
	} else {
		return 0
	}
}

func (m *Manager) GetDiamondRewardedAds() int32 {
	v, ok := m.others[KeyDiamondRewardedAds]
	if ok {
		return v.IntValue
	} else {
		return 0
	}
}

func (m *Manager) GetRewardedAdsCD() int32 {
	v, ok := m.others[KeyRewardedAdsCD]
	if ok {
		return v.IntValue
	} else {
		return 0
	}
}

func (m *Manager) GetMaxRewardedAdsPerDay() int32 {
	v, ok := m.others[KeyMaxRewardedAdsPerDay]
	if ok {
		return v.IntValue
	} else {
		return 0
	}
}

func (self *Manager) GetChannel(channel string) *sd.Channel {
	if channel != "" {
		return self.channels[channel]
	}

	return nil
}

func (self *Manager) GetFYSDKIAP(sku string) *sd.AndroidIAP {
	if sku != "" {
		return self.fysdkIAPs[sku]
	}

	return nil
}
