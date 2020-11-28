package protobufmapper

import (
	"fmt"
	"strconv"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/jasonblanchard/di-messages/packages/go/messages/notebook"
	"github.com/jasonblanchard/di-notebook/app"
	"github.com/jasonblanchard/di-notebook/di_messages/entry"
	"github.com/pkg/errors"
	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/protobuf/proto"
)

// ServiceOrigin service name
const ServiceOrigin = "notebook"

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
		Principal: &app.Principal{
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
	if err != nil {
		return nil, errors.Wrap(err, "Error converting string to int")
	}

	readEntryInput := &app.ReadEntryInput{
		Principal: &app.Principal{
			Type: app.PrincipalUSER,
			ID:   getEntryRequest.Context.Principal.Id,
		},
		ID: id,
	}

	return readEntryInput, nil
}

// ReadEntryRequestToReadEntryInput mapper
func ReadEntryRequestToReadEntryInput(readEntryRequest *notebook.ReadEntryRequest) (*app.ReadEntryInput, error) {
	id, err := strconv.Atoi(readEntryRequest.GetPayload().Id)
	if err != nil {
		return nil, errors.Wrap(err, "Error converting string to int")
	}

	readEntryInput := &app.ReadEntryInput{
		Principal: &app.Principal{
			Type: app.PrincipalUSER,
			ID:   readEntryRequest.GetContext().GetPrincipal().Id,
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

	if !e.UpdatedAt.IsZero() {
		getEntryResponse.Payload.UpdatedAt = timeToProtoTime(e.UpdatedAt)
	}

	output, err := proto.Marshal(getEntryResponse)
	if err != nil {
		errors.Wrap(err, "Wrror marshalling getEntryResponse")
	}

	return output, nil
}

// EntryToReadEntryResponse mapper
func EntryToReadEntryResponse(e *app.Entry, traceID string) ([]byte, error) {
	readEntryResponse := &notebook.ReadEntryResponse{
		Payload: &notebook.ReadEntryResponse_Payload{
			Id:        fmt.Sprintf("%d", e.ID),
			Text:      e.Text,
			CreatedAt: timeToProtoTime(e.CreatedAt),
			UpdatedAt: timeToNullableProtoTime(e.UpdatedAt),
		},
		Context: &notebook.ResponseContext{
			TraceId: traceID,
			Origin:  ServiceOrigin,
		},
	}

	output, err := proto.Marshal(readEntryResponse)
	if err != nil {
		errors.Wrap(err, "Wrror marshalling getEntryResponse")
	}

	return output, nil
}

func timeToProtoTime(time time.Time) *timestamp.Timestamp {
	seconds := time.Unix()

	if time.IsZero() {
		seconds = 0
	}

	return &timestamp.Timestamp{
		Seconds: seconds,
	}
}

func timeToNullableProtoTime(time time.Time) *notebook.NullableTimestamp {
	if time.IsZero() {
		return &notebook.NullableTimestamp{
			Value: &notebook.NullableTimestamp_Null{},
		}
	}
	return &notebook.NullableTimestamp{
		Value: &notebook.NullableTimestamp_Timestamp{
			Timestamp: timeToProtoTime(time),
		},
	}
}

// DeleteEntryRequestToDiscardEntryInput mapper
func DeleteEntryRequestToDiscardEntryInput(data []byte) (*app.DiscardEntryInput, error) {
	deleteEntryRequest := &entry.DeleteEntryRequest{}
	err := proto.Unmarshal(data, deleteEntryRequest)
	if err != nil {
		return nil, errors.Wrap(err, "Error unmarshalling request")
	}

	id, err := strconv.Atoi(deleteEntryRequest.Payload.Id)
	if err != nil {
		return nil, errors.Wrap(err, "Error converting string to int")
	}

	discardEntryInput := &app.DiscardEntryInput{
		Principal: &app.Principal{
			Type: app.PrincipalUSER,
			ID:   deleteEntryRequest.Context.Principal.Id,
		},
		ID: id,
	}

	return discardEntryInput, nil
}

// DiscardEntryToDeleteEntryResponse mapper
func DiscardEntryToDeleteEntryResponse() ([]byte, error) {
	response := &entry.DeleteEntryResponse{}
	data, err := proto.Marshal(response)
	if err != nil {
		return nil, errors.Wrap(err, "Error marshalling response")
	}
	return data, nil
}

// UpdateEntryRequestToChangeEntryInput mapper
func UpdateEntryRequestToChangeEntryInput(data []byte) (*app.ChangeEntryInput, error) {
	updateEntryRequest := &entry.UpdateEntryRequest{}
	err := proto.Unmarshal(data, updateEntryRequest)
	if err != nil {
		return nil, errors.Wrap(err, "Error unmarshalling request")
	}

	id, err := strconv.Atoi(updateEntryRequest.Payload.Id)
	if err != nil {
		return nil, errors.Wrap(err, "Error converting string to int")
	}

	changeEntryInput := &app.ChangeEntryInput{
		Principal: &app.Principal{
			Type: app.PrincipalUSER,
			ID:   updateEntryRequest.Context.Principal.Id,
		},
		ID:   id,
		Text: updateEntryRequest.Payload.Text,
	}

	return changeEntryInput, nil
}

// ChangeEntryOutputToUpdateEntryResponse mapper
func ChangeEntryOutputToUpdateEntryResponse(e *app.Entry) ([]byte, error) {
	updateEntryResponse := &entry.UpdateEntryResponse{
		Payload: &entry.UpdateEntryResponse_Payload{
			Id:        fmt.Sprintf("%d", e.ID),
			Text:      e.Text,
			CreatedAt: timeToProtoTime(e.CreatedAt),
			UpdatedAt: timeToProtoTime(e.UpdatedAt),
		},
	}

	output, err := proto.Marshal(updateEntryResponse)
	if err != nil {
		errors.Wrap(err, "Wrror marshalling updateEntryResponse")
	}

	return output, nil
}

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

// ListEntriesRequestToListEntriesInput mapper
func ListEntriesRequestToListEntriesInput(data []byte) (*app.ListEntriesInput, error) {
	request := &entry.ListEntriesRequest{}
	err := proto.Unmarshal(data, request)
	if err != nil {
		return nil, errors.Wrap(err, "Error unmarshalling request")
	}

	var after int
	if request.Payload.After == "" {
		after = 0
	} else {
		after, err = strconv.Atoi(request.Payload.After)
	}

	if err != nil {
		return nil, errors.Wrap(err, "Error converting string to int")
	}

	listEntriesInput := &app.ListEntriesInput{
		Principal: &app.Principal{
			Type: app.PrincipalTEST,
			ID:   request.Context.Principal.Id,
		},
		CreatorID: request.Payload.CreatorId,
		First:     int(request.Payload.First),
		After:     after,
	}

	return listEntriesInput, nil
}

// ListEntriesOutputToListEntriesResponse mapper
func ListEntriesOutputToListEntriesResponse(i *app.ListEntriesOutput) ([]byte, error) {
	response := &entry.ListEntriesResponse{
		PageInfo: &entry.ListEntriesResponse_PageInfo{
			TotalCount:  int32(i.Pagination.TotalCount),
			HasNextPage: i.Pagination.HasNextPage,
			StartCursor: fmt.Sprintf("%d", i.Pagination.StartCursor),
			EndCursor:   fmt.Sprintf("%d", i.Pagination.EndCursor),
		},
	}

	for _, entryInstance := range i.Entries {
		entity := &entry.ListEntriesResponse_Entity{
			Id:        fmt.Sprintf("%d", entryInstance.ID),
			Text:      entryInstance.Text,
			CreatorId: entryInstance.CreatorID,
			CreatedAt: timeToProtoTime(entryInstance.CreatedAt),
			UpdatedAt: timeToProtoTime(entryInstance.UpdatedAt),
		}

		response.Payload = append(response.Payload, entity)
	}

	responseData, err := proto.Marshal(response)
	if err != nil {
		return nil, errors.Wrap(err, "Error marshalling response")
	}

	return responseData, nil
}

// ToNotebookErrorResponse mapper
func ToNotebookErrorResponse(message string, code code.Code, traceID string) ([]byte, error) {
	errorResponse := &notebook.ErrorResponse{
		Status: &notebook.Status{
			Code:    int32(code),
			Message: message,
		},
		Context: &notebook.ResponseContext{
			TraceId: traceID,
			Origin:  ServiceOrigin,
		},
	}

	return proto.Marshal(errorResponse)
}
