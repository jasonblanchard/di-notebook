package main

import (
	"fmt"

	"github.com/jasonblanchard/di-messages/packages/go/messages/notebook"
	"github.com/jasonblanchard/di-notebook/app"
	"github.com/jasonblanchard/di-notebook/mappers/protobufmapper"
	"github.com/jasonblanchard/natsby"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
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

func (s *Service) handleReadEntry(c *natsby.Context) {
	readEntryRequest := &notebook.ReadEntryRequest{}
	err := proto.Unmarshal(c.Msg.Data, readEntryRequest)

	if err != nil {
		c.Err = errors.Wrap(err, "Error unmarshalling data")
		return
	}

	readEntryInput, err := protobufmapper.ReadEntryRequestToReadEntryInput(readEntryRequest)

	traceID := readEntryRequest.GetContext().TraceId

	c.Set("traceID", traceID)

	if err != nil {
		c.Err = errors.Wrap(err, "Error mapping request")
		return
	}

	entry, err := s.ReadEntry(readEntryInput)

	if err != nil {
		c.Err = errors.Wrap(err, "Error ReadEntry")
		code := errorToCode(errors.Cause(err))
		payload, err := protobufmapper.ToNotebookErrorResponse("Permission denied", code, traceID)
		if err != nil {
			c.Err = errors.Wrap(err, "ReadEntry Error mapping error")
			return
		}
		c.ByteReplyPayload = payload
		return
	}

	response, err := protobufmapper.EntryToReadEntryResponse(entry, traceID)
	if err != nil {
		c.Err = errors.Wrap(err, "Error mapping response")
		return
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

	if err != nil {
		c.Err = errors.Wrap(err, "Error ChangeEntry")
		return
	}

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

	if err != nil {
		c.Err = errors.Wrap(err, "Error ListEntries")
		return
	}

	response, err := protobufmapper.ListEntriesOutputToListEntriesResponse(output)

	if err != nil {
		c.Err = errors.Wrap(err, "Error mapping response")
		return
	}

	c.ByteReplyPayload = response
}

func errorHandler(s *Service) natsby.RecoveryFunc {
	return func(c *natsby.Context, err interface{}) {
		s.Logger.Error().Msg(fmt.Sprintf("%v", err))

		code := errorToCode(err)
		traceID := fmt.Sprintf("%v", c.Get("traceID"))

		payload, err := protobufmapper.ToNotebookErrorResponse("something went wrong", code, traceID)

		if err != nil {
			s.Logger.Error().Msg(fmt.Sprintf("%v", err))
			return
		}

		if c.Msg.Reply != "" {
			c.NatsConnection.Publish(c.Msg.Reply, payload)
		}
	}
}
