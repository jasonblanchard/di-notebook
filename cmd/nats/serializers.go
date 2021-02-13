package main

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// EntryRevision serializable revision
// This is an "egress interface" in the same way that the gRPC entities are "ingress interface".
// Defining them as such gives us the flexibility to serialize them differently and in different ways than the ingress interface or the internal entity representations,
// giving us an anti-corruption layer between ingress consumers, internal entities and egress consumers.
// In this case, the Entry and Principal entities are flattened for simpler analytics consumption. Arguable whether or not this lives here, but... it does for now.
type EntryRevision struct {
	ID         string
	Text       string
	CreatorID  string
	CreatedAt  *time.Time
	UpdatedAt  *time.Time
	DeleteTime *time.Time
	ActorType  string
	ActorID    string
}

func timestamppbToTimePointer(pbtime *timestamppb.Timestamp) *time.Time {
	var t *time.Time
	if pbtime.GetSeconds() != 0 {
		time := pbtime.AsTime()
		t = &time
	}

	return t
}
