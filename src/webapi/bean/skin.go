package bean

import (
	"log"

	"github.com/astaxie/beego/orm"
)

const (
	SkinSrcError      = iota - 1 //错误
	SkinSrcOld                   //旧的未知方式
	SkinSrcPay                   //充值支付购买
	SkinSrcDiamond               //钻石购买
	SkinSrcGold                  //金币购买
	SkinSrcCDKey                 //CDKey兑换
	SkinSrcTaskReward            //任务奖励
	SkinSrcActive                //活动赠送
	SkinSrcInit                  //初始获取
	SkinSrcGM                    //GM发放
)

type Skin struct {
	ID     int64  `orm:"column(id);auto;pk"` // 主键
	Uid    uint32 `orm:"column(uid);index"`  // 角色编号
	SkinID int32  `orm:"column(skin_id)"`    // 皮肤编号
	Source int    `orm:"column(source)"`
}

func (self Skin) TableName() string {
	return "skin"
}

// ---------------------------------------------------------------------------

func NewSkins(uid uint32, skins []int32, src int) (error, []*Skin) {
	result := make([]*Skin, 0, len(skins))
	for _, id := range skins {
		skin := &Skin{
			Uid:    uid,
			SkinID: id,
			Source: src,
		}
		result = append(result, skin)
	}

	_, e := defaultOrm.InsertMulti(1, result)
	if e != nil {
		log.Printf("角色%d插入皮肤记录时出错: %v", uid, e)
	}
	return nil, result
}

func LoadSkins(uid uint32) (error, []*Skin) {
	o := orm.NewOrm()

	var skins []*Skin

	qs := o.QueryTable("skin")
	_, e := qs.Filter("uid", uid).All(&skins)
	if e != nil {
		return e, nil
	}

	// log.Printf("角色%d载入%d条皮肤记录", uid, n)

	return nil, skins
}
