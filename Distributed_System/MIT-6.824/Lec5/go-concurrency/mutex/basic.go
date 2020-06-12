package main

import "sync"
//import "time"

func main() {
	counter := 0
	var wg sync.WaitGroup
	var mu sync.Mutex
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			mu.Lock()
			defer mu.Unlock()
			defer wg.Done()
			counter = counter + 1
		}()
	}
	wg.Wait()
	//time.Sleep(1 * time.Second)
	mu.Lock()
	println(counter)
	mu.Unlock()
}
