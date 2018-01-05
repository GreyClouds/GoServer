package config

var (
	_instance *Manager
)

func init() {
	_instance = newDisconfManager()
}

func Singleton() *Manager {
	return _instance
}

func Conf() *Manager {
	return _instance
}

func GetRankGroup(n int32, diff int32) int {
	return int(_instance.GetRankGroup(n, diff))
}

func GetRankTiersJson() string {
	return string(_instance.tiersTable)
}
