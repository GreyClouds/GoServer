package bean

import (
	"log"

	"github.com/astaxie/beego/orm"
)

type SkinTask struct {
	ID     int64  `orm:"column(id);auto;pk"` // 主键
	Uid    uint32 `orm:"column(uid);index"`  // 角色编号
	SkinID int32  `orm:"column(skin_id)"`    // 皮肤编号
	WinNum int32  `orm:"column(win_num)"`    // 使用皮肤胜场
	dirty  bool   `orm:"-"`                  // 是否脏数据
}

func (self *SkinTask) TableName() string {
	return "skin_task"
}

func (self *SkinTask) SetDirty() {
	self.dirty = true
}

func (self *SkinTask) ResetDirty() bool {
	if !self.dirty {
		return false
	}

	self.dirty = false

	return true
}

func NewSkinTask(uid uint32, skinID int32) (error, *SkinTask) {
	o := orm.NewOrm()

	skinTask := &SkinTask{
		Uid:    uid,
		SkinID: skinID,
		WinNum: 1,
	}

	_, e := o.Insert(skinTask)
	if e != nil {
		log.Printf("角色皮肤任务%d插入皮肤%d记录时出错: %v", uid, skinID, e)
		return e, nil
	}

	return nil, skinTask
}

func LoadSkinTasks(uid uint32) (error, []*SkinTask) {
	o := orm.NewOrm()

	var skinTasks []*SkinTask

	qs := o.QueryTable("skin_task")
	_, e := qs.Filter("uid", uid).All(&skinTasks)
	if e != nil {
		return e, nil
	}

	// log.Printf("角色%d载入%d条皮肤任务进度", uid, n)

	return nil, skinTasks
}

func UpdateSkinTask(bean *SkinTask) error {
	o := orm.NewOrm()
	_, e := o.Update(bean)
	if e != nil {
		return e
	}
	return nil
}
