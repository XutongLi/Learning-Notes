// 有正在并发运行的多个线程，它们要对某个共享变量进行更新
// 然后我有另一条线程，它正在等待该共享数据中的某个属性或某个条件变为true
// 该线程会等待直到这个条件变为true
// 这个工具叫做condition variable（条件变量）解决这个问题

package main

import "sync"
import "time"
import "math/rand"

func main() {
	rand.Seed(time.Now().UnixNano())

	count := 0
	finished := 0
	var mu sync.Mutex
	cond := sync.NewCond(&mu)
	// 当共享数据中的某个条件或某些属性变为true的时候，我们会通过condition variable（条件变量）来进行协调

	for i := 0; i < 10; i++ {
		go func() {
			vote := requestVote()
			mu.Lock()
			defer mu.Unlock()
			if vote {
				count++
			}
			finished++
			cond.Broadcast()
			// 当我们做某些要修改数据的事情时，我们调用cond.Broadcast()
			// 要在持有锁的时候才做这件事
			// 它所做的就是去唤醒正在等待这个condition variable（条件变量）的线程
		}()
	}

	// 当你检查条件时
	// 你拿到了锁
	// 接着，你总是先检查循环的条件
	// 接着，在循环内部，当条件为false时，你要去调用cond.Wait()
	// 只有当你持有锁时，才能去调用这个，然后以原子的方式将锁释放掉
	// 并将它自己放入正等待线程的列表中
	// 当我们调用完这个cond.Wait()后，然后我们会回到这个for循环的顶部
	mu.Lock()
	for count < 5 && finished != 10 {
		cond.Wait()
		// 另一边正在等待该共享数据上的某些条件变为true
		// 调用cond.Wait()
		// 为了让其他人能继续干活，它就会将手里的锁释放
		// 然后，它将自己这条线程添加到等待该condition variable（条件变量）的名单上（condition等待队列）
	}
	if count >= 5 {
		println("received 5+ votes!")
	} else {
		println("lost")
	}
	mu.Unlock()
}

func requestVote() bool {
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	return rand.Int() % 2 == 0
}
