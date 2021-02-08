package main

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// Principal serializable Principal
type Principal struct {
	Type string
	ID   string
}

// Entry - serializable text entry
type Entry struct {
	ID         string
	Text       string
	CreatorID  string
	CreatedAt  *time.Time
	UpdatedAt  *time.Time
	DeleteTime *time.Time
}

// Revision serializable revision
// This is an "egress interface" in the same way that the gRPC entities are "ingress interface".
// Defining them as such gives us the flexibility to serialize them differently and in different ways than the ingress interface or the internal entity representations,
// giving us an anti-corruption layer between ingress consumers, internal entities and egress consumers.
type Revision struct {
	Actor *Principal
	Entry *Entry
}

func timestamppbToTimePointer(pbtime *timestamppb.Timestamp) *time.Time {
	var t *time.Time
	if pbtime.GetSeconds() != 0 {
		time := pbtime.AsTime()
		t = &time
	}

	return t
}
