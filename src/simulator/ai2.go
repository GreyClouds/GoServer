package main

import (
	"log"
	"time"

	pbd "crazyant.com/deadfat/pbd/go"
)

type AI2 struct {
	state    int
	roomid   uint32
	ai       int32
	opponent uint32
}

func (self *AI2) Start(client *Client) {
	for {
		switch self.state {
		case 0:
			err := client.MatchBegin()
			if err != nil {
				println(err.Error())
				continue
			}

			self.state = 1
		case 1:
			err, isContinue, roomid, payload := client.MatchQuery()
			if err != nil {
				println(err.Error())
				continue
			}

			if isContinue {
				continue
			}

			p, ok := payload.(*pbd.MatchQueryResp)
			if !ok {
				// println("")
				continue
			}

			if roomid == 0 && p.GetVsAi() == 0 {
				log.Printf("匹配失败: %v", 1)
				self.state = 0
				continue
			}

			self.state = 2
			self.roomid = roomid
			self.ai = p.GetVsAi()
			self.opponent = p.GetOpponentUid()
		case 2:
			err := client.MatchReward(self.roomid, self.ai, self.opponent)
			if err != nil {
				log.Printf("MatchReward: %v", err)
				continue
			}

			self.state = 0
		}

		time.Sleep(2000 * time.Millisecond)
	}
}
