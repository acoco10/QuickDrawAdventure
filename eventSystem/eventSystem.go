package events

import (
	"reflect"
)

type Event interface{}

type Handler func(Event)

type EventHub struct {
	subscribers map[string][]Handler
}

func NewEventHub() *EventHub {
	return &EventHub{
		subscribers: make(map[string][]Handler),
	}
}

func (h *EventHub) Subscribe(eventType Event, handler Handler) {
	t := reflect.TypeOf(eventType).String()
	h.subscribers[t] = append(h.subscribers[t], handler)
}

func (h *EventHub) Publish(event Event) {
	t := reflect.TypeOf(event).String()
	for _, handler := range h.subscribers[t] {
		handler(event)
	}
}

type PropertyUpdate struct {
	Property string
	Value    bool
}

type DialogueEvent struct {
	Characters []string
}
