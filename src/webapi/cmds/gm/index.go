package gm

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	pkgAccount "webapi/account"
	pkgBean "webapi/bean"
	pkgCDKey "webapi/cdkey"
	pkgConfig "webapi/config"
)

type ISkeleton interface {
	AccountManager() *pkgAccount.Manager
}

type Error struct {
	Result bool   `json:"result"`
	ErrMSG string `json:"errmsg"`
}

type AddWholeSkinResp struct {
	Result bool    `json:"result"`
	Skins  []int32 `json:"skins"`
}

type AddDiamondResp struct {
	Ok      bool  `json:"ok"`
	Diamond int32 `json:"diamond"`
	Change  int32 `json:"change"`
}

type GenerateCDKeyResp struct {
	Result bool     `json:"result"`
	Gift   int      `json:"gift"`
	CDKeys []string `json:"cdkeys"`
}

func parseUint32(str string) uint32 {
	v, _ := strconv.ParseUint(str, 10, 32)
	return uint32(v)
}

func parseInt32(str string) int32 {
	v, _ := strconv.ParseInt(str, 10, 32)
	return int32(v)
}

func handleGMDiamond(skeleton ISkeleton, w http.ResponseWriter, r *http.Request) {
	uid := parseUint32(r.FormValue("uid"))
	value := parseInt32(r.FormValue("value"))
	role := skeleton.AccountManager().GetRole(uid)
	res := &AddDiamondResp{}
	var m *pkgBean.Money
	if role != nil {
		m = role.GetMoney()
	} else {
		_, r := pkgBean.LoadCharacter(uid)
		if r != nil {
			m = pkgBean.LoadMoney(uid)
		}
	}
	if m != nil {
		m.ChangeDiamond(value, pkgBean.DiamondChangeGM)
		res.Ok = true
		res.Diamond = m.GetDiamond()
		res.Change = value
	}
	log.Println("GM 修改钻石, 用户:", uid, "结果:", res)
	json.NewEncoder(w).Encode(res)
}

func Handle(skeleton ISkeleton, w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	strAction := strings.TrimSpace(r.FormValue("action"))
	log.Printf("[GM] %s %s", r.RemoteAddr, strAction)

	pConfManager := pkgConfig.Singleton()

	switch strAction {
	case "diamond":
		handleGMDiamond(skeleton, w, r)
	case "kick":
		uid := parseUint32(r.FormValue("uid"))
		if uid == 0 {
			http.Error(w, "invalid param", http.StatusBadRequest)
			return
		}

		pAccountManager := skeleton.AccountManager()
		if ok := pAccountManager.Kick(uid); !ok {
			resp := &Error{Result: false, ErrMSG: "对方不在线"}
			json.NewEncoder(w).Encode(resp)
			return
		}

		fmt.Fprintf(w, `{"result":true}`)

	case "removead":
		uid := parseUint32(r.FormValue("uid"))
		if uid == 0 {
			http.Error(w, "invalid param", http.StatusBadRequest)
			return
		}

		pAccountManager := skeleton.AccountManager()
		roleob := pAccountManager.GetRole(uid)
		if roleob == nil {
			roleob = pAccountManager.LoadRole(uid)
			if roleob == nil {
				resp := &Error{Result: false, ErrMSG: "加载离线数据时出错"}
				json.NewEncoder(w).Encode(resp)
				return
			}
		}

		roleob.Lock()
		roleob.RemoveAD()
		roleob.Unlock()

		fmt.Fprintf(w, `{"result":true}`)

	case "wholeskins":
		uid := parseUint32(r.FormValue("uid"))
		if uid == 0 {
			http.Error(w, "invalid param", http.StatusBadRequest)
			return
		}

		pAccountManager := skeleton.AccountManager()
		roleob := pAccountManager.GetRole(uid)
		if roleob == nil {
			roleob = pAccountManager.LoadRole(uid)
			if roleob == nil {
				resp := &Error{Result: false, ErrMSG: "加载离线数据时出错"}
				json.NewEncoder(w).Encode(resp)
				return
			}
		}

		resp := &AddWholeSkinResp{
			Result: true,
			Skins:  []int32{},
		}

		roleob.Lock()

		ids := pConfManager.GetCharacterIDList()

		for _, v := range ids {
			if ok := roleob.AddSkin(v, pkgBean.SkinSrcGM); ok {
				resp.Skins = append(resp.Skins, v)
			}
		}

		roleob.Unlock()

		json.NewEncoder(w).Encode(resp)
	case "cdkey":
		strChannelID := r.FormValue("channel")
		strGift := r.FormValue("gift")
		category, _ := strconv.ParseUint(r.FormValue("category"), 10, 32)
		num, _ := strconv.ParseInt(r.FormValue("num"), 10, 32)
		deadline, _ := strconv.ParseInt(r.FormValue("deadline"), 10, 64)

		if strChannelID == "" || strGift == "" || num <= 0 || num > 100000 {
			http.Error(w, "invalid param", http.StatusBadRequest)
			return
		}

		if deadline != 0 && deadline <= time.Now().Unix() {
			http.Error(w, "invalid param", http.StatusBadRequest)
			return
		}

		if ok := pkgBean.ValidCDKeyCategory(uint32(category)); !ok {
			http.Error(w, "invalid param", http.StatusBadRequest)
			return
		}

		log.Printf("准备生成兑换码: [类型%d] [渠道%s] [数量%d] [资源%s] [有效期%d]", category, strChannelID, num, strGift, deadline)

		gift := pkgCDKey.ParseGiftString(strGift)
		if len(gift) == 0 {
			http.Error(w, "invalid param", http.StatusBadRequest)
			return
		}

		err, id, cdkeys := pkgCDKey.Generate(uint32(category), strChannelID, int(num), deadline, gift)
		if err != nil {
			resp := &Error{Result: false, ErrMSG: err.Error()}
			json.NewEncoder(w).Encode(resp)
			return
		}

		resp := &GenerateCDKeyResp{
			Result: true,
			Gift:   id,
			CDKeys: cdkeys,
		}
		json.NewEncoder(w).Encode(resp)
	case "cdkey2":
		strChannelID := r.FormValue("channel")
		category, _ := strconv.ParseUint(r.FormValue("category"), 10, 32)
		gift, _ := strconv.ParseInt(r.FormValue("gift"), 10, 32)
		num, _ := strconv.ParseInt(r.FormValue("num"), 10, 32)
		deadline, _ := strconv.ParseInt(r.FormValue("deadline"), 10, 64)

		if strChannelID == "" || gift == 0 || num <= 0 || num > 100000 {
			http.Error(w, "invalid param", http.StatusBadRequest)
			return
		}

		if ok := pkgBean.ValidCDKeyCategory(uint32(category)); !ok {
			http.Error(w, "invalid param", http.StatusBadRequest)
			return
		}

		if deadline != 0 && deadline <= time.Now().Unix() {
			http.Error(w, "invalid param", http.StatusBadRequest)
			return
		}

		log.Printf("准备生成兑换码: [类型%d] [渠道%s] [数量%d] [礼包%d] [有效期%d]", category, strChannelID, num, gift, deadline)

		if err, exists := pkgBean.IsGiftIDExists(int(gift)); err != nil {
			resp := &Error{Result: false, ErrMSG: err.Error()}
			json.NewEncoder(w).Encode(resp)
			return
		} else {
			if !exists {
				http.Error(w, "invalid param", http.StatusBadRequest)
				return
			}
		}

		err, id, cdkeys := pkgCDKey.GenerateWithGiftID(uint32(category), strChannelID, int(num), deadline, int(gift))
		if err != nil {
			resp := &Error{Result: false, ErrMSG: err.Error()}
			json.NewEncoder(w).Encode(resp)
			return
		}

		resp := &GenerateCDKeyResp{
			Result: true,
			Gift:   id,
			CDKeys: cdkeys,
		}
		json.NewEncoder(w).Encode(resp)
	case "invite_code":
		strChannelID := r.FormValue("channel")
		num, _ := strconv.ParseInt(r.FormValue("num"), 10, 32)
		deadline, _ := strconv.ParseInt(r.FormValue("deadline"), 10, 64)

		if strChannelID == "" || num <= 0 || num > 100000 {
			http.Error(w, "invalid param", http.StatusBadRequest)
			return
		}

		if deadline != 0 && deadline <= time.Now().Unix() {
			http.Error(w, "invalid param", http.StatusBadRequest)
			return
		}

		log.Printf("准备生成邀请码: [渠道%s] [数量%d] [有效期%d]", strChannelID, num, deadline)

		err, cdkeys := pkgCDKey.GenerateInviteCode(strChannelID, int(num), deadline)
		if err != nil {
			resp := &Error{Result: false, ErrMSG: err.Error()}
			json.NewEncoder(w).Encode(resp)
			return
		}

		resp := &GenerateCDKeyResp{
			Result: true,
			CDKeys: cdkeys,
		}
		json.NewEncoder(w).Encode(resp)
	default:
		http.Error(w, "", http.StatusBadRequest)
	}
}
