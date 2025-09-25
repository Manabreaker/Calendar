package eventStore

import (
	"errors"
	"sync"

	"github.com/gibel/Calendar/model"
)

type EventStore struct {
	mu     sync.RWMutex
	events []model.Event
}

func (es *EventStore) Create(key string, value []byte) error {
	es.mu.Lock()
	defer es.mu.Unlock()
	event, err := decodeEvent(value)
	if err != nil {
		return err
	}
	// Проверка на уникальность ID
	for _, e := range es.events {
		if e.ID == key {
			return errors.New("event with this ID already exists")
		}
	}
	event.ID = key
	es.events = append(es.events, event)
	return nil
}

func (es *EventStore) Read(key string) ([]byte, error) {
	es.mu.RLock()
	defer es.mu.RUnlock()
	for _, e := range es.events {
		if e.ID == key {
			return encodeEvent(e)
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
			event.ID = key
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

// Вспомогательные функции для сериализации/десериализации Event
func decodeEvent(data []byte) (model.Event, error) {
	var event model.Event
	// Можно использовать json.Unmarshal, если Event сериализуется в JSON
	// return event, json.Unmarshal(data, &event)
	return event, nil // заглушка
}

func encodeEvent(event model.Event) ([]byte, error) {
	// Можно использовать json.Marshal, если Event сериализуется в JSON
	// return json.Marshal(event)
	return nil, nil // заглушка
}
