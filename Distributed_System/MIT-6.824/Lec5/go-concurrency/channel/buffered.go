// 实际上几乎没有问题需要使用buffer channel去解决
// buffered channel内部有一些存储空间
// 直到该空间被填满前，直到该空间被填满前
// 但一旦channel中的空间满了，那么它所表现得就和non-buffered channel一样了
// 也就是说，此后的发送会一直处于阻塞状态，直到对面开始接收数据，这样channel中的空间就会空出来了
// 从一个高级层面来讲，应该避免使用buffered channel
package main

import "time"
import "fmt"

func main() {
	c := make(chan bool, 1)
	go func() {
		time.Sleep(1 * time.Second)
		<-c
	}()
	start := time.Now()
	c <- true
	fmt.Printf("send took %v\n", time.Since(start))

	start = time.Now()
	c <- true
	fmt.Printf("send took %v\n", time.Since(start))
}
