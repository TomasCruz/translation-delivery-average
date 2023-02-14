package entities

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

var ErrEventIDPresent error = errors.New("event with the given ID already exists")

const (
	TranslationDeliveredEventName string = "translation_delivered"
	MsLayout                      string = "2006-01-02 15:04:05.000000"
)

type MicrosecondTime struct {
	T time.Time
}

func (ct *MicrosecondTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		ct.T = time.Time{}
		return
	}

	ct.T, err = time.Parse(MsLayout, s)
	return
}

func (ct *MicrosecondTime) MarshalJSON() ([]byte, error) {
	if ct.T.UnixNano() == int64(0) {
		return []byte("null"), nil
	}

	return []byte(fmt.Sprintf("\"%s\"", ct.T.Format(MsLayout))), nil
}

type Event struct {
	EventID   string          `json:"event_id" example:"5aa5b2f39f7254a75aa4"`
	EventName string          `json:"event_name" example:"translation_delivered"`
	EventTS   MicrosecondTime `json:"event_ts" example:"2018-12-26 18:12:19.903159"`
	Payload   string          `json:"payload" example:"{\"timestamp\": \"2018-12-26 18:11:08.509654\",\"translation_id\": \"5aa5b2f39f7254a75aa5\",\"source_language\": \"en\",\"target_language\": \"fr\",\"client_name\": \"airliberty\",\"event_name\": \"translation_delivered\",\"nr_words\": 30, \"duration\": 20}"`
}

type TranslationDeliveredEvent struct {
	Timestamp      MicrosecondTime `json:"timestamp" example:"2018-12-26 18:12:19.903159"`
	TranslationID  string          `json:"translation_id" example:"5aa5b2f39f7254a75aa4"`
	SourceLanguage string          `json:"source_language" example:"en"`
	TargetLanguage string          `json:"target_language" example:"fr"`
	ClientName     string          `json:"client_name" example:"airliberty"`
	EventName      string          `json:"event_name" example:"translation_delivered"`
	Duration       int             `json:"duration" example:"20"`
	NrWords        int             `json:"nr_words" example:"100"`
}

func NewTranslationDeliveredEventFromEvent(ev Event) (TranslationDeliveredEvent, error) {
	if ev.EventName != TranslationDeliveredEventName {
		return TranslationDeliveredEvent{}, errors.New(fmt.Sprintf("event %s not of '%s' type", ev.EventID, TranslationDeliveredEventName))
	}

	var tdEvent TranslationDeliveredEvent

	err := json.Unmarshal([]byte(ev.Payload), &tdEvent)
	if err != nil {
		return TranslationDeliveredEvent{}, err
	}

	return TranslationDeliveredEvent{
		Timestamp:      ev.EventTS,
		TranslationID:  ev.EventID,
		SourceLanguage: tdEvent.SourceLanguage,
		TargetLanguage: tdEvent.TargetLanguage,
		ClientName:     tdEvent.ClientName,
		EventName:      ev.EventName,
		Duration:       tdEvent.Duration,
		NrWords:        tdEvent.NrWords,
	}, nil
}
