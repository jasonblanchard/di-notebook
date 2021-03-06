package main

import (
	"fmt"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	notebook "github.com/jasonblanchard/di-apis/gen/pb-go/notebook/v2"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestEndToEnd(t *testing.T) {
	nc, err := nats.Connect(nats.DefaultURL)

	revision := &notebook.EntryRevision{
		Entry: &notebook.Entry{
			Text:      "a new one from the test",
			CreatedAt: ptypes.TimestampNow(),
			UpdatedAt: &timestamppb.Timestamp{},
			DeleteTime: &timestamppb.Timestamp{
				Seconds: 53456356,
			},
		},
		Actor: &notebook.Principal{
			Type: notebook.Principal_USER,
		},
	}
	data, err := proto.Marshal(revision)
	if err != nil {
		panic(err)
	}

	nc.Publish("data.mesh.notebook.v2.EntryRevision", data)

	r := &notebook.EntryRevision{}
	_ = proto.Unmarshal(data, r)
	fmt.Println(r.GetEntry().GetText())

	assert.Equal(t, true, true)
}
