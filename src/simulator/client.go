package main

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"

	pbd "crazyant.com/deadfat/pbd/go"
	pkgPhoenixPBD "yunjing.me/phoenix/pbd/go"
)

// ---------------------------------------------------------------------------

const (
	kAPISecret = "Lo01v8!P" //
)

func init() {
	http.DefaultClient.Timeout = 6 * time.Second
}

type Client struct {
	prefix     string
	prototypes map[uint16]proto.Message

	uid   uint32
	token []byte
	imei  string
}

func newClient(addr string, i string) *Client {
	self := &Client{
		prefix:     fmt.Sprintf("%s/o", addr),
		prototypes: make(map[uint16]proto.Message),
		imei:       i,
	}

	self.prototypes[uint16(pkgPhoenixPBD.MSG_ERROR)] = &pkgPhoenixPBD.Error{}
	self.prototypes[uint16(pkgPhoenixPBD.MSG_SERVER_LIST_AND_ACCESS_TOKEN)] = &pkgPhoenixPBD.AccessTokenAndServerList{}
	self.prototypes[uint16(pbd.GC_ID_GUEST_LOGIN_RESP)] = &pbd.GuestLoginResp{}
	self.prototypes[uint16(pbd.GC_ID_MATCH_BEGIN_RESP)] = &pbd.MatchBeginResp{}
	self.prototypes[uint16(pbd.GC_ID_MATCH_QUERY_RESP)] = &pbd.MatchQueryResp{}
	self.prototypes[uint16(pbd.GC_ID_MATCH_REWARD_RESP)] = &pbd.MatchRewardResp{}
	self.prototypes[uint16(pbd.GC_ID_SIGNIN_VIEW_RESP)] = &pbd.SigninViewResp{}
	self.prototypes[uint16(pbd.GC_ID_SIGNIN_RESP)] = &pbd.SigninResp{}
	self.prototypes[uint16(pbd.GC_ID_NICK_SET_RESP)] = &pbd.NickSetResp{}
	self.prototypes[uint16(pbd.GC_ID_ZONE_LIST_RESP)] = &pbd.ZoneListResp{}
	self.prototypes[uint16(pbd.GC_ID_ZONE_SET_RESP)] = &pbd.ZoneSetResp{}
	self.prototypes[uint16(pbd.GC_ID_MATCH_END_RESP)] = &pbd.MatchEndResp{}
	self.prototypes[uint16(pbd.GC_ID_MATCH_REPEAT_STATE_RESP)] = &pbd.MatchRepeatStateResp{}
	self.prototypes[uint16(pbd.GC_ID_IAP_VERIFY_RESP)] = &pbd.IAPVerifyResp{}
	self.prototypes[uint16(pbd.GC_ID_INVITE_CODE_RESP)] = &pbd.InviteCodeResp{}
	self.prototypes[uint16(pbd.GC_ID_CDKEY_EXCHANGE_RESP)] = &pbd.CDKeyExchangeResp{}
	self.prototypes[uint16(pbd.GC_ID_HEART_BEAT_RESP)] = &pbd.HeartBeat{}

	self.prototypes[uint16(pbd.GC_ID_RANK_BEGIN_RESP)] = &pbd.RankBeginResp{}
	self.prototypes[uint16(pbd.GC_ID_RANK_QUERY_RESP)] = &pbd.RankQueryResp{}
	self.prototypes[uint16(pbd.GC_ID_RANK_REWARD_RESP)] = &pbd.RankRewardResp{}
	self.prototypes[uint16(pbd.GC_ID_RANK_CANCEL_RESP)] = &pbd.RankCancelResp{}
	self.prototypes[uint16(pbd.GC_ID_RANK_INFO_RESP)] = &pbd.RankInfoResp{}
	self.prototypes[uint16(pbd.GC_ID_RANK_SELF_RESP)] = &pbd.RankSelfResp{}
	self.prototypes[uint16(pbd.GC_ID_RANK_TIERS_RESP)] = &pbd.RankTiersResp{}
	self.prototypes[uint16(pbd.GC_ID_RANK_SEASONREWARD_RESP)] = &pbd.RankSeasonRewardResp{}
	self.prototypes[uint16(pbd.GC_ID_SHOPPING_RESP)] = &pbd.ShoppingResp{}
	self.prototypes[uint16(pbd.GC_ID_SELF_MONEY_RESP)] = &pbd.SelfMoneyResp{}
	self.prototypes[uint16(pbd.GC_ID_AD_REWARD_RESP)] = &pbd.AdRewardResp{}

	return self
}

func (self *Client) Error(e *pbd.Error) error {
	var text string
	if args := e.Args; args != nil && len(args) > 0 {
		text = fmt.Sprintf("%d: %s", e.Code, strings.Join(args, ","))
	} else {
		text = fmt.Sprintf("%d", e.Code)
	}
	return errors.New(text)
}

func (self *Client) sign(rn uint16, uid uint32, token []byte, pid uint16, payload []byte) []byte {
	// 校验签名
	idx := 0
	total := 6 + 2 + len([]byte(kAPISecret))
	if uid != 0 && (token != nil && len(token) > 0) {
		total += len(token)
	}

	if payload != nil && len(payload) > 0 {
		total += len(payload)
	}

	raw := make([]byte, total)

	binary.LittleEndian.PutUint16(raw, rn)
	idx += 2

	binary.LittleEndian.PutUint32(raw[2:], uid)
	idx += 4

	// hash.Write(raw)
	if uid != 0 {
		if token == nil || len(token) == 0 {
			return []byte{}
		}

		copy(raw[idx:], token)
		idx += len(token)
	}

	// raw1 := make([]byte, 2)
	binary.LittleEndian.PutUint16(raw[idx:], pid)
	idx += 2
	// hash.Write(raw1)

	if payload != nil && len(payload) > 0 {
		// hash.Write(self.payload)
		copy(raw[idx:], payload)
		idx += len(payload)
	}

	// hash.Write([]byte(kAPISecret))
	copy(raw[idx:], []byte(kAPISecret))
	idx += len([]byte(kAPISecret))

	hash := md5.New()
	hash.Write([]byte(base64.StdEncoding.EncodeToString(raw)))

	return hash.Sum(nil)
}

func (self *Client) request(needAuth bool, id pbd.CG, payload proto.Message) (error, proto.Message) {
	var plain []byte
	var buf bytes.Buffer

	rn := uint16(rand.Intn(65536))

	{
		raw := make([]byte, 2)
		binary.LittleEndian.PutUint16(raw, rn)
		buf.Write(raw)
	}

	{
		raw := make([]byte, 4)
		binary.LittleEndian.PutUint32(raw, self.uid)
		buf.Write(raw)
	}

	if needAuth {
		buf.Write(self.token)
	}

	{
		raw := make([]byte, 2)
		binary.LittleEndian.PutUint16(raw, uint16(id))
		buf.Write(raw)
	}

	if payload != nil {
		plain, _ = proto.Marshal(payload)
		if len(plain) > 0 {
			buf.Write(plain)
		}
	}

	if needAuth {
		buf.Write(self.sign(rn, self.uid, self.token, uint16(id), plain))
	} else {
		buf.Write(self.sign(rn, 0, nil, uint16(id), plain))
	}

	c := &http.Client{
		Timeout: 6 * time.Second,
	}

	resp, err := c.Post(self.prefix, "", &buf)
	if err != nil {
		// log.Printf("服务请求时出错: %v", err)
		return err, nil
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("读取服务请求返回数据时出错: %v", err)
		return err, nil
	}
	if len(body) < 2 {
		log.Printf("读取服务请求%d返回内容过短：%v", id, body)
		return errors.New("1"), nil
	}

	respId := binary.LittleEndian.Uint16(body)

	if tpl, exists := self.prototypes[respId]; exists {
		packet := proto.Clone(tpl)
		if err := proto.Unmarshal(body[2:], packet); err != nil {
			log.Printf("解析协议%d时出错: %v", respId, err)
			return err, nil
		}

		return nil, packet
	} else {
		if len(body) > 2 {
			log.Printf("协议%d母本未定义", respId)
		}

		return nil, nil
	}
}

// ---------------------------------------------------------------------------

// 登陆
func (self *Client) Login(channel, version string) bool {
	payload := &pbd.GuestLogin{
		Channel:       proto.String(channel),
		Imei:          proto.String(self.imei),
		Os:            proto.Int32(0),
		ClientVersion: proto.String(version),
	}

	err, resp := self.request(false, pbd.CG_ID_GUEST_LOGIN, payload)
	if err != nil {
		log.Printf("登陆请求时出错: %v", err)
		return false
	}

	packet := resp.(*pbd.GuestLoginResp)

	if e := packet.GetErr(); e != nil {
		log.Printf("登陆请求时失败: %v", e)
		return false
	}

	self.uid = packet.GetUid()

	self.token = make([]byte, len(packet.GetToken()))
	copy(self.token, packet.GetToken())

	return true
}

func (self *Client) ViewSignTask() {
	err, resp := self.request(true, pbd.CG_ID_SIGNIN_VIEW, nil)
	if err != nil {
		log.Printf("查询签到进度请求时出错%d: %v", self.uid, err)
		return
	}

	packet := resp.(*pbd.SigninViewResp)
	_ = packet
}

func (self *Client) InviteCode() {
	payload := &pbd.InviteCodeReq{
		InviteCode: proto.String("9sS8sl20"),
	}

	err, resp := self.request(true, pbd.CG_ID_INVITE_CODE, payload)
	if err != nil {
		log.Printf("邀请码输入错误%d: %v", self.uid, err)
		return
	}

	packet := resp.(*pbd.InviteCodeResp)

	if e := packet.GetErr(); e != nil {
		if e.GetCode() == 1120 {
			return
		}

		log.Printf("邀请码输入错误%d: %v", self.uid, e)
	}
}

func (self *Client) CDKey() {
	payload := &pbd.CDKey{
		Id: proto.String("2ZQVFWCPK"),
	}

	err, resp := self.request(true, pbd.CG_ID_CDKEY_EXCHANGE, payload)
	if err != nil {
		log.Printf("兑换码输入错误%d: %v", self.uid, err)
		return
	}

	packet := resp.(*pbd.CDKeyExchangeResp)

	if e := packet.GetErr(); e != nil {
		if e.GetCode() == 1101 {
			return
		}

		log.Printf("兑换码输入错误%d: %v", self.uid, e)
	}
}

// 请求好友匹配
func (self *Client) MatchBegin() error {
	payload := &pbd.MatchBegin{
		Skin:     proto.Int32(1),
		Lag:      proto.Int32(int32(rand.Intn(400))),
		RoomCode: proto.Int32(-1),
		IsAgain:  proto.Bool(false),
	}

	err, resp := self.request(true, pbd.CG_ID_MATCH_BEGIN, payload)
	if err != nil {
		log.Printf("发起好友匹配请求时出错%d: %v", self.uid, err)
		return err
	}

	packet := resp.(*pbd.MatchBeginResp)

	if e := packet.GetErr(); e != nil {
		log.Printf("发起匹配请求时失败%d: %v", self.uid, e)
		return self.Error(e)
	}

	return nil
}

// 好友匹配查询
func (self *Client) MatchQuery() (error, bool, uint32, proto.Message) {
	payload := &pbd.MatchQuery{}

	err, resp := self.request(true, pbd.CG_ID_MATCH_QUERY, payload)
	if err != nil {
		log.Printf("询问好友匹配是否成功出错%d: %v", self.uid, err)
		return err, true, 0, nil
	}

	packet := resp.(*pbd.MatchQueryResp)

	switch packet.GetStatus() {
	case 0:
		// 在查找中
		return nil, true, 0, nil
	case 1:
		// 失败
		// println("failure")
		return nil, false, 0, nil
	case 2:
		// 成功
		// log.Printf("[房间%d][对手%d|%d|%s]", packet.GetRoomId(), packet.GetVsAi(), packet.GetOpponentUid(), packet.GetName())
		return nil, false, packet.GetRoomId(), packet
	default:
		return errors.New("unknown_resp_code"), true, 0, nil
	}
}

func (self *Client) MatchReward(roomid uint32, ai int32, opponent uint32) error {
	details := []*pbd.MatchMemberDetail{
		&pbd.MatchMemberDetail{
			Uid: proto.Uint32(opponent),
			Hp:  proto.Int32(1),
		},
		&pbd.MatchMemberDetail{
			Uid: proto.Uint32(self.uid),
			Hp:  proto.Int32(10),
		},
	}

	payload := &pbd.MatchReward{
		RoomId:  proto.Uint32(roomid),
		Result:  proto.Int32(1),
		VsAi:    proto.Int32(ai),
		Details: details,
	}
	err, resp := self.request(true, pbd.CG_ID_MATCH_REWARD, payload)
	if err != nil {
		log.Printf("匹配结算出错%d: %v", self.uid, err)
		return err
	}

	packet := resp.(*pbd.MatchRewardResp)

	if e := packet.GetErr(); e != nil {
		log.Printf("匹配结算失败%d: %v", self.uid, e)
		return self.Error(e)
	}

	return nil
}

func (self *Client) MatchEnd() {
	payload := &pbd.MatchEnd{
		RoomCode: proto.Uint32(uint32(rand.Intn(10000))),
	}
	err, resp := self.request(true, pbd.CG_ID_MATCH_END, payload)
	if err != nil {
		log.Printf("取消匹配出错%d: %v", self.uid, err)
		return
	}

	packet := resp.(*pbd.MatchEndResp)

	if e := packet.GetErr(); e != nil {
		log.Printf("取消匹配失败%d: %v", self.uid, e)
	}
}

// ---------------------------------------------------------------------------

// 请求rank匹配
func (self *Client) RankBegin() {
	payload := &pbd.RankBegin{
		Skin: proto.Int32(1),
	}
	err, resp := self.request(true, pbd.CG_ID_RANK_BEGIN, payload)
	if err != nil {
		log.Printf("发起RankBegin请求时出错: %v", err)
		return
	}
	packet := resp.(*pbd.RankBeginResp)
	_ = packet
	// log.Println(self.uid, "排位开始", packet)
	// for packet.GetStatus() == 4 {
	// 	time.Sleep(time.Minute)
	// 	err, resp = self.request(true, pbd.CG_ID_RANK_BEGIN, payload)
	// 	packet = resp.(*pbd.RankBeginResp)
	// }
}

func (self *Client) RankCancel() {
	payload := &pbd.RankCancel{}
	err, resp := self.request(true, pbd.CG_ID_RANK_CANCEL, payload)
	if err != nil {
		log.Printf("发起RankBegin请求时出错: %v", err)
		return
	}
	packet := resp.(*pbd.RankCancelResp)
	// log.Println(self.uid, "排位排队 取消", packet)
	_ = packet
}

// rank匹配查询
func (self *Client) RankQuery() {
	payload := &pbd.RankQuery{}

	err, resp := self.request(true, pbd.CG_ID_RANK_QUERY, payload)
	if err != nil {
		log.Printf("询问Rank匹配是否成功出错:uid= %v  err= %v", self.uid, err)
		return
	}
	// time.Sleep(time.Second)
	packet := resp.(*pbd.RankQueryResp)
	//log.Println(self.uid, "查询排位返回结果", packet)
	switch packet.GetStatus() {
	case 0:
		log.Println(self.uid, "RANK匹配.....")
		// self.RankQuery()
	case 1:
		// log.Println(self.uid, "RANK匹配超时")
	case 2:
		log.Println(self.uid, "RANK匹配成功")
		// self.RankReward(packet)
	}
}

//rank比赛结果
func (self *Client) RankReward(q *pbd.RankQueryResp) {
	roomUID := q.GetRoomUid()
	//result := roomUID%3 + 1
	payload := &pbd.RankReward{
		RoomUid: proto.Uint32(roomUID),
		Result:  proto.Int64(2),
	}

	err, resp := self.request(true, pbd.CG_ID_RANK_REWARD, payload)
	if err != nil {
		log.Printf("询问Rank比赛结果 出错:uid %v  %v", self.uid, err)
		return
	}

	packet := resp.(*pbd.RankRewardResp)
	// log.Println(self.uid, "排位比赛结果返回 :", packet)
	_ = packet
}

//rank赛季信息
func (self *Client) RankInfo() {
	payload := &pbd.RankInfo{}

	err, rsp := self.request(true, pbd.CG_ID_RANK_INFO, payload)
	if err != nil {
		log.Printf("询问赛季信息 出错:% %v", self.uid, err)
		return
	}

	packet := rsp.(*pbd.RankInfoResp)
	//log.Println(self.uid, "赛季信息返回 :", packet)
	if packet.GetStatus() != 0 {
		// time.Sleep(10 * time.Second)
	}
}

//自己的排位信息
func (self *Client) RankSelf() {
	payload := &pbd.RankSelf{}

	err, _ := self.request(true, pbd.CG_ID_RANK_SELF, payload)
	if err != nil {
		log.Printf("询问自己排位信息 出错:% %v", self.uid, err)
		return
	}

	//packet := rsp.(*pbd.RankSelfResp)
	//log.Println(self.uid, "排位信息 :", packet)
}

//获取排位段位表
func (self *Client) RankTiers() {
	payload := &pbd.RankTiers{}

	err, _ := self.request(true, pbd.CG_ID_RANK_TIERS, payload)
	if err != nil {
		log.Printf("询问排位段位信息 出错:% %v", self.uid, err)
		return
	}

	//packet := rsp.(*pbd.RankTiersResp)
	//log.Println(self.uid, "段位信息 : len =", len(packet.GetTableJson()))
}

//获取赛季奖励
func (self *Client) RankSeasonReward() {
	payload := &pbd.RankSeasonReward{}
	err, _ := self.request(true, pbd.CG_ID_RANK_SEASONREWARD, payload)
	if err != nil {
		log.Printf("获取赛季奖励 出错:%v %v", self.uid, err)
		return
	}

	//packet := rsp.(*pbd.RankSeasonRewardResp)
	//log.Println("获取赛季奖励", packet)
}

func (self *Client) SelfMoney() {
	err, resp := self.request(true, pbd.CG_ID_SELF_MONEY, nil)
	if err != nil {
		log.Printf("查询自己的货币时出错%d: %v", self.uid, err)
		return
	}

	packet := resp.(*pbd.SelfMoneyResp)
	_ = packet.GetMoney()
}

func (self *Client) ADReward() {
	err, resp := self.request(true, pbd.CG_ID_AD_REWARD, nil)
	if err != nil {
		log.Printf("领取广告奖励时出错%d: %v", self.uid, err)
		return
	}

	packet := resp.(*pbd.AdRewardResp)
	switch packet.GetErrorCode() {
	case 0:
	case 1:
	case 2:
	default:
		log.Printf("领取广告奖励时返回错误%d: %d", self.uid, packet.GetErrorCode())
	}
}
