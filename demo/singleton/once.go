package singleton

import (
	"sync"
	"sync/atomic"
)

type Once struct {
	done uint32

	mu sync.Mutex
}

func (o *Once) Do(f func()) {
	if atomic.LoadUint32(&o.done) == 0 {
		o.doSlow(f)
	}
}

func (o *Once) doSlow(f func()) {
	o.mu.Lock()
	defer o.mu.Unlock()

	if o.done == 0 {
		defer atomic.AddUint32(&o.done, 1)
		f()
	}
}
