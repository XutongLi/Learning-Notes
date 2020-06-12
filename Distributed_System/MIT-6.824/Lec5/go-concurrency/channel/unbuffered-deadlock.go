package main

func main() {
	c := make(chan bool)
	<-c
	c <- true
	// <-c
}
