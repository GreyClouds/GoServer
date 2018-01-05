package config

const (
	KeyInitialScore         = "InitialScore"
	KeyMaxMatchTime         = "MaxMatchTime"
	KeyGroupInterval        = "GroupInterval"
	KeyInitSkinIDList       = "InitialChara"         // 初始赠送皮肤
	KeyMaxTestBattles       = "MaxTestBattles"       // 新手初始与ai对战的最大场次
	KeyNoScoreLine          = "NoScoreLine"          // 超过这个分数后 与ai对战不再获得分数
	KeyRobotProtect         = "RobotProtect"         // 机器人保护连败场次
	KeyRobotWinsMin         = "RobotWinsMin"         // 机器人胜场数比例下限
	KeyRobotWinsMax         = "RobotWinsMax"         // 机器人胜场数比例上限
	KeyRankMaxMatchTime     = "RankMaxMatchTime"     // 排位匹配超时(秒)
	KeyRankGroupInterval    = "RankGroupInterval"    // 排位匹配间隔时间(秒)
	KeyGoldEveryTriWin      = "GoldEveryTriWin"      // 每三胜奖励金币
	KeyMaxTriWinReward      = "MaxTriWinReward"      // 每天每三胜奖励次数上限
	KeyDiamondRewardedAds   = "DiamondRewardedAds"   // 看视频广告钻石数量
	KeyRewardedAdsCD        = "RewardedAdsCD"        // 视频广告CD(秒)
	KeyMaxRewardedAdsPerDay = "MaxRewardedAdsPerDay" // 明天视频广告上限
	KeyFirstRechargeSkin    = "FirstRechargeRole"    // 首充赠送皮肤
)
