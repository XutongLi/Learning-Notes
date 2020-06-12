// 你们会遇上这种情况，即raft peer变成一个candidate(候选人)
// 它想给它的所有follower发送投票请求
// 最后，follower会返回信息给candidate，并表示它有没有投票给这个candidate
// 我们想让candidate以并行的方式去询问所有的peer，这样，它就可以尽可能快的赢得选举
// 但此处有些复杂的地方，比如，当我们以并行的方式询问所有peer时，我们不想等到它们全员回复后才下定决心要决定哪个成为leader
// 因为如果一个candidate得票数过半，因为它无需去等待直到它获取到所有人的响应

// 初始版本

package main

import "time"
import "math/rand"

func main() {
	rand.Seed(time.Now().UnixNano())

	count := 0	// 得票为yes的数量
	finished := 0	// finished表示的是我总共得到的响应的数量

	// 此处的想法是我想以并行的方式发送投票请求
	// 并跟踪我拿到了多少支持票
	// 以及统计我总共接收了多少个响应
	// 然后一旦我知道我是否赢得了选举
	// 然后我就可以做出决定并继续下去

	for i := 0; i < 10; i++ {
		go func() {
			vote := requestVote()
			if vote {
				count++
			}
			finished++
		}()
	}

	for count < 5 && finished != 10 {
		// wait
	}
	if count >= 5 {
		println("received 5+ votes!")
	} else {
		println("lost")
	}
}

func requestVote() bool {
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	return rand.Int() % 2 == 0
}
