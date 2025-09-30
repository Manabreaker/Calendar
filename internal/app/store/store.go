package store

import (
	"github.com/Manabreaker/Calendar/model"
	"time"
)

type Store interface {
	Create(value []byte) error
	Read(key string) ([]byte, error)
	Update(value []byte) error
	Delete(key string) error
	GetEventsInterval(start, end time.Time) []model.Event
	GetEventsToday() []model.Event
	GetEventsWeek() []model.Event
	GetEventsMonth() []model.Event
	GetAllEvents() []model.Event
}
