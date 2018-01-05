package activity

var DefaultActivityManager *ActivityManager

func init() {
	DefaultActivityManager = newActivityManager()
}

func newActivityManager() *ActivityManager {
	return &ActivityManager{
		activitys: make(map[int32]*Activity),
	}
}
