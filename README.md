# Go Synchronize Asynchronous Data (gosasd)
gosasd is a small package which aims to synchronize asynchronous data streams (channels) by simply
defining one receiving function or channel which takes the output of all pre-defined asynchronous channels. 

This is part of a larger personal project where i used this functionality with
~280 [Pusher](https://pusher.com/) streams which needed to be synchronized. So maybe
someone else can take advantage of this. Feel free to **fork it!** or leave suggestions for
improvements. 

# Install
```bash
go get -t -v https://github.com/xtrcode/gosasd
```

# Usage
```go
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