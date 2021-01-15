package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBearerHeaderToID(t *testing.T) {
	inputs := []struct {
		input    string
		expected string
		err      string
	}{
		{
			input:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1ZXNyVXVpZCI6ImJjZjRlMzYwLTJiZDQtNDFhMS1hOWQwLTc4NjU3N2UwMmY0YSIsImNzcmZUb2tlbiI6Ink1bXJqYU8zLU1CbWZsWE1EX3dGVDRWbUNJUDZqc0ZWRmVudyIsImlhdCI6MTYxMDY3MjA4Mn0.25KsVVzXwDA__D7fMyspDi4fJad7SPDjaOrcfcqiA2Q",
			expected: "bcf4e360-2bd4-41a1-a9d0-786577e02f4a",
		},
		{
			input:    "bogus",
			expected: "",
			err:      "Invalid token",
		},
	}

	for _, tc := range inputs {
		got, err := bearerHeaderToID(tc.input)
		assert.Equal(t, tc.expected, got, tc)
		if err != nil {
			assert.Equal(t, tc.err, err.Error())
		}
	}

}
