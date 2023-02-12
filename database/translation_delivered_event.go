package database

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/TomasCruz/translation-delivery-average/service"
)

func (pDB postgresDB) StoreTranslationDeliveredEvent(event service.TranslationDeliveredEvent) error {
	pTx, err := pDB.newTransaction()
	if err != nil {
		return err
	}
	defer pTx.commitOrRollbackOnError(&err)

	_, err = pDB.GetEventByID(event.TranslationID)
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
	} else {
		return service.ErrEventIDPresent
	}

	err = pDB.insertTranslationDeliveredEvent(event)
	if err != nil {
		return err
	}

	return nil
}

func (pDB postgresDB) GetEventByID(id string) (service.Event, error) {
	var (
		eventID   string
		eventName string
		eventTS   time.Time
		payload   string
	)

	err := pDB.db.QueryRow(`SELECT event_id, event_name, event_ts, payload FROM events WHERE event_id = $1`, id).
		Scan(&eventID, &eventName, &eventTS, &payload)
	if err != nil {
		return service.Event{}, err
	}

	return service.Event{
		EventID:   eventID,
		EventName: eventName,
		EventTS:   service.MicrosecondTime{T: eventTS},
		Payload:   payload,
	}, nil
}

// TODO
func (pDB postgresDB) ListTranslationDeliveredEvents(endMinute time.Time, windowSize int) ([]service.TranslationDeliveredEvent, error) {
	return nil, nil
}

func (pDB postgresDB) insertTranslationDeliveredEvent(event service.TranslationDeliveredEvent) error {
	sqlStatement := `INSERT INTO events (event_id, event_name, event_ts, payload) VALUES ($1, $2, $3, $4) returning event_id;`

	pTx, err := pDB.newTransaction()
	if err != nil {
		return err
	}
	defer pTx.commitOrRollbackOnError(&err)

	stmt, err := pTx.Prepare(sqlStatement)
	if err != nil {
		return err
	}
	defer stmt.Close()

	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	var id string
	err = stmt.QueryRow(event.TranslationID, event.EventName, event.Timestamp.T, payload).Scan(&id)
	if err != nil {
		return err
	}

	return nil
}
