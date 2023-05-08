package singleton

import (
	"fmt"
	"sync"
)

// Singleton 是单例模式接口，导出的
// 通过该接口可以避免 GetInstance 返回一个包私有类型的指针
type Singleton interface {
	Single()
}

// singleton 是单例模式类，包私有的
type singleton struct{}

func (s *singleton) Single() {
	fmt.Println("I am single")
}

var instance *singleton
var instanceOnce sync.Once

func GetSingleInstance() Singleton {
	instanceOnce.Do(func() {
		instance = &singleton{}
	})
	return instance
}
