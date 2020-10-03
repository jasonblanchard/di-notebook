package protobufmapper

import (
	"github.com/jasonblanchard/di-notebook/app"
	"github.com/jasonblanchard/di-notebook/di_messages/entry"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
)

// CreateEntryRequestToStartNewEntryInput mapping
func CreateEntryRequestToStartNewEntryInput(data []byte) (*app.StartNewEntryInput, error) {
	createEntryRequest := &entry.CreateEntryRequest{}
	err := proto.Unmarshal(data, createEntryRequest)
	if err != nil {
		return nil, errors.Wrap(err, "Error decoding CreateEntryRequest message")
	}

	startNewEntryInput := &app.StartNewEntryInput{
		Principle: &app.Principle{
			Type: app.PrincipleTypeUser,
			ID:   createEntryRequest.Context.Principal.Id,
		},
		Text:      createEntryRequest.Payload.Text,
		CreatorID: createEntryRequest.Payload.CreatorId,
	}

	return startNewEntryInput, nil
}

// IDToCreateEntryResponse mapping
func IDToCreateEntryResponse(id int) ([]byte, error) {
	createEntryResponse := &entry.CreateEntryResponse{
		Payload: &entry.CreateEntryResponse_Payload{
			Id: string(id),
		},
	}

	output, err := proto.Marshal(createEntryResponse)
	if err != nil {
		errors.Wrap(err, "Failed to marshall message")
	}

	return output, nil
}
