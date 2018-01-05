package activity

import (
	"testing"
	"time"
)

func TestRemainDays(t *testing.T) {
	now := int64(1509494399)

	for i := 0; i < 8; i++ {
		current := now + int64((i-4)*8*3600) + 8
		t1 := time.Unix(current, 0).Add(-8 * time.Hour)
		t2 := time.Unix(now, 0).Add(-8 * time.Hour)
		remain := int32((now - current + 86399) / 86400)
		switch i {
		case 0:
			if remain != 2 {
				t.Errorf("%s至%s剩余%d天", t1.Format("2006-01-02 15:04:05"), t2.Format("2006-01-02 15:04:05"), remain)
			}
		case 1, 2, 3:
			if remain != 1 {
				t.Errorf("%s至%s剩余%d天", t1.Format("2006-01-02 15:04:05"), t2.Format("2006-01-02 15:04:05"), remain)
			}
		default:
			if remain > 0 {
				t.Errorf("%s至%s剩余%d天", t1.Format("2006-01-02 15:04:05"), t2.Format("2006-01-02 15:04:05"), remain)
			}
		}

	}
}
