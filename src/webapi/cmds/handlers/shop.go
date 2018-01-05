package handlers

import (
	"log"

	"webapi/account"
	"webapi/bean"
	"webapi/config"
	"webapi/session"
	"webapi/skeleton"

	"time"

	sd "crazyant.com/deadfat/data/go"
	pkgProto "crazyant.com/deadfat/pbd/go"
	"github.com/golang/protobuf/proto"
)

const (
	OK             = 0
	ErrorID        = 1
	ErrorBalance   = 2
	ErrorCommodity = 3
	ErrorCD        = 1
	ErrorLimit     = 2
	RoleIDMin      = 1
	RoleIDMax      = 999
	GoldID         = 5001
	DiamondID      = 5002
)

// HandleSelfMoney 处理玩家发起的 获取自己的货币信息请求
func HandleSelfMoney(_ *skeleton.Skeleton, _ *session.Session, role *account.Role,
	_ proto.Message) (error, uint16, proto.Message) {

	rsp := &pkgProto.SelfMoneyResp{
		Money: role.SerializeMoney(),
	}

	return nil, uint16(pkgProto.GC_ID_SELF_MONEY_RESP), rsp
}

// HandleShopping 处理玩家发起的购买物品请求
func HandleShopping(_ *skeleton.Skeleton, _ *session.Session, role *account.Role, packet proto.Message) (error, uint16, proto.Message) {
	msg := packet.(*pkgProto.Shopping)
	resp := &pkgProto.ShoppingResp{}

	commodity, ok := config.Conf().GetCommodity(msg.Id)
	if ok {
		if checkRoleShopping(role, commodity) {
			shopping(role, commodity)
			resp.ErrorCode = proto.Int32(OK)
		} else {
			resp.ErrorCode = proto.Int32(ErrorBalance)
		}
	} else {
		resp.ErrorCode = proto.Int32(ErrorID)
	}
	resp.Money = role.SerializeMoney()

	return nil, uint16(pkgProto.GC_ID_SHOPPING_RESP), resp
}

func checkRoleShopping(r *account.Role, c *sd.Commodity) bool {
	if !checkRoleMoneyBalance(r.GetMoney(), c) {
		return false
	}
	if !c.Repeated {
		if c.ItemID <= RoleIDMax {
			return !r.SkinOwned(c.ItemID)
		}
	}
	return true
}

func checkRoleMoneyBalance(m *bean.Money, c *sd.Commodity) bool {
	switch c.MoneyType {
	case GoldID:
		return m.GetGold() >= c.MoneyPrice
	case DiamondID:
		return m.GetDiamond() >= c.MoneyPrice
	default:
		return false
	}
}

func shopping(role *account.Role, c *sd.Commodity) int32 {
	if shoppingPay(role.GetMoney(), c) {
		if shoppingDelivery(role, c) {
			return OK
		} else {
			log.Println("购买商品,[商品]配置错误:", c)
			return ErrorCommodity
		}
	} else {
		log.Println("购买商品,[价格]配置错误:", c)
		return ErrorCommodity
	}
}

func shoppingPay(m *bean.Money, c *sd.Commodity) bool {
	if c.MoneyPrice < 0 {
		return false
	}
	switch c.MoneyType {
	case GoldID:
		m.ChangeGold(-c.MoneyPrice)
	case DiamondID:
		m.ChangeDiamond(-c.MoneyPrice, commodityConsumptionType(c.ItemID))
	default:
		return false
	}
	return true
}

func commodityConsumptionType(i int32) int {
	switch {
	case GoldID == i:
		return bean.DiamondChangeBuyGold
	case i <= RoleIDMax:
		return bean.DiamondChangeBuySkin
	default:
		return bean.DiamondChangeError
	}
}

func shoppingDelivery(role *account.Role, c *sd.Commodity) bool {
	switch c.ItemID {
	case GoldID:
		role.GetMoney().ChangeGold(c.ItemAmount)
	default:
		if c.ItemID <= RoleIDMax {
			if c.MoneyType == DiamondID {
				role.AddSkin(c.ItemID, bean.SkinSrcDiamond)
			} else {
				role.AddSkin(c.ItemID, bean.SkinSrcGold)
			}
		} else {
			return false
		}
	}
	return true
}

func HandleAdReward(_ *skeleton.Skeleton, _ *session.Session, role *account.Role,
	_ proto.Message) (error, uint16, proto.Message) {
	adReward := role.GetAdReward()
	resp := &pkgProto.AdRewardResp{}
	now := time.Now().Unix()
	max := config.Conf().GetMaxRewardedAdsPerDay()
	if now < adReward.CD {
		resp.ErrorCode = proto.Int32(ErrorCD)
	} else {
		if adReward.Num < max {
			diamond := config.Conf().GetDiamondRewardedAds()
			cd := int64(config.Conf().GetRewardedAdsCD())
			role.DiamondAdd(diamond, bean.DiamondChangeTaskReward)
			adReward.Num++
			adReward.CD = now + cd
			adReward.Update()
			resp.ErrorCode = proto.Int32(0)
		} else {
			resp.ErrorCode = proto.Int32(ErrorLimit)
		}
	}
	resp.Money = role.SerializeMoney()
	resp.AdReward = proto.Bool(adReward.Num < max)
	resp.RemainT = proto.Int32(int32(adReward.CD - now))
	return nil, uint16(pkgProto.GC_ID_AD_REWARD_RESP), resp
}
