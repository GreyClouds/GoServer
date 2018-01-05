package account

import (
	"sync"
	"time"

	"webapi/bean"
	"webapi/config"
)

const (
	RankTiersKing = 0x60100
	rankListNum   = 100
	UnRanked      = 0
	TiersMark     = 0xFFFFFF
)

const (
	DRAW = iota + 1
	WIN
	LOSE
)

func RankCLear(r *Role, result int64) *Role {
	switch int(result) {
	case DRAW:
		RankCLearDraw(r)
	case WIN:
		RankCLearWin(r)
	case LOSE:
		RankCLearLose(r)
	default:

	}
	return r
}

func RankCLearDraw(r *Role) {
	r.rank.Draw++
	bean.DefaultKeeper().SendRankUpdate(r.rank)
}

func RankCLearWin(r *Role) {
	r.character.WinMatch()
	r.ClearWinReward()
	tiers, ok := config.Conf().GetRankTiers(r.rank.Tiers)
	if ok {
		if (r.rank.Tiers & 0xFFFFFF) < tiers.Up.Lv {
			if RankTiersKing <= tiers.Up.Lv {
				num := defaultRankListMgr.Rank(r.rank.UID, r.NickName(), tiers.Up.Lv)
				if num > 0 {
					r.rank.Tiers = tiers.Up.Lv | (int32(num) << 24)
				} else {
					r.rank.Tiers = tiers.Up.Lv
				}
			} else {
				r.rank.Tiers = tiers.Up.Lv
			}
			if r.rank.Tiers > r.rank.Max {
				r.rank.Max = r.rank.Tiers
			}
		}
	}
	r.rank.Win++
	bean.DefaultKeeper().SendRankUpdate(r.rank)
}

func RankCLearLose(r *Role) {
	r.character.LoseMatch()
	tiers, ok := config.Conf().GetRankTiers(r.rank.Tiers)
	if ok && (r.rank.Tiers&TiersMark) > tiers.Down.Lv {
		if r.rank.Tiers>>24 > 0 {
			num := defaultRankListMgr.Rank(r.rank.UID, r.NickName(), tiers.Down.Lv)
			if num > 0 {
				r.rank.Tiers = tiers.Down.Lv | (int32(num) << 24)
			} else {
				r.rank.Tiers = tiers.Down.Lv
			}
		} else {
			r.rank.Tiers = tiers.Down.Lv
		}
	}
	r.rank.Lose++
	bean.DefaultKeeper().SendRankUpdate(r.rank)
}

type RankListManager struct {
	mtx  sync.Mutex
	rank []*bean.RankList
	user sync.Map
}

var defaultRankListMgr = NewRankListManager()

func NewRankListManager() *RankListManager {
	return &RankListManager{
		rank: make([]*bean.RankList, 0, rankListNum),
	}
}

func (r *RankListManager) GetRankList(uid uint32) (*bean.RankList, bool) {
	v, ok := r.user.Load(uid)
	if ok {
		l, yes := v.(*bean.RankList)
		return l, yes
	} else {
		return nil, false
	}
}

func (r *RankListManager) Load() *RankListManager {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	rankList := bean.LoadRankList()
	l := len(rankList)
	r.rank = make([]*bean.RankList, l)
	for i := 0; i < l; i++ {
		temp := rankList[i]
		r.rank[temp.Rank-1] = temp
		r.user.Store(temp.UID, temp)
	}
	return r
}

func (r *RankListManager) Clean() *RankListManager {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	r.rank = make([]*bean.RankList, 0, rankListNum)
	r.user = sync.Map{}
	bean.RankListTruncate()
	return r
}

func (r *RankListManager) insertRankList(l *bean.RankList) {
	r.user.Store(l.UID, l)
	l.Insert()
}

func (r *RankListManager) deleteRankList(l *bean.RankList) {
	r.user.Delete(l.UID)
	l.Delete()
}

func (r *RankListManager) RankOld(old *bean.RankList, t int32) int {
	oldTiers := old.Tiers
	old.Tiers = t
	if oldTiers > t {
		r.rankListDown(old.Rank - 1)
	} else if oldTiers < t {
		r.rankListUp(old.Rank - 1)
	}
	return old.Rank
}

func (r *RankListManager) RankNew(rankList *bean.RankList) int {
	r.rank = append(r.rank, rankList)
	r.insertRankList(rankList)
	r.rankListUp(len(r.rank) - 1)
	return rankList.Rank
}

func (r *RankListManager) RankFoot(rankList *bean.RankList) int {
	pos := rankListNum - 1
	if r.rank[pos].Tiers < rankList.Tiers {
		r.deleteRankList(r.rank[pos])
		r.insertRankList(rankList)
		r.rank[pos] = rankList
		r.rankListUp(pos)
		return rankList.Rank
	} else {
		return UnRanked
	}
}

func (r *RankListManager) Rank(uid uint32, nick string, tiers int32) int {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	old, ok := r.GetRankList(uid)
	if ok {
		return r.RankOld(old, tiers)
	}
	rankList := bean.NewRankList(uid, nick, tiers, len(r.rank))
	if len(r.rank) >= rankListNum {
		return r.RankFoot(rankList)
	} else {
		return r.RankNew(rankList)
	}
}

func (r *RankListManager) rankListUp(n int) {
	p := r.rank[n]
	for i := n - 1; i >= 0; i-- {
		temp := r.rank[i]
		if temp.Tiers < p.Tiers {
			temp.Rank++
			r.rank[i+1] = temp
			temp.Update()
		} else {
			p.Rank = temp.Rank + 1
			r.rank[i+1] = p
			p.Update()
			return
		}
	}
	p.Rank = 1
	r.rank[0] = p
	p.Update()
}

func (r *RankListManager) rankListDown(n int) {
	p := r.rank[n]
	l := len(r.rank)
	for i := n + 1; i < l; i++ {
		temp := r.rank[i]
		if temp.Tiers > p.Tiers {
			temp.Rank--
			r.rank[i-1] = temp
			temp.Update()
		} else {
			p.Rank = temp.Rank - 1
			r.rank[i-1] = p
			p.Update()
			return
		}
	}
	p.Rank = l
	r.rank[l-1] = p
	p.Update()
}

func DefaultRankList() *RankListManager {
	return defaultRankListMgr
}

func RankCheatClean(rank *bean.Rank) {
	now := time.Now().Unix()
	rank.Mtx.Lock()
	if now >= rank.Suspend {
		rank.Suspend = 0
	}

	forbid := config.Conf().GetRankForbid(rank.Cheat)
	if now >= rank.ResetErr {
		rank.Sum = 0
		rank.Err = 0
		rank.ResetErr = now + int64(forbid.ResetTime)
	}

	if now >= rank.Check {
		t := now - rank.Check
		for t > 0 && rank.Cheat > 0 {
			rank.Cheat--
			next := config.Conf().GetRankForbid(rank.Cheat)
			t = t - int64(next.ProbationTime)
		}
		cur := config.Conf().GetRankForbid(rank.Cheat)
		rank.Check = now + int64(cur.ProbationTime)
	}
	bean.DefaultKeeper().SendRankUpdate(rank)
	rank.Mtx.Unlock()
}

func RankCheatAdd(rank *bean.Rank) {
	rank.Mtx.Lock()
	rank.Err++
	if config.Conf().GetRankCheatJudge(rank.Sum, rank.Err) {
		rank.Cheat++
		now := time.Now().Unix()
		forbid := config.Conf().GetRankForbid(rank.Cheat)
		rank.Suspend = now + int64(forbid.ForbidTime)
		rank.Sum = 0
		rank.Err = 0
		rank.ResetErr = now + int64(forbid.ResetTime)
		rank.Check = rank.Suspend + int64(forbid.ProbationTime)
	}
	bean.DefaultKeeper().SendRankUpdate(rank)
	rank.Mtx.Unlock()
}
