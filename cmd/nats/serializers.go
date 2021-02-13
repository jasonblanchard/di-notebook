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
	ID         string     `json:"id"`
	Text       string     `json:"text"`
	CreatorID  string     `json:"creator_id"`
	CreatedAt  *time.Time `json:"create_time"`
	UpdatedAt  *time.Time `json:"update_time"`
	DeleteTime *time.Time `json:"delete_time"`
	ActorType  string     `json:"actor_type"`
	ActorID    string     `json:"actor_id"`
}

func timestamppbToTimePointer(pbtime *timestamppb.Timestamp) *time.Time {
	var t *time.Time
	if pbtime.GetSeconds() != 0 {
		time := pbtime.AsTime()
		t = &time
	}

	return t
}
