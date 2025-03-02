package main

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Behaviour struct {
	actortype string
	fun       any
}

func echoFunc(msg string) string {
	return fmt.Sprintf("[echoer] echo: %s", msg)
}

func sendFunc(r *actorRegistry) func() error {
	return func() error {
		fmt.Println("[sender] triggered")
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
		fmt.Printf("[sender] echoerID: %s\n", echoerID)

		echoer, err := r.GetActor(echoerID)
		if err != nil {
			fmt.Printf("[sender] sender fail: %s\n", err.Error())
			return err
		}
		err = echoer.SendMessage(Message{
			Type:      NormalMsg,
			Content:   "intraactormsg",
			SenderID:  senderID,
			Timestamp: time.Now(),
		})
		if err != nil {
			panic(err)
		}

		fmt.Printf("[sender] messaged %q\n", echoerID)
		return nil
	}
}
