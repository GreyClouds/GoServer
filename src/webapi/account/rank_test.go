package account

import (
	"log"
	"strconv"
	"testing"
)

func TestRankListManager_Rank(t *testing.T) {
	for i := 0; i < 10; i++ {
		log.Println(DefaultRankList().Rank(uint32(i), "玩家"+strconv.Itoa(i), int32(10-i)))
	}
	log.Println(DefaultRankList().Rank(123, "玩家"+strconv.Itoa(123), int32(11)))
	log.Println(DefaultRankList().Rank(2, "玩家"+strconv.Itoa(2), int32(6)))
	log.Println(DefaultRankList().Rank(3, "玩家"+strconv.Itoa(3), int32(60)))

	log.Println("======", DefaultRankList().rank)
	DefaultRankList().user.Range(func(key, value interface{}) bool {
		log.Println(key, value)
		return true
	})
}
