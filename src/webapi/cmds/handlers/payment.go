package handlers

import (
	"log"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/huangqingcheng/go-iap/appstore"

	pbd "crazyant.com/deadfat/pbd/go"

	oaccount "webapi/account"
	obean "webapi/bean"
	pkgConfig "webapi/config"
	osession "webapi/session"
	oskeleton "webapi/skeleton"
)

func HandleIAPVerify(skeleton *oskeleton.Skeleton, _ *osession.Session, roleob *oaccount.Role, packet proto.Message) (error, uint16, proto.Message) {
	p := packet.(*pbd.IAPVerifyForm)
	receipt := p.Receipt

	resp := &pbd.IAPVerifyResp{}

	// 判断参数是否错误
	if receipt == "" {
		resp.Err = doError(kBadRequest)

		return nil, uint16(pbd.GC_ID_IAP_VERIFY_RESP), resp
	}

	client := appstore.New()
	r := appstore.IAPRequest{
		ReceiptData: receipt,
	}

	var t appstore.IAPResponse
	if e := client.Verify(r, &t); e != nil {
		log.Printf("验证IAP时出错: %v\n%s", e, receipt)
		resp.Err = doError(kInternelServerError)

		return nil, uint16(pbd.GC_ID_IAP_VERIFY_RESP), resp
	}

	uid := roleob.Uid()
	transaction := t.Receipt.OriginalTransactionID
	productID := t.Receipt.ProductID

	// 判断配置表是否存在
	pConfManager := skeleton.ConfigManager()
	sku, ok := pConfManager.GetAppStoreSKU(productID)
	if !ok {
		log.Printf("AppStore支付失败[%d]：目标编号%s配置不存在", uid, productID)
		resp.Err = doError(kInternelServerError)

		return nil, uint16(pbd.GC_ID_IAP_VERIFY_RESP), resp
	}

	if err, exists := obean.ExistsAppStorePayment(uid, transaction); err != nil {
		log.Printf("验证AppStore支付订单是否重复时出错：%v", err)
		resp.Err = doError(kInternelServerError)

		return nil, uint16(pbd.GC_ID_IAP_VERIFY_RESP), resp
	} else {
		if exists {
			log.Printf("AppStore支付订单已存在%d: %s", uid, transaction)
			resp.Err = doError(kPaymentAlreadyExists)

			return nil, uint16(pbd.GC_ID_IAP_VERIFY_RESP), resp
		}
	}

	// 判断是否重复领取
	err, bean := obean.NewAndSaveAppStorePayment(time.Now(), uid, transaction, productID)
	if err != nil {
		log.Printf("插入AppStore支付订单时出错：%v", err)
		resp.Err = doError(kInternelServerError)

		return nil, uint16(pbd.GC_ID_IAP_VERIFY_RESP), resp
	}

	if sku.RemoveAds > 0 {
		roleob.RemoveAD()
	}

	isFirstRecharge := sku.IsFirstRedeem && roleob.IsFirstRecharge()

	// 资源编号
	var resources []*pbd.Resource
	for _, id := range sku.ID {
		switch {
		case id == DiamondID:
			if isFirstRecharge {
				roleob.DiamondAdd(sku.Num*2, obean.DiamondChangePay)
				resources = append(resources, &pbd.Resource{
					Id:    proto.Int32(id),
					Value: proto.Int32(sku.Num * 2),
				})
			} else {
				roleob.DiamondAdd(sku.Num, obean.DiamondChangePay)
				resources = append(resources, &pbd.Resource{
					Id:    proto.Int32(id),
					Value: proto.Int32(sku.Num),
				})
			}
		case id >= RoleIDMin && id <= RoleIDMax:
			if !roleob.AddSkin(id, obean.SkinSrcPay) {
				continue
			}

			resources = append(resources, &pbd.Resource{
				Id:    proto.Int32(id),
				Value: proto.Int32(sku.Num),
			})
		}
	}

	if isFirstRecharge {
		id := pConfManager.GetInt(pkgConfig.KeyFirstRechargeSkin)

		if ok := roleob.AddSkin(id, obean.SkinSrcPay); ok {
			resources = append(resources, &pbd.Resource{
				Id:    proto.Int32(id),
				Value: proto.Int32(1),
			})
		}

		roleob.CleanFirstRechargeFlag()
	}

	resp.RemoveAd = proto.Bool(roleob.IsRemoveAD())
	resp.Id = proto.Uint32(uint32(bean.ID))
	if len(resources) > 0 {
		resp.Resources = resources
	}
	resp.Money = roleob.SerializeMoney()

	return nil, uint16(pbd.GC_ID_IAP_VERIFY_RESP), resp
}
