package protobufmapper

import (
	"fmt"
	"strconv"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
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

	// TODO: Come up with better validation logic
	if createEntryRequest.Context == nil {
		return nil, errors.New("Validation error")
	}

	startNewEntryInput := &app.StartNewEntryInput{
		Principle: &app.Principal{
			Type: app.PrincipalUSER,
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
			Id: fmt.Sprintf("%d", id),
		},
	}

	output, err := proto.Marshal(createEntryResponse)
	if err != nil {
		errors.Wrap(err, "Failed to marshall message")
	}

	return output, nil
}

// GetEntryRequestToReadEntryInput mapper
func GetEntryRequestToReadEntryInput(data []byte) (*app.ReadEntryInput, error) {
	getEntryRequest := &entry.GetEntryRequest{}
	err := proto.Unmarshal(data, getEntryRequest)
	if err != nil {
		return nil, errors.Wrap(err, "Error unmarshalling data")
	}

	id, err := strconv.Atoi(getEntryRequest.Payload.Id)

	readEntryInput := &app.ReadEntryInput{
		Principle: &app.Principal{
			Type: app.PrincipalUSER,
			ID:   getEntryRequest.Context.Principal.Id,
		},
		ID: id,
	}

	return readEntryInput, nil
}

// EntryToGetEntryResponse mapper
func EntryToGetEntryResponse(e *app.Entry) ([]byte, error) {
	getEntryResponse := &entry.GetEntryResponse{
		Payload: &entry.GetEntryResponse_Payload{
			Id:        fmt.Sprintf("%d", e.ID),
			Text:      e.Text,
			CreatedAt: timeToProtoTime(e.CreatedAt),
			UpdatedAt: timeToProtoTime(e.UpdatedAt),
		},
	}

	output, err := proto.Marshal(getEntryResponse)
	if err != nil {
		errors.Wrap(err, "Wrror marshalling getEntryResponse")
	}

	return output, nil
}

func timeToProtoTime(time time.Time) *timestamp.Timestamp {
	return &timestamp.Timestamp{
		Seconds: time.Unix(),
	}
}
