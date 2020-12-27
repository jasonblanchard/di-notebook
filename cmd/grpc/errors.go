package main

import (
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type notFound interface {
	NotFound() bool
}

type unauthorized interface {
	Unauthorized() bool
}

func isUnauthorized(err error) bool {
	ue, ok := err.(unauthorized)
	return ok && ue.Unauthorized()
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
	case isUnauthorized(errors.Cause(err)):
		return status.Error(codes.PermissionDenied, "Permission denied")
	default:
		return status.Error(codes.Unknown, "Unknown")
	}
}
