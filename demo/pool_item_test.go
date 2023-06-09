package demo

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
)

func TestPool(t *testing.T) {
	pool := sync.Pool{
		New: func() any {
			// 创建函数，sync.Pool会回调
			return &User{}
		},
	}
	u1 := pool.Get().(*User)
	u1.ID = 12
	u1.Name = "Tom"
	// 一通操作
	// 放回去之前要先重置掉
	u1.Reset()
	pool.Put(u1)

	u2 := pool.Get().(*User)
	fmt.Println(u2)
}

type User struct {
	ID   uint64
	Name string
}

func (u *User) Reset() {
	u.ID = 0
	u.Name = ""
}

func TestWaitGroup(t *testing.T) {
	wg := sync.WaitGroup{}
	var result int64
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(delta int) {
			defer wg.Done()
			atomic.AddInt64(&result, int64(delta))
		}(i)
	}
	wg.Wait()
}
