package main

import (
	"encoding/json"
	"fmt"

	"github.com/jasonblanchard/natsby"
)

func (s *Service) handleDebug(c *natsby.Context) {
	s.Logger.Info().Msg(fmt.Sprintf("%v", c.Msg))
	revision, err := bytesToEntryRevision(c.Msg.Data)
	if err != nil {
		c.Err = err
		return
	}
	s.Logger.Info().Msg(fmt.Sprintf("%v", revision))

	r := &EntryRevision{
		ID:         revision.GetEntry().GetId(),
		Text:       revision.GetEntry().GetText(),
		CreatorID:  revision.GetEntry().GetCreatorId(),
		CreatedAt:  timestamppbToTimePointer(revision.GetEntry().GetCreatedAt()),
		UpdatedAt:  timestamppbToTimePointer(revision.GetEntry().GetUpdatedAt()),
		DeleteTime: timestamppbToTimePointer(revision.GetEntry().GetDeleteTime()),
		ActorType:  revision.GetActor().Type.String(),
		ActorID:    revision.GetActor().GetId(),
	}

	serialized, _ := json.Marshal(r)
	s.Logger.Info().Msg(string(serialized))
}

func errorHandler(s *Service) natsby.RecoveryFunc {
	return func(c *natsby.Context, err interface{}) {
		s.Logger.Error().Msg(fmt.Sprintf("%v", err))

		if err != nil {
			s.Logger.Error().Msg(fmt.Sprintf("%v", err))
			return
		}

		if c.Msg.Reply != "" {
			c.NatsConnection.Publish(c.Msg.Reply, []byte(""))
		}
	}
}
