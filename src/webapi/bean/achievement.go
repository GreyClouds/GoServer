package bean

import (
	"github.com/astaxie/beego/orm"
	"fmt"
	"time"
)

type AchievementAttr struct {
	AchieveName       string    `orm:"pk;column(achieveName)"`           // 成就属性名称 + 角色编号
	AchieveValue int32 `orm:"column(achieveValue)"`       // 成就属性值
}

type AchievementUnLock struct {
	AchieveName       string    `orm:"pk;column(achieveName)"`           // 成就名称 + 角色编号
	Date string `orm:"column(data)"`       // 成就解锁日期字符串
}


func (self *AchievementAttr) TableName() string {
	return "achievement"
}

// ---------------------------------------------------------------------------
func GetAchievement(uid uint32, name string) int32 {
	achieveKey := fmt.Sprintf("%s%s", name, uid);
	var achievement =  AchievementAttr{AchieveName:achieveKey}
	err := defaultOrm.Read(&achievement)
	if err == nil{
		return  achievement.AchieveValue
	}else if(err == orm.ErrNoRows){
		return 0
	}else{
		return 0
	}
}

func SetAchievement(uid uint32, name string, value int32) {
	achieveKey := fmt.Sprintf("%s%s", name, uid);
	var achievement =  AchievementAttr{AchieveName:achieveKey, AchieveValue:value}
	defaultOrm.InsertOrUpdate(&achievement)
}

func UnLockAchievement(uid uint32, name string) string{
	achieveKey := fmt.Sprintf("%s%s", name, uid);
	timestamp := time.Now().Unix()
	tm := time.Unix(timestamp, 0);
	dateString := tm.Format("2006/01/02")
	var achievement =  AchievementUnLock{AchieveName:achieveKey, Date:dateString}
	defaultOrm.InsertOrUpdate(&achievement)
	return dateString
}

func GetUnLockAchieveDate(uid uint32, name string) string {
	achieveKey := fmt.Sprintf("%s%s", name, uid);
	var achievement =  AchievementUnLock{AchieveName:achieveKey}
	err := defaultOrm.Read(&achievement)
	if err == nil{
		return  achievement.Date
	}else if(err == orm.ErrNoRows){
		return ""
	}else{
		return ""
	}
}


