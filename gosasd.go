// The contents of this file is free and unencumbered software released into the
// public domain. For more information, please refer to <http://unlicense.org/>

package gosasd

import (
	"log"
	"sync"
)

type GenericChan chan interface{}

type Synchronizer struct {
	wg             sync.WaitGroup
	signal         chan bool
	asyncChannels  map[interface{}][]Chan
	groupPipelines map[interface{}]PipelineChan
	closeCounter   int
	globalPipeline PipelineChan
	logging        bool
}

type Chan struct {
	C          GenericChan
	Identifier interface{}
}

func NewSyncronizer(signal chan bool, debug bool) *Synchronizer {
	return &Synchronizer{
		wg:             sync.WaitGroup{},
		logging:        debug,
		signal:         signal,
		asyncChannels:  make(map[interface{}][]Chan),
		groupPipelines: make(map[interface{}]PipelineChan),
		closeCounter:   0,
		globalPipeline: nil,
	}
}

func (s *Synchronizer) Log(str ...string) {
	if s.logging {
		log.Println(str)
	}
}

func (s *Synchronizer) SetGlobalPipeline(pc PipelineChan) {
	s.globalPipeline = pc
}

func (s *Synchronizer) SetGroupPipeline(group interface{}, pc PipelineChan) {
	s.groupPipelines[group] = pc
}

// AddAsyncChannel
// group (Optional) will be set to "global_internal" if nil
func (s *Synchronizer) AddAsyncChannel(group, identifier interface{}, rc GenericChan) {
	if group == nil {
		group = "global_internal"
		s.SetGroupPipeline(group, s.globalPipeline)
	}

	s.asyncChannels[group] = append(s.asyncChannels[group], Chan{rc, identifier})
}

func (s *Synchronizer) Sync() func() {
	// Loop over each group and its synced
	// pipeline
	for group, pipeline := range s.groupPipelines {
		for _, asyncChan := range s.asyncChannels[group] {
			go s.processPipeline(asyncChan, pipeline)
		}
	}

	return func() {
		s.wg.Wait()
		s.Log("Sending signal to finish operation")
		s.signal <- true
	}
}
