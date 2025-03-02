package main

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
)

type ActorRegistry interface {
	CreateActor(behavior func(string) string) (*actor, error)
	GetActor(id uuid.UUID) (*actor, error)
	ListActors() map[uuid.UUID]string
}

type actorRegistry struct {
	actors map[uuid.UUID]*actor
	mu     sync.RWMutex
}

func (r *actorRegistry) CreateActor(behavior Behaviour) (*actor, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := uuid.New()
	a := &actor{
		id:       id,
		mailbox:  make(chan Message),
		behavior: behavior,
		wg:       sync.WaitGroup{},
		ctx:      nil,
		cancel:   nil,
	}

	r.actors[id] = a

	return a, nil
}

func (r *actorRegistry) CreateActors(behaviors []Behaviour) error {
	for _, b := range behaviors {
		_, err := r.CreateActor(b)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *actorRegistry) GetActor(id uuid.UUID) (*actor, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	a, ok := r.actors[id]
	if !ok {
		return nil, fmt.Errorf("no such actor: %s", id)
	}
	return a, nil
}

func (r *actorRegistry) ListActors() map[uuid.UUID]string {
	r.mu.Lock()
	defer r.mu.Unlock()

	as := make(map[uuid.UUID]string, len(r.actors))
	for u, a := range r.actors {
		as[u] = a.behavior.actortype
	}
	return as
}
