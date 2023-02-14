package service

import (
	"time"

	"github.com/TomasCruz/translation-delivery-average/entities"
)

// Database is an interface through which to talk with datastore
type Database interface {
	StoreTranslationDeliveredEvent(event entities.Event) error
	ListTranslationDeliveredEvents(startMinute, endMinute time.Time) ([]entities.Event, error)
	GetEventByID(id string) (entities.Event, error)
	GetFirstTranslationDeliveredEventTime() (time.Time, error)
	GetLastTranslationDeliveredEventTime() (time.Time, error)
}
