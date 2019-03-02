// The contents of this file is free and unencumbered software released into the
// public domain. For more information, please refer to <http://unlicense.org/>

package main

import (
	"fmt"
	"gosasd"
)

func main() {
	async1 := make(chan interface{})
	async2 := make(chan interface{})
	async3 := make(chan interface{})
	async4 := make(chan interface{})

	signal := make(chan bool)
	globalPipeline := make(chan gosasd.PipelinePayload)
	groupOnePipeline := make(chan gosasd.PipelinePayload)
	groupTwoPipeline := make(chan gosasd.PipelinePayload)

	sync := gosasd.NewSyncronizer(signal, true)

	// no group, identifier: chan1
	sync.AddAsyncChannel(nil, "chan1", async1)

	// 2 element group
	// identifier: chan2
	sync.AddAsyncChannel("group1", "chan2", async2)
	// identifier: chan3
	sync.AddAsyncChannel("group1", "chan3", async3)

	// 1 element group, identifier: chan4
	sync.AddAsyncChannel("group2", "chan4", async4)

	// receive all payloads
	sync.SetGlobalPipeline(globalPipeline)
	sync.SetGroupPipeline("group1", groupOnePipeline)
	sync.SetGroupPipeline("group2", groupTwoPipeline)

	go sync.Sync()()

	{
		go func() {
			for i := 0; i < 5; i++ {
				async1 <- i
			}
			close(async1)
		}()

		go func() {
			for i := 10; i < 15; i++ {
				async2 <- i
			}
			close(async2)
		}()

		go func() {
			for i := 20; i < 25; i++ {
				async3 <- i
			}
			close(async3)
		}()

		go func() {
			for i := 30; i < 35; i++ {
				async4 <- i
			}
			close(async4)
		}()

	}

LOOP:
	for {
		select {
		case payload := <-globalPipeline:
			fmt.Println("Global:", payload.Identifier, payload.Data)

		case payload := <-groupOnePipeline:
			fmt.Println("Group 1:", payload.Identifier, payload.Data)

		case payload := <-groupTwoPipeline:
			fmt.Println("Group 2:", payload.Identifier, payload.Data)

		case _ = <-signal:
			fmt.Println("Signal to close operation")
			break LOOP
		}
	}
}
