package main

import (
	"fmt"
	"log"

	"math/rand"

	"gopkg.in/go-playground/pool.v3"
)

type Simulator struct {
	addr string // 服务器连接地址
}

func NewSimulator(addr string) *Simulator {
	s := &Simulator{
		addr: addr,
	}

	return s
}

func (s *Simulator) Start(concurrence int) {
	p := pool.NewLimited(uint(concurrence))
	defer p.Close()

	batch := p.Batch()

	go func() {
		for i := 0; i < concurrence; i++ {
			batch.Queue(doClientNewAndStart(s.addr, int(rand.Int63n(100000))))
		}

		batch.QueueComplete()
	}()

	for v := range batch.Results() {
		if err := v.Error(); err != nil {
			// handle errorStart
			// maybe call batch.Cancel()

			log.Printf("任务执行中出错: %v", err)
			continue
		}

		// use result value
		log.Printf("%v", v.Value().(bool))
	}
}

func doClientNewAndStart(addr string, i int) pool.WorkFunc {
	return func(wu pool.WorkUnit) (interface{}, error) {
		imei := fmt.Sprintf("xxxx:%d", i)
		client := newClient(addr, imei)
		ok := client.Login("100000", "1.5.1")

		if !ok {
			log.Printf("设备登陆失败: %s", imei)
			return false, nil
		}

		// log.Printf("设备登陆成功: %s", imei)

		if wu.IsCancelled() {
			log.Printf("设备任务被取消: %s", imei)
			return true, nil
		}
		new(AI2).Start(client)

		return true, nil
	}
}
