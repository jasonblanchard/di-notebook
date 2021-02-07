package main

import (
	"github.com/gogo/protobuf/proto"
	notebook "github.com/jasonblanchard/di-apis/gen/pb-go/notebook/v2"
	"github.com/pkg/errors"
)

func bytesToEntryRevision(data []byte) (*notebook.EntryRevision, error) {
	revision := &notebook.EntryRevision{}
	err := proto.Unmarshal(data, revision)
	if err != nil {
		return nil, errors.Wrap(err, "Error unmarshalling EntryREvision")
	}

	return revision, nil
}
