// The contents of this file is free and unencumbered software released into the
// public domain. For more information, please refer to <http://unlicense.org/>

// The whole process of distributing incoming asynchronous streams
// into one straight kind-of synchronized stream is called Pipeline
// This files defines helpers and structures.

package gosasd

import (
	"fmt"
	"reflect"
)

type PipelineChan chan PipelinePayload

// PipelinePayload contains the actual data that is outgoing.
// It ships the original payload (Data) and an identifier
// to distinguish over which channel we received the payload
type PipelinePayload struct {
	// Original payload received from channel
	Data interface{}
	// (Optional) Identifier for origin channel
	Identifier interface{}
}

// Is checks if identifiers are deeply equal
// Refer to: https://golang.org/pkg/reflect/#DeepEqual
func (pp *PipelinePayload) Is(_id interface{}) bool {
	return reflect.DeepEqual(pp.Identifier, _id)
}

func (s *Synchronizer) processPipeline(c Chan, group PipelineChan) {
	s.wg.Add(1)
	defer s.wg.Done()

	for {
		data, isOpen := <-c.C
		if !isOpen {
			s.Log("Async channel", fmt.Sprintf("%v", c.Identifier), "got closed")
			break
		}

		payload := PipelinePayload{
			Data:       data,
			Identifier: c.Identifier,
		}

		if s.globalPipeline != nil && s.globalPipeline != group {
			s.globalPipeline <- payload
		}

		if group != nil {
			group <- payload
		}
	}

	s.Log("Leaving routine for ", fmt.Sprintf("%v", c.Identifier))
}
