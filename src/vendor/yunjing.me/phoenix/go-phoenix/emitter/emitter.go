package emitter

const (
	EVENT_REGISTER_LISTENER   = EventID(1) // 事件监听注册
	EVENT_UNREGISTER_LISTENER = EventID(2) // 事件监听移除
)

type EventID uint32

type Event struct {
	ID   EventID
	Args []interface{}
}

// Emitter - listeners container
type Emitter struct {
	queue     chan<- *Event
	listeners map[EventID][]Listener
}

// Listener - callback container and whether it will run once or not
type Listener struct {
	callback func(...interface{})
	once     bool
}

// NewEventEmitter() - create a new instance of Emitter
func NewEventEmitter(queue chan<- *Event) *Emitter {
	return &Emitter{
		queue:     queue,
		listeners: make(map[EventID][]Listener),
	}
}

// Finialize() - free memory from an emitter instance
func (self *Emitter) Finialize() {
	self.listeners = nil
}

func (self *Emitter) GetEventListenerList(id EventID) []Listener {
	arr, exists := self.listeners[id]
	if !exists {
		return nil
	}

	return arr
}

func (self *Emitter) GetEventListenerListCount(id EventID) int {
	listeners := self.GetEventListenerList(id)
	if listeners == nil {
		return 0
	}

	return len(listeners)
}

func (self *Emitter) addListener(once bool, id EventID, callback func(...interface{})) *Emitter {
	if _, exists := self.listeners[id]; !exists {
		self.listeners[id] = []Listener{}
	}

	self.listeners[id] = append(self.listeners[id], Listener{callback, once})
	self.EmitSync(EVENT_REGISTER_LISTENER, id, callback)
	return self
}

// AddListener() - register a new listener on the specified event
func (self *Emitter) AddListener(id EventID, callback func(...interface{})) *Emitter {
	return self.addListener(false, id, callback)
}

// On() - register a new listener on the specified event
func (self *Emitter) On(id EventID, callback func(...interface{})) *Emitter {
	return self.addListener(false, id, callback)
}

// Once() - register a new one-time listener on the specified event
func (self *Emitter) Once(id EventID, callback func(...interface{})) *Emitter {
	return self.addListener(true, id, callback)
}

// RemoveListener() - remove the specified callback from the specified event's listeners
func (self *Emitter) RemoveListener(id EventID, callback func(...interface{})) *Emitter {
	arr, exists := self.listeners[id]
	if !exists {
		return self
	}

	for k, v := range arr {
		if &v.callback == &callback {
			self.listeners[id] = append(self.listeners[id][:k], self.listeners[id][k+1:]...)
			self.EmitSync(EVENT_UNREGISTER_LISTENER, id, callback)
			return self
		}
	}

	return self
}

// EmitSync() - run all listeners of the specified event in synchronous mode
func (self *Emitter) EmitSync(id EventID, args ...interface{}) *Emitter {
	arr := self.GetEventListenerList(id)
	if arr == nil {
		return self
	}

	for _, v := range arr {
		if v.once {
			self.RemoveListener(id, v.callback)
		}

		if args != nil && len(args) > 0 {
			v.callback(args...)
		} else {
			v.callback()
		}
	}

	return self
}

// EmitAsync() - run all listeners of the specified event in asynchronous mode using goroutines
func (self *Emitter) EmitAsync(id EventID, args ...interface{}) *Emitter {
	arr := self.GetEventListenerList(id)
	if arr == nil {
		return self
	}

	for _, v := range arr {
		if v.once {
			self.RemoveListener(id, v.callback)
		}

		if args != nil && len(args) > 0 {
			go v.callback(args...)
		} else {
			go v.callback()
		}
	}

	return self
}

// EmitAsync2() - run all listeners of the specified event in asynchronous mode using skeleton goroutine
func (self *Emitter) EmitAsync2(id EventID, args ...interface{}) *Emitter {
	evt := &Event{
		ID: id,
	}

	if args != nil && len(args) > 0 {
		evt.Args = args
	}

	self.queue <- evt

	return self
}
