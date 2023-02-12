package service

import "time"

const TranslationDeliveredEventName string = "translation_delivered"

type TranslationDeliveredEvent struct {
	Timestamp      time.Time `json:"timestamp" example:"2018-12-26 18:12:19.903159"`
	TranslationID  string    `json:"translation_id" example:"5aa5b2f39f7254a75aa4"`
	SourceLanguage string    `json:"source_language" example:"en"`
	TargetLanguage string    `json:"target_language" example:"fr"`
	ClientName     string    `json:"client_name" example:"airliberty"`
	EventName      string    `json:"event_name" example:"translation_delivered"`
	Duration       int       `json:"duration" example:"20"`
	NrWords        int       `json:"nr_words" example:"100"`
}
