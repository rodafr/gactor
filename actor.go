package main

import (
	"context"
	"fmt"
	"time"
)

func (actor *actor) Start(ctx context.Context) error {
	actor.ctx, actor.cancel = context.WithCancel(ctx)
	actor.wg.Add(1)
	go func() {
		defer actor.wg.Done()

		for {
			select {
			case <-actor.ctx.Done():
				return
			case msg := <-actor.mailbox:
				if strfunc, ok := actor.behavior.fun.(func(string) string); ok {
					response := strfunc(msg.Content)
					fmt.Printf("actor %q processed: %q -> response: %q \n",
						actor.behavior.actortype, msg.Content, response)
				}
			}
		}
	}()
	return nil
}

func (actor *actor) Stop() error {
	fmt.Printf("%s stopping ...\n", actor.behavior.actortype)
	actor.cancel()
	actor.wg.Wait()
	close(actor.mailbox)
	return nil
}

func (actor *actor) SendMessage(msg Message) error {
	select {
	case actor.mailbox <- msg:
		return nil
	case <-time.After(time.Second):
		return fmt.Errorf("mailbox timeout")
	}
}

func (actor *actor) Execute() {
	fmt.Printf("actor %q executed func \n", actor.id)
	if callfunc, ok := actor.behavior.fun.(func() error); ok {
		err := callfunc()
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}
