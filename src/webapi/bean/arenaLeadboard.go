package bean

import (
	"log"
	"github.com/astaxie/beego/orm"
	"time"
	"fmt"
	pkgProto "crazyant.com/deadfat/pbd/hero"
)

type ArenaLeaderboard struct {
	Uid     uint32 `orm:"pk;column(uid)"`  // 角色编号
	Score   uint32  `orm:"column(score)"`   // 车间波次
	Updated int64 `orm:"column(updated)"` //更新时间(当前时间)
}

func (self *ArenaLeaderboard) TableName() string {
	return "arenaLeaderboard"
}

// ---------------------------------------------------------------------------
var arenaLeaderboardCacheList []*ArenaLeaderboard = make([]*ArenaLeaderboard, 0)
//获取数据库中所有的ArenaLeaderboardtiaomu，并且做好排序，准备好排行榜服务
func ArenaLeaBoardLoadAndSort(){
	t := time.Now()
	//nowsec := t.Unix()

	var arenaLeaderboards []ArenaLeaderboard
	num , err := defaultOrm.QueryTable(new(ArenaLeaderboard)).Limit(-1).All(&arenaLeaderboards)
	if err != nil {
		panic(fmt.Sprintf("载入车间战斗排行榜缓存数据时出错: %v", err))
	}

	if(num == 0){
		return
	}
	for i := int64(0); i < num; i++ {
		item := arenaLeaderboards[i]
		arenaLeaderboardCacheList = append(arenaLeaderboardCacheList, &item)
	}
	highToLowsort(arenaLeaderboardCacheList, 0 , len(arenaLeaderboardCacheList) -1)

	long := time.Since(t)
	log.Println("bean.ArenaLeaBoardLoadAndSort 消耗时间:", long, num)
}


func DoUpdateArenaScore(uid uint32, score uint32) (bool, uint32){
	t := time.Now()

	var oldIndex int
	arenaLeaderboard := ArenaLeaderboard{Uid:uid}
	rErr := defaultOrm.Read(&arenaLeaderboard)
	cond := orm.NewCondition()

	if rErr == nil{
		cond1 := cond.And("score__gt", arenaLeaderboard.Score)
		cond2 := cond.And("score__eq", arenaLeaderboard.Score).And("updated__lt", arenaLeaderboard.Updated)
		cond3 := cond.And("score__eq", arenaLeaderboard.Score).And("updated__eq", arenaLeaderboard.Updated).And("uid__lt", arenaLeaderboard.Uid)

		lastCond := cond.AndCond(cond1).OrCond(cond2).OrCond(cond3)
		temIndex, qErr := defaultOrm.QueryTable(new(ArenaLeaderboard)).SetCond(lastCond).Count()
		checkError("UpdateArenaScore_queryOldIndex,错误:", qErr)
		oldIndex = int(temIndex)
		oldItem := arenaLeaderboardCacheList[oldIndex]
		oldItem.Score = score
		oldItem.Updated = time.Now().Unix()

		arenaLeaderboard.Score = score;
		arenaLeaderboard.Updated = time.Now().Unix()
		_, uErr := defaultOrm.Update(&arenaLeaderboard)
		checkError("UpdateArenaScore_update,错误:", uErr)
	}else if(rErr == orm.ErrNoRows){
		arenaLeaderboard.Score = score
		arenaLeaderboard.Updated = time.Now().Unix()
		oldIndex = len(arenaLeaderboardCacheList)
		arenaLeaderboardCacheList = append(arenaLeaderboardCacheList, &arenaLeaderboard)
		_, iErr := defaultOrm.Insert(&arenaLeaderboard)
		checkError("UpdateArenaScore_insert,错误:", iErr)
	}else{
		panic(fmt.Sprintf("DoUpdateArenaScore 获取玩家当前分数错误: %v", rErr))
	}

	cond1 := cond.And("score__gt", arenaLeaderboard.Score)
	cond2 := cond.And("score__eq", arenaLeaderboard.Score).And("updated__lt", arenaLeaderboard.Updated)
	cond3 := cond.And("score__eq", arenaLeaderboard.Score).And("updated__eq", arenaLeaderboard.Updated).And("uid__lt", arenaLeaderboard.Uid)

	lastCond := cond.AndCond(cond1).OrCond(cond2).OrCond(cond3)
	newIndex, qErr := defaultOrm.QueryTable(new(ArenaLeaderboard)).SetCond(lastCond).Count()
	checkError("UpdateArenaScore_queryNewIndex,错误:", qErr)

	if(oldIndex != int(newIndex)){
		arenaLeaderboardCacheList[oldIndex], arenaLeaderboardCacheList[newIndex] = arenaLeaderboardCacheList[newIndex], arenaLeaderboardCacheList[oldIndex]
	}
	long := time.Since(t)
	log.Println("bean.DoUpdateArenaScore 消耗时间:", long)
	return true, uint32(newIndex+1)
	}

	//包括start，不包括end
func GetArenaLeaderboard(start int, end int) []*pkgProto.LearderboardInfo {
	lastRank := len(arenaLeaderboardCacheList) + 1
	if(end > lastRank){
		end = lastRank
	}
	if(start > end){
		return make([]*pkgProto.LearderboardInfo, 0)
	}
	answerList := make([]*pkgProto.LearderboardInfo, 0, end-start)
	temArenaLeaderboardList := arenaLeaderboardCacheList[start-1:end-1]
	for i, item := range temArenaLeaderboardList{
		ch := SimRole{Uid: item.Uid}
		defaultOrm.Read(&ch)
		name := ch.Name
		answerList = append(answerList, &pkgProto.LearderboardInfo{int32(item.Score), uint32(start+i), name})
	}
	return answerList
}


func GetRank(uid uint32) uint32 {
	t := time.Now()

	arenaLeaderboard := ArenaLeaderboard{Uid:uid}
	rErr := defaultOrm.Read(&arenaLeaderboard)
	if rErr == nil{
		cond := orm.NewCondition()
		cond1 := cond.And("score__gt", arenaLeaderboard.Score)
		cond2 := cond.And("score__eq", arenaLeaderboard.Score).And("updated__lt", arenaLeaderboard.Updated)
		cond3 := cond.And("score__eq", arenaLeaderboard.Score).And("updated__eq", arenaLeaderboard.Updated).And("uid__lt", arenaLeaderboard.Uid)

		lastCond := cond.AndCond(cond1).OrCond(cond2).OrCond(cond3)

		log.Println("GetRank 读取出来的数据:", arenaLeaderboard.Score, arenaLeaderboard.Updated)

		temIndex, qErr := defaultOrm.QueryTable(new(ArenaLeaderboard)).SetCond(lastCond).Count()
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

func GetScore(uid uint32) uint32 {
	arenaLeaderboard := ArenaLeaderboard{Uid:uid}
	rErr := defaultOrm.Read(&arenaLeaderboard)
	if rErr == nil{
		return arenaLeaderboard.Score
	}else if(rErr == orm.ErrNoRows){

		return 0
	}else{
		panic(fmt.Sprintf("DoUpdateArenaScore 获取玩家当前分数错误: %v", rErr))
	}

}

//排序：积分从高到底，积分相同，updated小的靠前
func highToLowsort(arr []*ArenaLeaderboard, start int, end int) {
	var (
		key  *ArenaLeaderboard = arr[start]
		low  int               = start
		high int               = end
	)
	for {
		for low < high {
			if arr[high].Score > key.Score ||
				(arr[high].Score == key.Score && arr[high].Updated < key.Updated) ||
				(arr[high].Score == key.Score && arr[high].Updated < key.Updated && arr[high].Uid < key.Uid){
				arr[low] = arr[high]
				break
			}
			high--
		}
		for low < high {
			if arr[low].Score < key.Score ||
				(arr[low].Score == key.Score && arr[low].Updated > key.Updated) ||
				(arr[low].Score == key.Score && arr[low].Updated > key.Updated && arr[low].Uid > key.Uid){
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
		highToLowsort(arr, start, low-1)
	}
	if high+1 < end {
		highToLowsort(arr, high+1, end)
	}
}