package main

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/google/uuid"
)

func main() {
	// Create registry to keep track of actors
	r := &actorRegistry{
		actors: make(map[uuid.UUID]*actor),
	}

	// Define wanted behaviors
	bs := []Behaviour{
		{
			actortype: "echoMsgActor",
			fun:       echoFunc,
		},
		{
			actortype: "sendMsgActor",
			fun:       sendFunc(r),
		},
	}

	err := r.CreateActors(bs)
	if err != nil {
		panic(err)
	}

	// List actors
	var senderID, echoerID uuid.UUID
	actorIDs := r.ListActors()
	for a, t := range actorIDs {
		if t == "echoMsgActor" {
			echoerID = a
		}
		if t == "sendMsgActor" {
			senderID = a
		}
	}

	fmt.Printf("echoer: %s\nsender: %s\n", echoerID, senderID)
	fmt.Printf("goroutines: %d\n", runtime.NumGoroutine())
	fmt.Println("starting actors ...")

	// Start actors
	for actor := range actorIDs {
		a, err := r.GetActor(actor)
		if err != nil {
			panic(err)
		}
		err = a.Start(context.Background())
		if err != nil {
			panic(err)
		}
	}
	fmt.Printf("goroutines: %d\n", runtime.NumGoroutine())

	// call execute on sendMsgActor
	sender, err := r.GetActor(senderID)
	if err != nil {
		panic(err)
	}
	sender.Execute()

	fmt.Printf("goroutines: %d\n", runtime.NumGoroutine())

	time.Sleep(time.Second * 2) // Let processing complete

	fmt.Println("stopping actors ...")
	// stop actors
	for actor := range actorIDs {
		a, err := r.GetActor(actor)
		if err != nil {
			panic(err)
		}
		err = a.Stop()
		if err != nil {
			panic(err)
		}
	}
	fmt.Printf("goroutines: %d\n", runtime.NumGoroutine())
}
