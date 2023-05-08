package demo

import (
	"fmt"
	"sync"
)

// PublicResource 你永远不知道你的用户拿了它会干啥
// 他即使不用 PublicResourceLock 你也毫无办法
// 如果你用这个resource，一定要用锁
var PublicResource map[string]string
var PublicLock sync.Mutex

// privateResource 要好一点，祈祷你的同事会来看你的注释，知道要用锁
// 很多库都是这么写的，我也写了很多类似的代码
var privateResource map[string]string
var privateLock sync.Mutex

// safeResource 很棒，所有的期望对资源的操作都只能通过定义在上 safeResource 上的方法来进行
type safeResource struct {
	resource map[string]string
	lock     sync.Mutex
}

func (s *safeResource) DoSomethingToResource() {
	s.lock.Lock()
	defer s.lock.Unlock()
}

type SafeMap[K comparable, V any] struct {
	values map[K]V
	lock   sync.RWMutex
}

// 已经有 key，返回对应的值，然后loaded = true
// 没有，则放进去，返回loaded = false
// goroutine 1 => ("key1", 1)
// goroutine 2 => ("key1", 2)
func (s *SafeMap[K, V]) LoadOrStore(key K, newValue V) (V, bool) {
	s.lock.RLock()
	oldVal, ok := s.values[key]
	s.lock.RUnlock()
	if ok {
		return oldVal, true
	}

	s.lock.Lock()
	defer s.lock.Unlock()
	// double check 加了写锁再检查一遍
	oldVal, ok = s.values[key]
	if ok {
		return oldVal, true
	}
	s.values[key] = newValue
	return newValue, false
}

type MyBiz struct {
	once sync.Once
}

// 只被执行一次
func (m *MyBiz) Init() {
	m.once.Do(func() {

	})
}

type singleton struct {
}

func (s *singleton) Single() {
	fmt.Println("I am single")
}

var instance *singleton
var instanceOnce sync.Once

// GetSingleInstance 返回接口
func GetSingleInstance() *singleton {
	instanceOnce.Do(func() {
		instance = &singleton{}
	})
	return instance
}
