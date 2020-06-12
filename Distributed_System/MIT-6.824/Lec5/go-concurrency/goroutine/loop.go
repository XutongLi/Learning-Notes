// 如果你需要候选人们投票（以选出谁是primary或leader），我们想同时从所有followers（从机）那里投票，而不是一个接一个地逐个进行（当primary有问题后，这些followers们需要站出来投个票选primary）
// 或者，类似的，leader可能想给所有的follower发送追加内容条目
// 此处的i所代表的可能是我们所试着发送的follower的索引号

package main

import "sync"

func main() {
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(x int) {
			sendRPC(x)
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func sendRPC(i int) {
	println(i)
}
