package gosasd

import (
	"log"
	"strconv"
	"sync"
)

type Synchronizer struct {
	wg            sync.WaitGroup
	asyncChannels []chan interface{}
	closeCounter  int
	outputChannel chan interface{}
	receiverFunc  Receiver
	receiverChan  chan interface{}
	logging       bool
}

type Receiver func(i interface{}) bool

func NewSyncronizer(debug bool) *Synchronizer {
	return &Synchronizer{
		wg:            sync.WaitGroup{},
		logging:       debug,
		asyncChannels: nil,
		closeCounter:  0,
		outputChannel: make(chan interface{}),
		receiverFunc:  nil,
		receiverChan:  nil,
	}
}

func (s *Synchronizer) Log(str ...string) {
	if s.logging {
		log.Println(str)
	}
}

func (s *Synchronizer) SetReceiverFunction(r Receiver) {
	s.receiverFunc = r
}

func (s *Synchronizer) SetReceiverChannel(rc chan interface{}) {
	s.receiverChan = rc
}

func (s *Synchronizer) AddAsyncChannel(c chan interface{}) {
	s.asyncChannels = append(s.asyncChannels, c)
}

func (s *Synchronizer) Sync() func() {
	for _, asyncChan := range s.asyncChannels {
		s.wg.Add(1)
		go func(c chan interface{}) {
			defer s.wg.Done()

			for {
				data, isOpen := <-c
				if !isOpen {
					s.closeCounter++

					s.Log("Async channel number", strconv.Itoa(s.closeCounter), "got closed")

					if s.closeCounter == len(s.asyncChannels) {
						s.Log("All async channels are closed! Closing receiver channel!")

						if s.receiverChan != nil {
							close(s.receiverChan)
						}
					}

					break
				}

				if s.receiverFunc != nil {
					s.receiverFunc(data)
				} else if s.receiverChan != nil {
					s.receiverChan <- data
				}
			}
		}(asyncChan)
	}

	return func() {
		s.wg.Wait()
	}
}
