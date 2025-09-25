package eventStore

import (
	"encoding/json"
	"errors"
	"sync"

	"github.com/Manabreaker/Calendar/model"
)

type EventStore struct {
	mu     sync.RWMutex
	events []model.Event
}

func (es *EventStore) Create(value []byte) error {
	es.mu.Lock()
	defer es.mu.Unlock()
	event, err := decodeEvent(value)
	if err != nil {
		return err
	}
	// Проверка на уникальность ID
	for _, e := range es.events {
		if e.ID == event.ID {
			return errors.New("event with this ID already exists")
		}
	}
	es.events = append(es.events, event)
	return nil
}

func (es *EventStore) Read(key string) ([]byte, error) {
	es.mu.RLock()
	defer es.mu.RUnlock()
	for _, e := range es.events {
		if e.ID == key {
			return json.Marshal(e)
		}
	}
	return nil, errors.New("event not found")
}

func (es *EventStore) Update(key string, value []byte) error {
	es.mu.Lock()
	defer es.mu.Unlock()
	event, err := decodeEvent(value)
	if err != nil {
		return err
	}
	for i, e := range es.events {
		if e.ID == key {
			es.events[i] = event
			return nil
		}
	}
	return errors.New("event not found")
}

func (es *EventStore) Delete(key string) error {
	es.mu.Lock()
	defer es.mu.Unlock()
	for i, e := range es.events {
		if e.ID == key {
			es.events = append(es.events[:i], es.events[i+1:]...)
			return nil
		}
	}
	return errors.New("event not found")
}

func decodeEvent(data []byte) (model.Event, error) {
	var event model.Event
	err := json.Unmarshal(data, &event)
	if err != nil {
		return model.Event{}, err
	}
	// Дополнительная валидация
	if !event.Validate() {
		return model.Event{}, errors.New("invalid event data")
	}

	return event, nil
}
