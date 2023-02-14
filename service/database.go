package service

import "time"

// Database is an interface through which to talk with datastore
type Database interface {
	StoreTranslationDeliveredEvent(event Event) error
	ListTranslationDeliveredEvents(startMinute, endMinute time.Time) ([]Event, error)
	GetEventByID(id string) (Event, error)
	GetFirstTranslationDeliveredEventTime() (time.Time, error)
	GetLastTranslationDeliveredEventTime() (time.Time, error)
}
