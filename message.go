package main

import (
	"time"

	"github.com/google/uuid"
)

type MessageType string

const (
	NormalMsg MessageType = "normal"
	SystemMsg MessageType = "system"
	ErrorMsg  MessageType = "error"
)

type Message struct {
	Type      MessageType
	Content   string
	SenderID  uuid.UUID
	Timestamp time.Time
}
