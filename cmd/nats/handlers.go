package main

import (
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
