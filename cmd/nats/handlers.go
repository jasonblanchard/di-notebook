package main

import (
	"github.com/jasonblanchard/di-notebook/mappers/protobufmapper"
	"github.com/jasonblanchard/natsby"
	"github.com/pkg/errors"
)

func (s *Service) handleCreateEntry(c *natsby.Context) {
	startNewEntryInput, err := protobufmapper.CreateEntryRequestToStartNewEntryInput(c.Msg.Data)
	if err != nil {
		c.Err = errors.Wrap(err, "Error mapping request")
		return
	}

	id, err := s.StartNewEntry(startNewEntryInput)
	if err != nil {
		c.Err = errors.Wrap(err, "Error calling StartNewEntry")
		return
	}

	payload, err := protobufmapper.IDToCreateEntryResponse(id)
	if err != nil {
		c.Err = errors.Wrap(err, "Error mapping response")
		return
	}

	c.ByteReplyPayload = payload
}
