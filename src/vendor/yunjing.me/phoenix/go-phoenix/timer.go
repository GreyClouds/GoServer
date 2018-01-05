package phoenix

type internalTimer struct {
	id       uint32
	holder   interface{}
	delay    int
	interval int
	callback func()
	end      bool
	pos      int
}

func newTimer(id uint32, holder interface{}, delay, interval int, callback func()) *internalTimer {
	return &internalTimer{
		id:       id,
		holder:   holder,
		delay:    delay,
		interval: interval,
		callback: callback,
		end:      false,
		pos:      0,
	}
}

func (self *internalTimer) OnTick() {
	if self.end {
		return
	}

	// if self.delay >= self.pos {
	// 	self.delay - self.pos
	// }

	// self.pos++

	// if self.delay > 0 {
	// 	self.delay--
	// 	if self.delay == 0 {
	// 		(self.callback)()
	// 		if self.interval == 0 {
	// 			self.end = true
	// 		}
	// 	}

	// 	return
	// }
}

func (self *internalTimer) IsEnd() bool {
	return self.end
}
