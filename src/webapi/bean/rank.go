package bean

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/astaxie/beego/orm"
)

const (
	NoResult = iota
	SeatLeft
	SeatRight
)

const (
	up = iota + 1
	ins
	del
)

type RankData struct {
	Season int32 `orm:"column(season)"`
	Tiers  int32 `orm:"column(tiers)"`
	Max    int32 `orm:"column(max)"`
	Win    int32 `orm:"column(win)"`
	Lose   int32 `orm:"column(lose)"`
	Draw   int32 `orm:"column(draw)"`
}

type Rank struct {
	op     int        `orm:"-"`
	Status int        `orm:"-"`
	Mtx    sync.Mutex `orm:"-"`
	UID    uint32     `orm:"column(uid);pk"`
	RankData
	Updated  time.Time `orm:"auto_now;type(datetime);column(updated)"`
	Suspend  int64     `orm:"column(suspend)"`
	Check    int64     `orm:"column(check)"`
	ResetErr int64     `orm:"column(reset_err)"`
	Cheat    int       `orm:"column(cheat)"`
	Sum      int32     `orm:"column(sum)"`
	Err      int32     `orm:"column(err)"`
}

var defaultOrm orm.Ormer

func DefaultORM() {
	defaultOrm = orm.NewOrm()
}

func NewRank(uid uint32) *Rank {
	return &Rank{
		UID: uid,
	}
}

func (r *Rank) Update(field ...string) *Rank {
	_, err := defaultOrm.Update(r, field...)
	checkError("bean/rank.go:func (r *Rank) Update(field ...string) *Rank", err)
	return r
}

func (r *Rank) Read() *Rank {
	err := defaultOrm.Read(r)
	checkError("bean/rank.go:func (r *Rank) Read() *Rank", err)
	return r
}

func (r *Rank) clone() *Rank {
	c := *r
	return &c
}

func (r *Rank) SumAdd() *Rank {
	r.Mtx.Lock()
	r.Sum++
	r.Mtx.Unlock()
	return r
}

func LoadSeasonRank() []*Rank {
	var r []*Rank
	_, err := defaultOrm.QueryTable("rank").Limit(-1).All(&r)
	checkError("bean/rank.go:LoadSeasonRank  ", err)
	return r
}

type RankHistory struct {
	ID  int64  `orm:"column(id);pk;auto"`
	UID uint32 `orm:"column(uid);index"`
	RankData
	Updated time.Time `orm:"auto_now;type(datetime);column(updated)"`
}

type RankReward struct {
	ID      int64             `orm:"column(id);pk;auto"`
	UID     uint32            `orm:"column(uid);index"`
	Season  int32             `orm:"column(season)"`
	Tiers   int32             `orm:"column(tiers)"`
	Reward  []*RankRewardItem `orm:"reverse(many)"`
	Updated time.Time         `orm:"auto_now;type(datetime);column(updated)"`
}

func (r *RankReward) Rewards() map[int32]int32 {
	results := make(map[int32]int32)

	l := len(r.Reward)
	for i := 0; i < l; i++ {
		v := r.Reward[i]

		results[v.RewardID] = v.RewardValue
	}

	return results
}

type RankRewardItem struct {
	ID          int64       `orm:"column(id);pk;auto"`
	RankReward  *RankReward `orm:"rel(fk)"`
	RewardID    int32       `orm:"column(reward_id)"`
	RewardValue int32       `orm:"column(reward_value)"`
	Updated     time.Time   `orm:"auto_now;type(datetime);column(updated)"`
}

func (r *RankReward) Delete() {
	defaultOrm.Delete(r)
	l := len(r.Reward)
	for i := 0; i < l; i++ {
		defaultOrm.Delete(r.Reward[i])
	}
}

func (r *RankReward) RelatedReward() *RankReward {
	_, err := defaultOrm.QueryTable("rank_reward_item").
		Filter("RankReward", r.ID).RelatedSel().All(&r.Reward)
	checkError("bean/rank.go:func (r *RankReward) RelatedReward() *RankReward", err)
	return r
}

func checkError(i string, err error) {
	if nil != err && err.Error() != "<QuerySeter> no row found" {
		log.Println("bean/rank.go:", i, " : ", err.Error())
	}
}

func LoadUserSeasonReward(uid uint32) []*RankReward {
	var r []*RankReward
	_, err := defaultOrm.QueryTable("rank_reward").Filter("UID", uid).All(&r)
	checkError("bean/rank.go:func LoadUserSeasonReward(uid uint32) []*RankReward", err)
	l := len(r)
	for i := 0; i < l; i++ {
		r[i].RelatedReward()
	}
	return r
}

type RankRoom struct {
	ID              uint32    `orm:"column(id);pk"`
	RoomID          uint32    `orm:"column(room_id)"`
	BeginTime       time.Time `orm:"type(datetime);column(begin_time)"`
	Season          int32     `orm:"column(season)"`
	UID             uint32    `orm:"column(uid);"`
	RivalUID        uint32    `orm:"column(rival_uid)"`
	RivalRankRoomID uint32    `orm:"column(rival_rank_room_id)"`
	Skin            int32     `orm:"column(skin)"`
	Result          int64     `orm:"column(result)"`
	Seat            int       `orm:"column(seat)"`
	BeforeTiers     int32     `orm:"column(before_tiers)"`
	AfterTiers      int32     `orm:"column(after_tiers)"`
	Flag            int64     `orm:"column(flag)"`
	Updated         time.Time `orm:"auto_now;type(datetime);column(updated)"`
	op              int       `orm:"-"`
}

type RankList struct {
	op      int       `orm:"-"`
	UID     uint32    `orm:"column(uid);pk"`
	Rank    int       `orm:"column(rank);index"`
	Nick    string    `orm:"column(nick)"`
	Tiers   int32     `orm:"column(tiers)"`
	Updated time.Time `orm:"auto_now;type(datetime);column(updated)"`
}

func NewRankList(uid uint32, nick string, tiers int32, r int) *RankList {
	return &RankList{
		UID:   uid,
		Nick:  nick,
		Tiers: tiers,
		Rank:  r + 1,
	}
}

func (r *RankList) clone() *RankList {
	c := *r
	return &c
}

func (r *RankList) Update() *RankList {
	c := r.clone()
	c.op = up
	DefaultKeeper().SendRankList(c)
	return r
}

func (r *RankList) Insert() *RankList {
	c := r.clone()
	c.op = ins
	DefaultKeeper().SendRankList(c)

	return r
}

func (r *RankList) Delete() *RankList {
	c := r.clone()
	c.op = del
	DefaultKeeper().SendRankList(c)
	return r
}

func LoadRankList() []*RankList {
	var r []*RankList
	_, err := defaultOrm.QueryTable("rank_list").All(&r)
	checkError("bean/rank.go:func LoadRankList() []*RankList", err)
	return r
}

func RankListTruncate() {
	_, err := defaultOrm.QueryTable("rank_list").Exclude("UID", -1).Delete()
	checkError("清空排行榜表记录", err)
}

func RankResultCleanTask(season int32) {
	query := fmt.Sprintf("DELETE FROM rank_room WHERE season < %v", season)
	f := func() {
		res, err := defaultOrm.Raw(query).Exec()
		if err != nil {
			log.Printf("清理历史排位房间结算时出错: %v", err)
			return
		}
		i, _ := res.RowsAffected()
		log.Printf("清理历史排位房间结算共 %d 条", i)
	}
	DefaultKeeper().SendCron(Once, f)
}
