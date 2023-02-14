package service

import (
	"encoding/json"
	"log"

	"github.com/TomasCruz/translation-delivery-average/entities"
)

// ProcessInput instantiates and stores an event from each of jsonLines
func (svc Service) ProcessInput(jsonLines []string) error {
	var err error

	for _, line := range jsonLines {
		err = svc.processSingleEvent(line)
		if err != nil {
			return err
		}
	}

	return nil
}

// processSingleEvent sends an event to a corresponding processing function; directing done by this function is dealt by Kafka through topics in a worker app
func (svc Service) processSingleEvent(line string) error {
	var event entities.Event

	err := json.Unmarshal([]byte(line), &event)
	if err != nil {
		return err
	}

	if event.EventName == entities.TranslationDeliveredEventName {
		var translationDeliveredEvent entities.TranslationDeliveredEvent

		err = json.Unmarshal([]byte(line), &translationDeliveredEvent)
		if err != nil {
			return err
		}

		err = svc.processTranslationDeliveredEvent(line, translationDeliveredEvent)
		if err != nil {
			if err == entities.ErrEventIDPresent {
				log.Printf("event %s already processed, skipping\n", translationDeliveredEvent.TranslationID)
				return nil
			} else {
				return err
			}
		}
	}

	return nil
}

// processTranslationDeliveredEvent stores and processes (nothing in this case) the event; this function would be used in the worker app described in main.go comment
func (svc Service) processTranslationDeliveredEvent(line string, tdEvent entities.TranslationDeliveredEvent) error {
	// construct the event; keep original as payload
	event := entities.Event{
		EventID:   tdEvent.TranslationID,
		EventName: tdEvent.EventName,
		EventTS:   tdEvent.Timestamp,
		Payload:   line,
	}

	// store the event
	err := svc.db.StoreTranslationDeliveredEvent(event)
	if err != nil {
		return err
	}

	// process the event (nothing in this case)

	// done
	return nil
}
