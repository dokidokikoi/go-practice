package main

import "github.com/google/wire"

func InitializeEvent(phrase string) (Event, error) {
	wire.Build(ModelSet, NewEvent, NewMessage)
	return Event{}, nil
}
