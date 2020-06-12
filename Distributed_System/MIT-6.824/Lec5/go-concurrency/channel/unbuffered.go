// 如果有人试着通过channel发送东西
// 然而，当没有人去接收对面所发的东西时
// 那么这条线程就会阻塞，直到有人准备去接收数据为止
// 接着，在此时进行同步，将数据交接到接收者这一方

// 对另一方也是如此
// 如果有人试着从一个channel上接收数据
// 然而此时没有人发送数据
// 直到有另一个Goroutine准备往这个channel上发送数据，那么接收者这边的线程才不会被阻塞
// 此时，发送会同步进行（synchronously）

package main

import "time"
import "fmt"

func main() {
	c := make(chan bool)
	go func() {
		time.Sleep(1 * time.Second)
		<-c
	}()
	start := time.Now()
	c <- true // blocks until other goroutine receives
	fmt.Printf("send took %v\n", time.Since(start))
}
