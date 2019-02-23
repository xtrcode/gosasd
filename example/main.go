package main

import (
	"fmt"
	"gosasd"
	"time"
)

func receiver(i interface{}) bool {
	fmt.Println(i)

	return true
}

func main() {
	async1 := make(chan interface{})
	async2 := make(chan interface{})
	receiver := make(chan interface{})

	synchronize := gosasd.NewSyncronizer(true)
	synchronize.AddAsyncChannel(async1)
	synchronize.AddAsyncChannel(async2)

	//synchronize.SetReceiverFunction(receiver)
	synchronize.SetReceiverChannel(receiver)
	defer synchronize.Sync()()

	go func() {
		for i := 0; i < 10; i++ {
			async1 <- "a"
			time.Sleep(time.Second)
		}
		close(async1)
	}()

	go func() {
		for i := 0; i < 5; i++ {
			async2 <- "b"
			time.Sleep(time.Second)
		}
		close(async2)
	}()

	// comment this loop for using a receiving function
	for {
		data, isOpen := <-receiver
		if !isOpen {
			break
		}

		fmt.Println(data)
	}

}
