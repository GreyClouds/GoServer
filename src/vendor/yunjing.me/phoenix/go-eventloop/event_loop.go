package eventloop

import (
	"fmt"
	"math"
	"runtime"
	"sync"
)

type EventLoop struct {
	workers uint
	work    chan *workUnit
	cancel  chan struct{}
	closed  bool
	m       sync.RWMutex
}

// NewLimited returns a new limited pool instance
func New(workers uint) *EventLoop {
	if workers == 0 {
		panic("invalid workers '0'")
	}

	p := &EventLoop{
		workers: workers,
	}

	p.initialize()

	return p
}

func (p *EventLoop) initialize() {
	p.work = make(chan *workUnit, p.workers)
	p.cancel = make(chan struct{})
	p.closed = false

	// fire up worker here
	go p.newWorker()
}

// passing work and cancel channels to newWorker() to avoid any potential race condition
// betweeen p.work read & write
func (p *EventLoop) newWorker() {
	var wu *workUnit

	defer func() {
		if err := recover(); err != nil {
			trace := make([]byte, 1<<16)
			n := runtime.Stack(trace, true)

			s := fmt.Sprintf(errRecovery, err, string(trace[:int(math.Min(float64(n), float64(7000)))]))

			iwu := wu
			iwu.err = &ErrRecovery{s: s}
			close(iwu.done)

			println(s)
		}
	}()

	var value interface{}
	var err error

	for {
		select {
		case wu = <-p.work:

			// possible for one more nilled out value to make it
			// through when channel closed, don't quite understad the why
			if wu == nil {
				continue
			}

			// support for individual WorkUnit cancellation
			// and batch job cancellation
			if wu.cancelled.Load() == nil {
				value, err = wu.fn(wu)

				wu.writing.Store(struct{}{})

				// need to check again in case the WorkFunc cancelled this unit of work
				// otherwise we'll have a race condition
				if wu.cancelled.Load() == nil && wu.cancelling.Load() == nil {
					wu.value, wu.err = value, err

					// who knows where the Done channel is being listened to on the other end
					// don't want this to block just because caller is waiting on another unit
					// of work to be done first so we use close
					close(wu.done)
				}
			}

		case <-p.cancel:
			return
		}
	}
}

// Queue queues the work to be run, and starts processing immediately
func (p *EventLoop) Execute(fn WorkFunc) WorkUnit {

	w := &workUnit{
		done: make(chan struct{}),
		fn:   fn,
	}

	go func() {
		p.m.RLock()
		if p.closed {
			w.err = &ErrPoolClosed{s: errClosed}
			if w.cancelled.Load() == nil {
				close(w.done)
			}
			p.m.RUnlock()
			return
		}

		p.work <- w

		p.m.RUnlock()
	}()

	return w
}

// Reset reinitializes a pool that has been closed/cancelled back to a working state.
// if the pool has not been closed/cancelled, nothing happens as the pool is still in
// a valid running state
func (p *EventLoop) Reset() {

	p.m.Lock()

	if !p.closed {
		p.m.Unlock()
		return
	}

	// cancelled the pool, not closed it, pool will be usable after calling initialize().
	p.initialize()
	p.m.Unlock()
}

func (p *EventLoop) closeWithError(err error) {

	p.m.Lock()

	if !p.closed {
		close(p.cancel)
		close(p.work)
		p.closed = true
	}

	for wu := range p.work {
		wu.cancelWithError(err)
	}

	p.m.Unlock()
}

// Close cleans up the pool workers and channels and cancels any pending
// work still yet to be processed.
// call Reset() to reinitialize the pool for use.
func (p *EventLoop) Close() {

	err := &ErrPoolClosed{s: errClosed}
	p.closeWithError(err)
}