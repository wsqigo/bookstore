package demo

import (
	"io"
)

type TaskPool struct {
	ch chan struct{}
}

func NewTaskPool(limit int) *TaskPool {
	t := &TaskPool{
		ch: make(chan struct{}, limit),
	}
	//提前准备好了令牌
	for i := 0; i < limit; i++ {
		t.ch <- struct{}{}
	}
	return t
}

func (t *TaskPool) Do(f func()) {
	token := <-t.ch
	// 异步执行
	go func() {
		defer func() {
			t.ch <- token
		}()
		f()
	}()
}

type TaskPoolWithCache struct {
	cache chan func()
	close chan struct{}
}

func NewTaskPoolWithCache(limit int, cacheSize int) *TaskPoolWithCache {
	t := &TaskPoolWithCache{
		cache: make(chan func(), cacheSize),
		close: make(chan struct{}, 1),
	}

	// 提前把goroutine开好
	for i := 0; i < limit; i++ {
		go func() {
			for {
				// 在 goroutine 里面不断尝试从 cache 里面拿到任务
				select {
				case task := <-t.cache:
					task()
				case <-t.close:
					return
				}
			}
		}()
	}
	return t
}

func (t *TaskPoolWithCache) Do(f func()) {
	t.cache <- f
}

func (t *TaskPoolWithCache) Close() error {
	close(t.cache)
	close(t.close)
	return nil
}

type ReadWriter interface {
	io.Reader
	io.Writer
}

type Shape interface {
	perimeter() float64
	area() float64
}
