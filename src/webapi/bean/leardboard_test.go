package bean

import (
	"log"
	"testing"
	"fmt"
)

func Test_Leardboard(t *testing.T){
	DefaultORM()
	//ClearArenaLeaderboard()
	WriteInDDB()
	//ArenaLeaBoardLoadAndSort()
	//log.Println("获得9000的排名：", GetRank(9000))
	//getInfoList := GetArenaLeaderboard(0,50)
	//for _, item := range getInfoList{
	//	log.Println("获取指定排名数据：", item.score, item.name, item.rank)
	//}
}

//2017/12/07 14:41:41 bean.ArenaLeaBoardLoadAndSort 消耗时间: 4.5419961s 100000
//2017/12/07 14:41:41 bean.DoUpdateArenaScore 消耗时间: 50.1708ms
//2017/12/07 14:41:41 GetRank 读取出来的数据: 10000 1512628901
//2017/12/07 14:41:41 bean.GetRank 消耗时间: 24.0633ms

func WriteInDDB(){
	dateNum := 90000
	for i := dateNum; i > 1000 ; i--{
		name := fmt.Sprintf("testName%v", i);
		arenaLeaderboard := ArenaLeaderboard{uint32(i), int32(i), int64(i)}
		simRole := SimRole{name, uint32(i)}
		_, iErr := defaultOrm.Insert(&arenaLeaderboard)
		_, iErr2 := defaultOrm.Insert(&simRole)
		checkError("WriteInDDB,错误:", iErr)
		checkError("WriteInDDB,错误2:", iErr2)
	}
	log.Println("数据写入完成：", dateNum)
}

func ClearArenaLeaderboard(){
	defaultOrm.Raw("DELETE From arenaLeaderboard").Exec()
	defaultOrm.Raw("DELETE From arenaLeaderboard").Exec()
	cc, _ := defaultOrm.QueryTable(new(ArenaLeaderboard)).Exclude("Uid", -1).Count()
	log.Println("清空排行榜--剩余记录条数:", cc)

}

//func Test_LeardoardUpdate(t *testing.T){
//	DefaultORM()
//	ArenaLeaBoardLoadAndSort()
//	DoUpdateArenaScore(9000, 20000)
//	log.Println("获得9000的排名：", GetRank(9000))
//		getInfoList := GetArenaLeaderboard(1,10)
//		for _, item := range getInfoList{
//			log.Println("获取指定排名数据：", item.Score, item.Name, item.Rank)
//		}
//}
