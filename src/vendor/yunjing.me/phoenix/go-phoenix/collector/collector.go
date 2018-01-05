package collector

import (
	"time"
	// protodata "yunjing.me/phoenix/pbd/go/collect"
)

// 数据收集客户端
type Collector struct {
	client      *Client
	zoneId      uint32
	serviceType uint32
	serviceId   uint32
}

func New(addr string) *Collector {
	self := &Collector{
		client: newClient(addr),
	}

	return self
}

// 逻辑服调用
// 参数 happen[发生时间]
// 参数 id[事件ID]
// 参数 args[具体逻辑参数],例子:["abc",123,[]byte("abc")]
func (self *Collector) Do(happen time.Time, id uint32, args ...interface{}) error {
	raw, err := self.encode(self.zoneId, self.serviceType, self.serviceId, happen, id, args...)
	if err != nil {
		return err
	}

	self.client.Send(raw)

	return nil
}
