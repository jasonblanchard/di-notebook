package main

import (
	notebook "github.com/jasonblanchard/di-apis/gen/pb-go/notebook/v2"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
)

func bytesToEntryRevision(data []byte) (*notebook.EntryRevision, error) {
	revision := &notebook.EntryRevision{}
	err := proto.Unmarshal(data, revision)
	if err != nil {
		return nil, errors.Wrap(err, "Error unmarshalling EntryRevision")
	}

	return revision, nil
}
