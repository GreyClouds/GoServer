package config

//import (
//	"testing"
//
//	pkgDisconf "yunjing.me/phoenix/go-phoenix/disconf"
//)

//func TestIsActivityStart(t *testing.T) {
//	manager := NewDisconfManager()
//	manager.Preload(pkgDisconf.New("/data/deadfat/data/"))
//
//	// 测试线上
//	state, ids := manager.IsActivityStart(1, "1.5.0")
//	if state != 1 {
//		t.Errorf("Online version Activity NOT open: %v", ids)
//	}
//
//	// 测试线上新版本
//	state, ids = manager.IsActivityStart(1, "1.5.1")
//	if state != 1 {
//		t.Errorf("New version Client Activity NOT open: %v", ids)
//	}
//
//	// 测试旧版本: 旧版本不开放活动
//	state, ids = manager.IsActivityStart(1, "")
//	if state != 0 {
//		t.Errorf("Old version Client Activity is open: %v", ids)
//	}
//
//	// 测试未来新颁布的版本
//	state, ids = manager.IsActivityStart(1, "1.5.2")
//	if state != 1 {
//		t.Errorf("Develop version Client Activity NOT open: %v", ids)
//	}
//
//}
