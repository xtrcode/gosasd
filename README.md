# Go Synchronize Asynchronous Data (gosasd)
gosasd is a small package which aims to synchronize asynchronous data streams (channels) by simply
defining one channel which takes the output of all asynchronous channels. 

This is part of a larger personal project where i used this functionality with
~280 [Pusher](https://pusher.com/) streams which needed to be "synchronized". So maybe
someone else can take advantage of this. Feel free to **fork it!** or leave suggestions for
improvements. 

*The whole process of distributing incoming asynchronous streams into one straight kind-of synchronized stream is 
called Pipeline*

# Features
- Global streams (no group assigned)
- Grouped streams
- Receive payloads per group (group pipeline)
- Identify payload origin
- Receive payloads of all groups (global pipeline)
- Signal channel

# Install
```bash
go get -t -v https://github.com/xtrcode/gosasd
```

# Usage
```go
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
``` 
Output:
```bash
Global: chan4 30
Group 2: chan4 30
Global: chan3 20
2019/03/02 13:45:21 [Async channel chan3 got closed]
2019/03/02 13:45:21 [Leaving routine for  chan3]
2019/03/02 13:45:21 [Async channel chan1 got closed]
Group 1: chan3 20
Global: chan3 21
Group 1: chan3 21
2019/03/02 13:45:21 [Leaving routine for  chan1]
Global: chan2 10
2019/03/02 13:45:21 [Async channel chan2 got closed]
Global: chan4 31
2019/03/02 13:45:21 [Leaving routine for  chan2]
Global: chan3 22
Global: chan1 0
2019/03/02 13:45:21 [Async channel chan4 got closed]
Group 1: chan2 10
2019/03/02 13:45:21 [Leaving routine for  chan4]
Global: chan2 11
2019/03/02 13:45:21 [Sending signal to finish operation]
Group 1: chan2 11
Group 2: chan4 31
Global: chan4 32
Global: chan1 1
Global: chan2 12
Group 1: chan3 22
Global: chan3 23
Group 1: chan3 23
Group 2: chan4 32
Global: chan4 33
Group 1: chan2 12
Global: chan1 2
Global: chan3 24
Group 1: chan3 24
Group 2: chan4 33
Global: chan2 13
Global: chan1 3
Global: chan1 4
Global: chan4 34
Group 1: chan2 13
Global: chan2 14
Group 1: chan2 14
Group 2: chan4 34
Signal to close operation

Process finished with exit code 0
```  
# (UN)LICENSE
This is free and unencumbered software released into the public domain.

Anyone is free to copy, modify, publish, use, compile, sell, or
distribute this software, either in source code form or as a compiled
binary, for any purpose, commercial or non-commercial, and by any
means.

In jurisdictions that recognize copyright laws, the author or authors
of this software dedicate any and all copyright interest in the
software to the public domain. We make this dedication for the benefit
of the public at large and to the detriment of our heirs and
successors. We intend this dedication to be an overt act of
relinquishment in perpetuity of all present and future rights to this
software under copyright law.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
IN NO EVENT SHALL THE AUTHORS BE LIABLE FOR ANY CLAIM, DAMAGES OR
OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
OTHER DEALINGS IN THE SOFTWARE.

For more information, please refer to <http://unlicense.org/>