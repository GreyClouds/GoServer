package bean

import (
"log"
"github.com/astaxie/beego/orm"
"time"
"fmt"
pkgProto "crazyant.com/deadfat/pbd/hero"
)

type ChallLeaderboard struct {
	Uid     uint32 `orm:"pk;column(uid)"`  // 角色编号
	Score   uint64  `orm:"column(score)"`   // 挑战花费时间
	Updated int64 `orm:"column(updated)"` //更新时间(当前时间)
}

func (self *ChallLeaderboard) TableName() string {
	return "challengeLeaderboard"
}

// ---------------------------------------------------------------------------
var challLeaderboardCacheList []*ChallLeaderboard = make([]*ChallLeaderboard, 0)
//获取数据库中所有的tiaomu，并且做好排序，准备好排行榜服务
func ChallLeaBoardLoadAndSort(){
	t := time.Now()
	var challengeLeaderboards []ChallLeaderboard
	num , err := defaultOrm.QueryTable(new(ChallLeaderboard)).Limit(-1).All(&challengeLeaderboards)
	if err != nil {
		panic(fmt.Sprintf("载入特殊挑战排行榜缓存数据时出错: %v", err))
	}

	if(num == 0){
		return
	}
	for i := int64(0); i < num; i++ {
		item := challengeLeaderboards[i]
		challLeaderboardCacheList = append(challLeaderboardCacheList, &item)
	}
	lowToHighSort(challLeaderboardCacheList, 0 , len(challLeaderboardCacheList) -1)

	long := time.Since(t)
	log.Println("bean.ChallLeaBoardLoadAndSort 消耗时间:", long, num)
}


func DoUpdateChallScore(uid uint32, score uint64) (bool, uint32){
	t := time.Now()

	var oldIndex int
	challLeaderboard := ChallLeaderboard{Uid:uid}
	rErr := defaultOrm.Read(&challLeaderboard)
	cond := orm.NewCondition()
	if rErr == nil{
		cond1 := cond.And("score__lt", challLeaderboard.Score)
		cond2 := cond.And("score__eq", challLeaderboard.Score).And("updated__lt", challLeaderboard.Updated)
		cond3 := cond.And("score__eq", challLeaderboard.Score).And("updated__eq", challLeaderboard.Updated).And("uid__lt", challLeaderboard.Uid)

		lastCond := cond.AndCond(cond1).OrCond(cond2).OrCond(cond3)
		temIndex, qErr := defaultOrm.QueryTable(new(ChallLeaderboard)).SetCond(lastCond).Count()
		checkError("UpdateArenaScore_queryOldIndex,错误:", qErr)
		oldIndex = int(temIndex)
		oldItem := challLeaderboardCacheList[oldIndex]
		oldItem.Score = score
		oldItem.Updated = time.Now().Unix()

		challLeaderboard.Score = score;
		challLeaderboard.Updated = time.Now().Unix()
		_, uErr := defaultOrm.Update(&challLeaderboard)
		checkError("DoUpdateChallScore_update,错误:", uErr)
	}else if(rErr == orm.ErrNoRows){
		challLeaderboard.Score = score
		challLeaderboard.Updated = time.Now().Unix()
		oldIndex = len(challLeaderboardCacheList)
		challLeaderboardCacheList = append(challLeaderboardCacheList, &challLeaderboard)
		_, iErr := defaultOrm.Insert(&challLeaderboard)
		checkError("DoUpdateChallScore_insert,错误:", iErr)
	}else{
		panic(fmt.Sprintf("DoUpdateChallScore 获取玩家当前分数错误: %v", rErr))
	}

	cond1 := cond.And("score__lt", challLeaderboard.Score)
	cond2 := cond.And("score__eq", challLeaderboard.Score).And("updated__lt", challLeaderboard.Updated)
	cond3 := cond.And("score__eq", challLeaderboard.Score).And("updated__eq", challLeaderboard.Updated).And("uid__lt", challLeaderboard.Uid)

	lastCond := cond.AndCond(cond1).OrCond(cond2).OrCond(cond3)

	newIndex, qErr := defaultOrm.QueryTable(new(ChallLeaderboard)).SetCond(lastCond).Count()
	checkError("DoUpdateChallScore_queryNewIndex,错误:", qErr)

	if(oldIndex != int(newIndex)){
		challLeaderboardCacheList[oldIndex], challLeaderboardCacheList[newIndex] = challLeaderboardCacheList[newIndex], challLeaderboardCacheList[oldIndex]
	}
	long := time.Since(t)
	log.Println("bean.DoUpdateArenaScore 消耗时间:", long)
	return true, uint32(newIndex+1)
}

//包括start，不包括end
func GetChallLeaderboard(start int, end int) []*pkgProto.LearderboardInfo {
	lastRank := len(challLeaderboardCacheList) + 1
	if(end > lastRank){
		end = lastRank
	}
	if(start > end){
		return make([]*pkgProto.LearderboardInfo, 0)
	}
	answerList := make([]*pkgProto.LearderboardInfo, 0, end-start)
	temArenaLeaderboardList := challLeaderboardCacheList[start-1:end-1]
	for i, item := range temArenaLeaderboardList{
		ch := SimRole{Uid: item.Uid}
		defaultOrm.Read(&ch)
		name := ch.Name
		answerList = append(answerList, &pkgProto.LearderboardInfo{int32(item.Score), uint32(start+i), name})
	}
	return answerList
}


func GetChallRank(uid uint32) uint32 {
	t := time.Now()

	challLeaderboard := ChallLeaderboard{Uid:uid}
	rErr := defaultOrm.Read(&challLeaderboard)
	if rErr == nil{
		cond := orm.NewCondition()
		cond1 := cond.And("score__lt", challLeaderboard.Score)
		cond2 := cond.And("score__eq", challLeaderboard.Score).And("updated__lt", challLeaderboard.Updated)
		cond3 := cond.And("score__eq", challLeaderboard.Score).And("updated__eq", challLeaderboard.Updated).And("uid__lt", challLeaderboard.Uid)

		lastCond := cond.AndCond(cond1).OrCond(cond2).OrCond(cond3)
		log.Println("GetRank 读取出来的数据:", challLeaderboard.Score, challLeaderboard.Updated)

		temIndex, qErr := defaultOrm.QueryTable(new(ChallLeaderboard)).SetCond(lastCond).Count()
		checkError("UpdateArenaScore_queryOldIndex,错误:", qErr)

		long := time.Since(t)
		log.Println("bean.GetRank 消耗时间:", long)

		return uint32(temIndex+1)
	}else if(rErr == orm.ErrNoRows){
		return 0
	}else{
		panic(fmt.Sprintf("DoUpdateArenaScore 获取玩家当前分数错误: %v", rErr))
	}

}

func GetChallScore(uid uint32) uint64 {
	challLeaderboard := ChallLeaderboard{Uid:uid}
	rErr := defaultOrm.Read(&challLeaderboard)
	if rErr == nil{
		return challLeaderboard.Score
	}else if(rErr == orm.ErrNoRows){
		return 0
	}else{
		panic(fmt.Sprintf("DoUpdateArenaScore 获取玩家当前分数错误: %v", rErr))
	}
}

//排序：积分从低到高，积分相同，updated小的靠前
func lowToHighSort(arr []*ChallLeaderboard, start int, end int) {
	var (
		key  *ChallLeaderboard = arr[start]
		low  int               = start
		high int               = end
	)
	for {
		for low < high {
			if arr[high].Score < key.Score ||
				(arr[high].Score == key.Score && arr[high].Updated < key.Updated) ||
				(arr[high].Score == key.Score && arr[high].Updated == key.Updated && arr[high].Uid < key.Uid){
				arr[low] = arr[high]
				break
			}
			high--
		}
		for low < high {
			if arr[low].Score > key.Score ||
				(arr[low].Score == key.Score && arr[low].Updated > key.Updated) ||
				(arr[low].Score == key.Score && arr[low].Updated == key.Updated && arr[low].Uid > key.Uid){
				arr[high] = arr[low]
				break
			}
			low++
		}
		if low >= high {
			arr[low] = key
			break
		}
	}
	if low-1 > start {
		lowToHighSort(arr, start, low-1)
	}
	if high+1 < end {
		lowToHighSort(arr, high+1, end)
	}
}