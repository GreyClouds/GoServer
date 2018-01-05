package bean

import (
	"log"
	"sync"
	"time"

	"github.com/astaxie/beego/orm"

	. "webapi/common"
)

const (
	DefaultChannelID = "default"
)

type Announcement struct {
	ID                int64  `orm:"column(id);auto;pk"`      // 主键
	Channel           string `orm:"column(channel);size(8)"` // 渠道编号
	Illustration      string `orm:"column(illu_all)"`        // 默认插图(英文)
	IllustrationZH_TW string `orm:"column(illu_zh_tw)"`      // 繁体插图
	IllustrationZH_CN string `orm:"column(illu_zh_cn)"`      // 简体插图
	Link              string `orm:"column(link)"`            // 外部链接
	Force             bool   `orm:"column(force)"`           // 是否强制显示
}

func (self *Announcement) TableName() string {
	return "announcement"
}

func (self *Announcement) GetIllustration(lang int32) string {
	switch lang {
	case ZH_CN:
		if self.IllustrationZH_CN != "" {
			return self.IllustrationZH_CN
		}
	case ZH_TW:
		if self.IllustrationZH_TW != "" {
			return self.IllustrationZH_TW
		}
	case EN_US:
		return self.Illustration
	}

	return self.Illustration
}

// ---------------------------------------------------------------------------

var gCacheAnnouncements []*Announcement
var gCacheAnnouncementLock sync.RWMutex

func PreloadAnnouncement() {
	doAnnouncementLoad(true)

	go func() {
		t := time.NewTicker(30 * time.Second)
		for {
			select {
			case <-t.C:
				doAnnouncementLoad(false)
			}
		}
	}()
}

func doAnnouncementLoad(verbose bool) {
	o := orm.NewOrm()

	var arr []*Announcement
	num, err := o.QueryTable(new(Announcement)).Limit(-1).All(&arr)
	if err != nil {
		log.Printf("更新公告缓存数据时出错: %v", err)
		return
	}

	gCacheAnnouncementLock.Lock()
	if num > 0 {
		gCacheAnnouncements = make([]*Announcement, num)
		for i := int64(0); i < num; i++ {
			gCacheAnnouncements[i] = arr[i]
		}
	} else {
		gCacheAnnouncements = nil
	}
	gCacheAnnouncementLock.Unlock()

	if verbose {
		log.Printf("预载入%d条公告记录", num)
	}
}

func GetAnnouncements() []*Announcement {
	gCacheAnnouncementLock.RLock()
	defer gCacheAnnouncementLock.RUnlock()

	if gCacheAnnouncements != nil {
		arr := make([]*Announcement, len(gCacheAnnouncements))
		for i := 0; i < len(gCacheAnnouncements); i++ {
			arr[i] = gCacheAnnouncements[i]
		}

		return arr
	}

	return nil
}
