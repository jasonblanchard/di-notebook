package main

import (
	"fmt"

	"github.com/jasonblanchard/di-notebook/app"
	"github.com/jasonblanchard/di-notebook/di_messages/entry"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
)

// ChangeEntryOutputToInfoEntryUpdated mapper
func ChangeEntryOutputToInfoEntryUpdated(e *app.Entry) ([]byte, error) {
	infoEntryUpdated := &entry.InfoEntryUpdated{
		Payload: &entry.InfoEntryUpdated_Payload{
			Id:        fmt.Sprintf("%d", e.ID),
			Text:      e.Text,
			CreatedAt: timeToProtoTime(e.CreatedAt),
			UpdatedAt: timeToProtoTime(e.UpdatedAt),
		},
	}

	output, err := proto.Marshal(infoEntryUpdated)
	if err != nil {
		errors.Wrap(err, "Wrror marshalling infoEntryUpdated")
	}

	return output, nil
}
