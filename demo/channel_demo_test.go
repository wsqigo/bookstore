package demo

import (
	"fmt"
	"testing"
	"time"
)

type Broker struct {
	ch        chan string
	consumers []func(s string)
}

func (b *Broker) Subscribe(consume func(s string)) {
	b.consumers = append(b.consumers, consume)
}

func (b *Broker) Produce(msg string) {
	b.ch <- msg
}

func (b *Broker) Start() {
	go func() {
		for {
			s := <-b.ch
			for _, c := range b.consumers {
				c(s)
			}
		}
	}()
}

type Consumer struct {
	ch chan string
}

func TestNilChannel(t *testing.T) {
	ch1 := make(chan int)
	ch2 := make(chan int)

	go func() {
		time.Sleep(5 * time.Second)
		ch1 <- 5
		close(ch1)
	}()

	go func() {
		time.Sleep(7 * time.Second)
		ch2 <- 7
		close(ch2)
	}()

	var ok1, ok2 bool
	for {
		select {
		case x := <-ch1:
			ok1 = true
			fmt.Println(x)
		case x := <-ch2:
			ok2 = true
			fmt.Println(x)
		}
		if ok1 && ok2 {
			break
		}
	}
	fmt.Println("program end")
}
