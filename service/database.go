package service

import "time"

// Database is an interface through which to talk with datastore
type Database interface {
	StoreTranslationDeliveredEvent(event TranslationDeliveredEvent) error
	ListTranslationDeliveredEvents(endMinute time.Time, windowSize int) ([]TranslationDeliveredEvent, error)
	GetEventByID(id string) (Event, error)
	GetFirstTranslationDeliveredEventTime() (time.Time, error)
	GetLastTranslationDeliveredEventTime() (time.Time, error)
}
