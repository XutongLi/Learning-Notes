//我们想去做些定时任务，直到某种事情发生
//例如，这里想启动一个raft，然后定期发送心跳（heartbeats ）信号
//但当我们在这个raft实例上调用.kill时，你们实际上会想去关闭所有的goroutine
//So，你不会让所有这些Goroutines依然运行在后台
//可以通过使用一个共享变量来这个控制线程，即决定这个Goroutine是否该寿终正寝
package main

import "time"
import "sync"

var done bool
var mu sync.Mutex

func main() {
	time.Sleep(1 * time.Second)
	println("started")
	go periodic()
	time.Sleep(5 * time.Second) // wait for a while so we can observe what ticker does
	mu.Lock()
	done = true
	mu.Unlock()
	println("cancelled")
	time.Sleep(3 * time.Second) // observe no output
}

func periodic() {
	for {
		println("tick")
		time.Sleep(1 * time.Second)
		mu.Lock()
		if done {
			return
		}
		mu.Unlock()
	}
}
