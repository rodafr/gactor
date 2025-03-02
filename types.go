package main

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
)

type ActorRegistry interface {
	CreateActor(behavior func(string) string) (*actor, error)
	GetActor(id uuid.UUID) (*actor, error)
	ListActors() []*actor
}

type actorRegistry struct {
	actors map[uuid.UUID]*actor
	mu     sync.RWMutex
}

type Actor interface {
	Start(ctx context.Context) error
	Stop() error
	SendMessage(Message) error
	GetID() uuid.UUID
	Execute()
}

type actor struct {
	id       uuid.UUID
	mailbox  chan Message
	behavior Behaviour
	wg       sync.WaitGroup
	ctx      context.Context
	cancel   context.CancelFunc
}

type Behaviour struct {
	actortype string
	fun       any
}

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
