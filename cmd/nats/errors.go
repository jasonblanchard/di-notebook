package main

import (
	"google.golang.org/genproto/googleapis/rpc/code"
)

type unauthorized interface {
	Unauthorized() bool
}

func isUnauthorized(err interface{}) bool {
	ue, ok := err.(unauthorized)
	return ok && ue.Unauthorized()
}

func errorToCode(err interface{}) code.Code {
	if isUnauthorized(err) {
		return code.Code_PERMISSION_DENIED
	}

	return code.Code_UNKNOWN
}
