// Code auto-generated by cfgeditor.
// DO NOT EDIT!
package sdata

type Rank struct {
	ArenaLevel         int32           // 大段位
	SubLevel           int32           // 小段位
	LevelMaxStar       int32           // 段位总进度
	ResetLevel         []int32         // 重置段位
	ScoreProtect       int32           // 积分保护
	WinGoldRandomRange []int32         // 胜场金币随机范围
	MaxDailyGold       int32           // 日金币上限
	SeasonReward       map[int32]int32 // 赛季奖励
	ArenaChest         map[int32]int32 // 竞技场宝箱
	MatchChest         map[int32]int32 // 匹配宝箱
}
