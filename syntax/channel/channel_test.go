package channel

import (
	"math/rand"
	"testing"
	"time"
)

func TestChannel(t *testing.T) {
	// 一个可以缓存容量为 1 的，放 int 类型数据的 channel
	ch := make(chan int, 1)
	// 写入数据
	ch <- 123
	// 打印读取出来的数据
	println(<-ch)
	close(ch)
}

func TestChannelClose(t *testing.T) {
	// 一个可以缓存容量为 1 的，放 int 类型数据的 channel
	ch := make(chan int, 1)
	// 写入数据
	ch <- 0
	val, ok := <-ch
	t.Log(val, ok)

	close(ch)
	val, ok = <-ch
	t.Log(val, ok)
}

func TestForLoop(t *testing.T) {
	// 这个是没有缓存的
	ch := make(chan int)
	go func() {
		for i := 0; i < 3; i++ {
			ch <- i
			time.Sleep(time.Second)
		}
		close(ch)
	}()
	for val := range ch {
		t.Log(val)
	}
	t.Log("发送完毕")
}

func TestSelect(t *testing.T) {
	ch1 := make(chan int, 1)
	ch2 := make(chan int, 1)
	go func() {
		time.Sleep(time.Second)
		ran := rand.Int()
		if ran%2 == 1 {
			ch1 <- ran
		} else {
			ch2 <- ran
		}
	}()
	select {
	case val := <-ch1:
		t.Log(val)
	case val := <-ch2:
		t.Log(val)
	}
}
