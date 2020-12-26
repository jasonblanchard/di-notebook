package main

import (
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type notFound interface {
	NotFound() bool
}

func isNotFound(err error) bool {
	ue, ok := err.(notFound)
	return ok && ue.NotFound()
}

// MapError Maps errors between app and gRPC
func MapError(err error) error {
	switch {
	case isNotFound(errors.Cause(err)):
		return status.Error(codes.NotFound, "Not found")
	default:
		return status.Error(codes.Unknown, "Unknown")
	}
}
