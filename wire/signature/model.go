package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/wire"
)

type Message string

type Greeter interface {
	Greeter() Message
	Grumpy() bool
}

type myGreeter struct {
	Message Message
	grumpy  bool
}

func (g myGreeter) Greeter() Message {
	if g.Grumpy() {
		return Message("Go away!")
	}
	return g.Message
}

func (g myGreeter) Grumpy() bool {
	return g.grumpy
}

type Event struct {
	Greeter Greeter
}

func (e Event) Start() {
	msg := e.Greeter.Greeter()
	fmt.Println(msg)
}

func NewMessage(phrase string) Message {
	return Message(phrase)
}

func NewGreeter(m Message) Greeter {
	var grumpy bool
	if time.Now().Unix()%2 == 0 {
		grumpy = true
	}
	return myGreeter{Message: m, grumpy: grumpy}
}

func NewEvent(g Greeter) (Event, error) {
	if g.Grumpy() {
		return Event{}, errors.New("could not create event: event greeter is grumpy")
	}
	return Event{Greeter: g}, nil
}

var ModelSet = wire.NewSet(NewGreeter)
