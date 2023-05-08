package demo

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

type MyPool struct {
	p      sync.Pool
	maxCnt int32
	cnt    int32
}

func (m *MyPool) Get() any {
	return m.p.Get()
}

func (m *MyPool) Put(val any) {
	// 大对象不放回去
	if unsafe.Sizeof(val) > 1024 {
		return
	}

	// 超过数量限制
	if atomic.LoadInt32(&m.cnt) >= atomic.LoadInt32(&m.maxCnt) {
		return
	}
	m.p.Put(val)
}
