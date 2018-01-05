package trace

const (
	ProcessStatusReport = 1000 // 进程状态上报
	UserLogin           = 1001 // 用户登录[channel,imei,way,userid,uid,ip,version,happen]
	GameMatch           = 1002 // 游戏匹配[category,roomid,uid1,uid2]
	NormalMatchReport   = 1003 // 普通匹配状态上报[join,leave,waitings]
)

const (
	NormalMatch      = uint32(1) // 普通匹配
	OrientationMatch = uint32(2) // 定向匹配
	ArenaMatch       = uint32(3) // 排位
)
