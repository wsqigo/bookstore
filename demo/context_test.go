package demo

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestParentValueCtx(t *testing.T) {
	ctx := context.Background()
	context.WithCancel(ctx)
	childCtx := context.WithValue(ctx, "key1", 123)
	ccCtx := context.WithValue(childCtx, "key2", 124)
	val := childCtx.Value("key2")
	fmt.Println(val)
	val = ccCtx.Value("key2")
	fmt.Println(val)
}

func TestTimeoutExample(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	bsChan := make(chan struct{})
	go func() {
		slowBusiness()
		bsChan <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		fmt.Println("timeout")
	case <-bsChan:
		fmt.Println("business end")
	}

}

func slowBusiness() {
	time.Sleep(2 * time.Second)
}

func TestTimeoutTimeAfter(t *testing.T) {
	bsChan := make(chan struct{})
	go func() {
		slowBusiness()
		bsChan <- struct{}{}
	}()

	timer := time.AfterFunc(time.Second, func() {
		fmt.Println("timeout")
	})
	<-bsChan
	fmt.Println("business end")
	timer.Stop()
}

func Test_test(t *testing.T) {
	fmt.Println("wsqigo")
}

// A canceler is a context type that can be canceled directly. The
// implementations are *cancelCtx and *timerCtx.
type canceler interface {
	cancel(removeFromParent bool, err, cause error)
	Done() <-chan struct{}
}

type cancelCtx struct {
	context.Context

	done     atomic.Value
	mu       sync.Mutex
	children map[canceler]struct{}
}

func (c *cancelCtx) Done() <-chan struct{} {
	d := c.done.Load()
	if d != nil {
		return d.(chan struct{})
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	d = c.done.Load()
	if d == nil {
		d = make(chan struct{})
		c.done.Store(d)
	}
	return d.(chan struct{})
}

func Step1(ctx context.Context) {
	var db *sql.DB
	db.ExecContext(ctx, "UPDATE XXXX", 1)
}

type Once struct {
	done uint32
	m    sync.Mutex
}

func (o *Once) DO(f func()) {
	if atomic.LoadUint32(&o.done) == 0 {
		o.doSlow(f)
	}
}

func (o *Once) doSlow(f func()) {
	o.m.Lock()
	defer o.m.Unlock()
	if o.done == 0 {
		defer atomic.StoreUint32(&o.done, 1)
		f()
	}
}
