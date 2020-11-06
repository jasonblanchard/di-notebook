package main

import (
	"github.com/jasonblanchard/di-notebook/app"
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

func (s *Service) handleGetEntry(c *natsby.Context) {
	readEntryInput, err := protobufmapper.GetEntryRequestToReadEntryInput(c.Msg.Data)
	if err != nil {
		c.Err = errors.Wrap(err, "Error mapping request")
		return
	}

	entry, err := s.ReadEntry(readEntryInput)

	if err != nil {
		c.Err = errors.Wrap(err, "Error ReadEntry")
		return
	}

	response, err := protobufmapper.EntryToGetEntryResponse(entry)
	if err != nil {
		c.Err = errors.Wrap(err, "Error mapping response")
	}

	c.ByteReplyPayload = response
}

func (s *Service) handleUpdateEntry(c *natsby.Context) {
	changeEntryInput, err := protobufmapper.UpdateEntryRequestToChangeEntryInput(c.Msg.Data)
	if err != nil {
		c.Err = errors.Wrap(err, "Error mapping request")
		return
	}

	entry, err := s.ChangeEntry(changeEntryInput, func(entry *app.Entry) {
		infoEntryUpdatedPayload, err := protobufmapper.ChangeEntryOutputToInfoEntryUpdated(entry)
		if err != nil {
			c.Err = errors.Wrap(err, "Error mapping info")
			return
		}
		c.NatsConnection.Publish("provisional.info.entry.updated", infoEntryUpdatedPayload)
	})

	response, err := protobufmapper.ChangeEntryOutputToUpdateEntryResponse(entry)
	if err != nil {
		c.Err = errors.Wrap(err, "Error mapping response")
		return
	}

	c.ByteReplyPayload = response
}

func (s *Service) handleDeleteEntry(c *natsby.Context) {
	deleteEntryInput, err := protobufmapper.DeleteEntryRequestToDiscardEntryInput(c.Msg.Data)
	if err != nil {
		c.Err = errors.Wrap(err, "Error mapping request")
		return
	}

	err = s.DiscardEntry(deleteEntryInput)
	if err != nil {
		c.Err = errors.Wrap(err, "Error discarding entry")
		return
	}

	response, err := protobufmapper.DiscardEntryToDeleteEntryResponse()
	if err != nil {
		c.Err = errors.Wrap(err, "Error mapping response")
		return
	}

	c.ByteReplyPayload = response
}

func (s *Service) handleListEntries(c *natsby.Context) {
	listEntriesInput, err := protobufmapper.ListEntriesRequestToListEntriesInput(c.Msg.Data)
	if err != nil {
		c.Err = errors.Wrap(err, "Error mapping request")
		return
	}

	output, err := s.ListEntries(listEntriesInput)

	response, err := protobufmapper.ListEntriesOutputToListEntriesResponse(output)

	if err != nil {
		c.Err = errors.Wrap(err, "Error mapping response")
		return
	}

	c.ByteReplyPayload = response
}
