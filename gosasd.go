// The contents of this file is free and unencumbered software released into the
// public domain. For more information, please refer to <http://unlicense.org/>

package gosasd

import (
	"github.com/eapache/channels"
	"log"
	"sync"
)

type Synchronizer struct {
	wg             sync.WaitGroup
	signal         chan bool
	asyncChannels  map[interface{}][]ChannelContainer
	groupPipelines map[interface{}]PipelineChan
	closeCounter   int
	globalPipeline PipelineChan
	logging        bool
}

type ChannelContainer struct {
	C          channels.SimpleOutChannel
	Identifier interface{}
}

func NewSyncronizer(signal chan bool, debug bool) *Synchronizer {
	return &Synchronizer{
		wg:             sync.WaitGroup{},
		logging:        debug,
		signal:         signal,
		asyncChannels:  make(map[interface{}][]ChannelContainer),
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
func (s *Synchronizer) AddAsyncChannel(group, identifier interface{}, c channels.SimpleOutChannel) {
	if group == nil {
		group = "global_internal"
		s.SetGroupPipeline(group, s.globalPipeline)
	}

	s.asyncChannels[group] = append(s.asyncChannels[group], ChannelContainer{c, identifier})
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
