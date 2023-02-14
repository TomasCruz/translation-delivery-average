package database

import (
	"database/sql"
	"time"

	"github.com/TomasCruz/translation-delivery-average/entities"
)

func (pDB postgresDB) StoreTranslationDeliveredEvent(event entities.Event) error {
	pTx, err := pDB.newTransaction()
	if err != nil {
		return err
	}
	defer pTx.commitOrRollbackOnError(&err)

	_, err = pDB.GetEventByID(event.EventID)
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
	} else {
		return entities.ErrEventIDPresent
	}

	err = pDB.insertEvent(event)
	if err != nil {
		return err
	}

	return nil
}

func (pDB postgresDB) GetEventByID(id string) (entities.Event, error) {
	var (
		eventID   string
		eventName string
		eventTS   time.Time
		payload   string
	)

	err := pDB.db.QueryRow(`SELECT event_id, event_name, event_ts, payload FROM events WHERE event_id = $1`, id).
		Scan(&eventID, &eventName, &eventTS, &payload)
	if err != nil {
		return entities.Event{}, err
	}

	return entities.Event{
		EventID:   eventID,
		EventName: eventName,
		EventTS:   entities.MicrosecondTime{T: eventTS},
		Payload:   payload,
	}, nil
}

func (pDB postgresDB) ListTranslationDeliveredEvents(startMinute, endMinute time.Time) ([]entities.Event, error) {
	sqlQuery := `SELECT event_id, event_name, event_ts, payload
				 FROM	events
				 WHERE	event_name = $1 AND
				 		event_ts >= $2 AND
						event_ts <  $3`

	rows, err := pDB.db.Query(sqlQuery, entities.TranslationDeliveredEventName, startMinute, endMinute)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		eventID   string
		eventName string
		eventTS   time.Time
		payload   string
	)

	var events []entities.Event
	for rows.Next() {
		if err = rows.Scan(&eventID, &eventName, &eventTS, &payload); err != nil {
			return nil, err
		}

		events = append(events, entities.Event{
			EventID:   eventID,
			EventName: eventName,
			EventTS:   entities.MicrosecondTime{T: eventTS},
			Payload:   payload,
		})
	}

	return events, nil
}

func (pDB postgresDB) insertEvent(event entities.Event) error {
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

	var id string
	err = stmt.QueryRow(event.EventID, event.EventName, event.EventTS.T, event.Payload).Scan(&id)
	if err != nil {
		return err
	}

	return nil
}

func (pDB postgresDB) GetFirstTranslationDeliveredEventTime() (time.Time, error) {
	var eventTS time.Time

	err := pDB.db.QueryRow(`SELECT MIN(event_ts) FROM events`).Scan(&eventTS)
	if err != nil {
		return time.Time{}, err
	}

	return eventTS, nil
}

func (pDB postgresDB) GetLastTranslationDeliveredEventTime() (time.Time, error) {
	var eventTS time.Time

	err := pDB.db.QueryRow(`SELECT MAX(event_ts) FROM events`).Scan(&eventTS)
	if err != nil {
		return time.Time{}, err
	}

	return eventTS, nil
}
