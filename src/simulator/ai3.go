package main

import (
	"math/rand"
	"time"
)

type AI3 struct {
}

func (self *AI3) Start(client *Client) {
	for {
		switch rand.Intn(13 * 10) {
		case 0:
			client.InviteCode()
		case 1:
			client.RankInfo()
		case 2:
			client.RankSeasonReward()
		case 3:
			client.RankTiers()
		case 4:
			client.RankSelf()
		case 5:
			client.RankBegin()
		case 6:
			client.RankCancel()
		case 7:
			client.RankQuery()
		case 8:
			client.MatchBegin()
		case 9:
			client.MatchQuery()
		case 10:
			client.MatchEnd()
		case 11:
			client.ViewSignTask()
		case 12:
			client.SelfMoney()
		case 13:
			client.CDKey()
		case 14:
			client.ADReward()
		default:

		}

		time.Sleep(time.Second)
	}
}
