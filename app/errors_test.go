package app

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type unauthorized interface {
	Unauthorized() bool
}

func isUnauthorized(err error) bool {
	ue, ok := err.(unauthorized)
	return ok && ue.Unauthorized()
}

func TestUnauthorized(t *testing.T) {
	tests := []struct {
		input error
		want  bool
	}{
		{input: errors.New("oops"), want: false},
		{input: &UnauthorizedError{s: "Can't touch this"}, want: true},
	}

	for i, tc := range tests {
		got := isUnauthorized(tc.input)
		assert.Equal(t, got, tc.want, fmt.Sprintf("case: %v, input: %v, want: %v", i, tc.input, tc.want))
	}
}
