package fysdk

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	pkgBean "webapi/bean"
)

type ISkeleton interface {
	OnPaid(*pkgBean.AndroidPayment)
}

// 支付通知自定义字段
type FYSDKPaymentExt struct {
	SKU string `json:"sku"`
}

func PayNotify(skeleton ISkeleton, w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	if ok := CheckSign(r.Form, GetPaySecret()); !ok {
		log.Printf("收到FYSDK订单验签失败: %v", r.Form)
		http.Error(w, "sign_failure", http.StatusBadRequest)
		return
	}

	var ext FYSDKPaymentExt
	if v := r.FormValue("app_callback_ext"); v != "" {
		err := json.Unmarshal([]byte(v), &ext)
		if err != nil {
			log.Printf("收到FYSDK订单: 自定义参数格式不合法")
			http.Error(w, "ext_invalid", http.StatusBadRequest)
			return
		}
	}

	bean := &pkgBean.AndroidPayment{
		OrderID:  r.FormValue("order_id"),                   // SDK订单号
		UUID:     r.FormValue("uuid"),                       // SDK唯一用户编号
		ZoneID:   parseInt32(r.FormValue("app_zone_id")),    // 游戏大区编号
		UID:      parseUint32(r.FormValue("app_player_id")), // 游戏账号
		SKU:      ext.SKU,                                   // 商品编号
		Amount:   parseInt32(r.FormValue("pay_amount")),     // 支付数量
		PayTime:  parseInt64(r.FormValue("pay_time")),       // 订单支付时间
		Sandbox:  r.FormValue("sandbox") == "sandbox",       // 是否测试订单
		Happen:   parseInt64(r.FormValue("time")),           // 发生时间
		Achieved: false,                                     // 是否已认领
	}

	if err := pkgBean.InsertAndroidPayment(bean); err != nil {
		if exists := strings.HasPrefix(err.Error(), "Error 1062: Duplicate entry"); exists {
			// 重复的订单
			fmt.Fprintf(w, "ok")
			return
		}

		http.Error(w, "sign_failure", http.StatusInternalServerError)
		return
	}

	skeleton.OnPaid(bean)

	fmt.Fprintf(w, "ok")
	return
}
