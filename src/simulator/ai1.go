package main

import (
	"time"
)

type AI1 struct {
}

const t = time.Millisecond * 300

func (self *AI1) Start(client *Client) {
	for {
		time.Sleep(t)
		client.RankInfo()
		time.Sleep(t)
		client.RankSeasonReward()
		time.Sleep(t)
		client.RankTiers()
		time.Sleep(t)
		client.RankSelf()
		time.Sleep(t)
		client.RankBegin()
		time.Sleep(t)
		//client.RankCancel()
		time.Sleep(t)
		client.RankQuery()
	}

	time.Sleep(time.Second)
}
