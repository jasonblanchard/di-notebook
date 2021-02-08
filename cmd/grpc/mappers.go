package main

import (
	"fmt"

	notebook "github.com/jasonblanchard/di-apis/gen/pb-go/notebook/v2"
	"github.com/jasonblanchard/di-notebook/app"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
)

// EntryToEntryRevision mapper
func EntryToEntryRevision(e *app.Entry, p *notebook.Principal) ([]byte, error) {
	entryRevision := &notebook.EntryRevision{
		Entry: &notebook.Entry{
			Id:         fmt.Sprintf("%d", e.ID),
			Text:       e.Text,
			CreatedAt:  timeToProtoTime(e.CreatedAt),
			UpdatedAt:  timeToProtoTime(e.UpdatedAt),
			DeleteTime: timeToProtoTime(*e.DeleteTime),
		},
		Actor: p,
	}

	output, err := proto.Marshal(entryRevision)
	if err != nil {
		errors.Wrap(err, "Wrror marshalling entryRevision")
	}

	return output, nil
}
