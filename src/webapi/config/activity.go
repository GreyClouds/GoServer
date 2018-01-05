package config

type ActivityConf struct {
	ID      int32
	StartTS int64
	EndTS   int64
	Skins   map[int32]int32
}
