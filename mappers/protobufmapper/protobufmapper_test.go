package protobufmapper

import (
	"testing"

	"github.com/jasonblanchard/di-notebook/di_messages/entry"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

func TestCreateEntryRequestToStartNewEntryInput(t *testing.T) {
	createEntryRequest := &entry.CreateEntryRequest{
		Context: &entry.RequestContext{
			Principal: &entry.Principal{
				Type: entry.Principal_USER,
				Id:   "123",
			},
		},
		Payload: &entry.CreateEntryRequest_Payload{
			Text:      "Hello, world",
			CreatorId: "123",
		},
	}

	b, err := proto.Marshal(createEntryRequest)
	if err != nil {
		panic(err)
	}

	output, err := CreateEntryRequestToStartNewEntryInput(b)
	assert.Equal(t, output.Principle.ID, "123")
	assert.Equal(t, output.Text, "Hello, world")
	assert.Equal(t, output.CreatorID, "123")
}

func TestBadCreateEntryRequestToStartNewEntryInput(t *testing.T) {
	createEntryRequest := &entry.Error{}

	b, err := proto.Marshal(createEntryRequest)
	if err != nil {
		panic(err)
	}

	_, err = CreateEntryRequestToStartNewEntryInput(b)
	assert.NotNil(t, err)
}
