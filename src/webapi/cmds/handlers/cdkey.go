package handlers

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"

	pbd "crazyant.com/deadfat/pbd/go"

	oaccount "webapi/account"
	pkgBean "webapi/bean"
	osession "webapi/session"
	oskeleton "webapi/skeleton"
)

func HandleCDKeyExchange(skeleton *oskeleton.Skeleton, _ *osession.Session, roleob *oaccount.Role, packet proto.Message) (error, uint16, proto.Message) {
	payload := packet.(*pbd.CDKey)
	resp := &pbd.CDKeyExchangeResp{}

	cdkey := strings.ToUpper(payload.Id)
	uid := roleob.Uid()

	// 判断参数是否合理
	if len(cdkey) < 3 || len(cdkey) > 9 {
		resp.Err = doError(kCDKeyInvalid)
		return nil, uint16(pbd.GC_ID_CDKEY_EXCHANGE_RESP), resp
	}

	// 查询兑换码是否存在
	err, beanCDKey, beanGift := pkgBean.QueryCDKey(cdkey)
	if err != nil {
		resp.Err = doError(kCDKeyInvalid)
		return nil, uint16(pbd.GC_ID_CDKEY_EXCHANGE_RESP), resp
	}

	// 判断是否存在
	if beanCDKey == nil || beanGift == nil {
		resp.Err = doError(kCDKeyNotExists)
		return nil, uint16(pbd.GC_ID_CDKEY_EXCHANGE_RESP), resp
	}

	// 判断是否满足限定渠道
	if roleob.GetChannelID() != beanCDKey.ChannelID {
		resp.Err = doError(kCDKeyChannelNotMatch, beanCDKey.ChannelID)
		return nil, uint16(pbd.GC_ID_CDKEY_EXCHANGE_RESP), resp
	}

	// 判断是否过期
	if beanCDKey.Deadline > 0 && time.Now().Unix() > beanCDKey.Deadline {
		resp.Err = doError(kCDKeyIsTimeout, fmt.Sprintf("%d", beanCDKey.Deadline))
		return nil, uint16(pbd.GC_ID_CDKEY_EXCHANGE_RESP), resp
	}

	// 判断是否已被使用
	if beanCDKey.Uid != 0 {
		resp.Err = doError(kCDKeyIsAchieved)
		return nil, uint16(pbd.GC_ID_CDKEY_EXCHANGE_RESP), resp
	}

	switch beanCDKey.Category {
	case pkgBean.NormalCDKey:
		// 判断是否已领取过同一礼包的兑换码
		err, exists := pkgBean.IsExchangeGiftCDKEY(uid, beanCDKey.Gift)
		if err != nil {
			log.Printf("查询兑换码%s是否已被%d领取过时出错: %v", cdkey, uid, err)
			resp.Err = doError(kInternelServerError)
			return nil, uint16(pbd.GC_ID_CDKEY_EXCHANGE_RESP), resp
		}

		if exists {
			resp.Err = doError(kCDKeyGiftAlreadyUse)
			return nil, uint16(pbd.GC_ID_CDKEY_EXCHANGE_RESP), resp
		}
	case pkgBean.DailyCDKey:
		// 判断是否已领取过同一礼包的兑换码
		err, exists := pkgBean.IsExchangeDailyGiftCDKey(uid, beanCDKey.Gift)
		if err != nil {
			log.Printf("查询兑换码%s是否已被%d领取过时出错: %v", cdkey, uid, err)
			resp.Err = doError(kInternelServerError)
			return nil, uint16(pbd.GC_ID_CDKEY_EXCHANGE_RESP), resp
		}

		if exists {
			resp.Err = doError(kCDKeyGiftAlreadyUse)
			return nil, uint16(pbd.GC_ID_CDKEY_EXCHANGE_RESP), resp
		}
	}

	// 标记兑换成功
	beanCDKey.SetAchieved(uid)
	if err := pkgBean.UpdateCDKey(beanCDKey); err != nil {
		log.Printf("更新兑换码%s被%d领取记录时出错: %v", cdkey, uid, err)
		resp.Err = doError(kInternelServerError)
		return nil, uint16(pbd.GC_ID_CDKEY_EXCHANGE_RESP), resp
	}

	// 资源结算
	pAccountManager := skeleton.AccountManager()
	rewards := pAccountManager.Reward(roleob, beanGift.GetResources(), pkgBean.DiamondChangeCDKey)
	if rewards != nil && len(rewards) > 0 {
		// 兑换码的返回内容为原始发送物品列表
		for id, count := range beanGift.GetResources() {
			if count <= 0 {
				continue
			}

			resp.Resources = append(resp.Resources, &pbd.Resource{
				Id:    proto.Int32(id),
				Value: proto.Int32(count),
			})
		}
	}

	return nil, uint16(pbd.GC_ID_CDKEY_EXCHANGE_RESP), resp
}
