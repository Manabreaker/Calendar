package eventStore

import (
	"encoding/json"
	"errors"
	"github.com/Manabreaker/Calendar/model"
	"sync"
	"time"
)

var (
	errorNotFound           = errors.New("event not found")
	errorEventAlreadyExists = errors.New("event with this ID already exists")
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
			return errorEventAlreadyExists
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
	return nil, errorNotFound
}

func (es *EventStore) Update(value []byte) error {
	event, err := decodeEvent(value)
	if err != nil {
		return err
	}
	es.mu.Lock()
	defer es.mu.Unlock()
	for i, e := range es.events {
		if e.ID == event.ID {
			es.events[i] = event
			return nil
		}
	}
	return errorNotFound
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
	return errorNotFound
}

// GetEventsInterval возвращает все события из интервала [start, end]
func (es *EventStore) GetEventsInterval(start, end time.Time) []model.Event {
	var result []model.Event
	es.mu.RLock()
	defer es.mu.RUnlock()
	// формат даты в событии
	layout := "2006-01-02"

	for _, e := range es.events {
		eventTime, err := time.Parse(layout, e.Date)
		if err != nil {
			continue
		}

		// проверяем попадание в интервал включительно
		if (eventTime.Equal(start) || eventTime.After(start)) &&
			(eventTime.Equal(end) || eventTime.Before(end)) {
			result = append(result, e)
		}
	}

	return result
}

func (es *EventStore) GetEventsToday() []model.Event {
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	end := start.AddDate(0, 0, 1).Add(-time.Nanosecond)
	return es.GetEventsInterval(start, end)
}

func (es *EventStore) GetEventsWeek() []model.Event {
	now := time.Now()
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7 // Воскресенье
	}
	start := time.Date(now.Year(), now.Month(), now.Day()-weekday+1, 0, 0, 0, 0, now.Location())
	end := start.AddDate(0, 0, 7).Add(-time.Nanosecond)
	return es.GetEventsInterval(start, end)
}

func (es *EventStore) GetEventsMonth() []model.Event {
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	end := start.AddDate(0, 1, 0).Add(-time.Nanosecond)
	return es.GetEventsInterval(start, end)
}

func (es *EventStore) GetAllEvents() []model.Event {
	var result []model.Event
	es.mu.RLock()
	defer es.mu.RUnlock()

	for _, e := range es.events {
		result = append(result, e)

	}

	return result
}

func NewStore() *EventStore {
	return &EventStore{
		events: make([]model.Event, 0),
	}
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
