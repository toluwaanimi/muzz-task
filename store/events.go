package store

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

type eventStore struct {
	eventsChannel chan Event
	subscribers   map[string]SubscriptionHandler
	logger        *logrus.Logger
}

func (e eventStore) Publish(topic string, data []byte) error {
	event := Event{topic: topic, data: data}
	e.eventsChannel <- event
	return nil
}

func (e eventStore) Subscribe(topic string, handler SubscriptionHandler) error {
	if _, ok := e.subscribers[topic]; ok {
		return fmt.Errorf("duplicate subscription for topic: %s", topic)
	}

	e.subscribers[topic] = handler
	return nil
}

func NewEventStore(logger *logrus.Logger) EventStore {
	evStore := &eventStore{
		eventsChannel: make(chan Event, 10),
		subscribers:   make(map[string]SubscriptionHandler),
		logger:        logger,
	}

	go evStore.eventsToRespectiveSubscriptionHandlerConsumer()
	return evStore
}

// eventsToRespectiveSubscriptionHandlerConsumer does the following:
func (e eventStore) eventsToRespectiveSubscriptionHandlerConsumer() {
	for {
		event, ok := <-e.eventsChannel
		if !ok {
			continue
		}

		handler, ok := e.subscribers[event.topic]
		if !ok {
			e.logger.Warnf("No registered subscriber for %s topic", event.topic)
			continue
		}

		go func(ev Event, handler SubscriptionHandler) {
			if err := handler(ev); err != nil {
				e.logger.WithError(err).Errorf("%s returned an error when handling event.", ev.Topic())
			}
		}(event, handler)
	}
}
