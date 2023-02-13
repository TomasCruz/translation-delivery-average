package service

import (
	"fmt"
	"strings"
	"time"
)

const (
	MinLayout = "2006-01-02 15:04:05"
)

type MinuteTime struct {
	time.Time
}

func (ct *MinuteTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		ct = &MinuteTime{}
		return
	}

	t, err := time.Parse(MinLayout, s)
	ct = &MinuteTime{t}
	return
}

func (ct *MinuteTime) MarshalJSON() ([]byte, error) {
	if ct.UnixNano() == int64(0) {
		return []byte("null"), nil
	}

	return []byte(fmt.Sprintf("\"%s\"", ct.Format(MinLayout))), nil
}

type Average struct {
	Date                time.Time `json:"date" example:"2018-12-26 18:12:19"`
	AverageDeliveryTime float64   `json:"average_delivery_time" example:"25.5"`
}

// CalculateAverages extract times of first and last event from the DB,
// then for each minute between first and last time calculates average for the particular minute, and returns slice of Average
//
// The sequential approach used here only works for a limited number of minutes involved. For great number,
// the results of the calculation would be stored in it's own DB table (call it calculated_averages),
// one row per minute, instead of being returned from the function.
func (svc Service) CalculateAverages() ([]Average, error) {
	firstMin, err := svc.db.GetFirstTranslationDeliveredEventTime()
	if err != nil {
		return nil, err
	}
	firstMinRounded := time.Date(firstMin.Year(), firstMin.Month(), firstMin.Day(), firstMin.Hour(), firstMin.Minute(), 0, 0, firstMin.Location())

	lastMin, err := svc.db.GetLastTranslationDeliveredEventTime()
	if err != nil {
		return nil, err
	}
	lastMinRounded := time.Date(lastMin.Year(), lastMin.Month(), lastMin.Day(), lastMin.Hour(), lastMin.Minute(), 0, 0, lastMin.Location())

	// first one is the special case, i.e. no events,
	// average for that one is 0 instead of duration_sum / nr_preceding_events, (avoid division by 0)
	numMinutesToProcess := lastMinRounded.Sub(firstMinRounded)/time.Minute + 2 // additional one is for average of duration for last minute, i.e. last + 1 - first + 1
	averages := make([]Average, 0, numMinutesToProcess)
	averages = append(averages, Average{Date: firstMinRounded, AverageDeliveryTime: 0})

	currentlyProcessedMinute := firstMinRounded
	for currentlyProcessedMinute.Before(lastMin) {
		currentlyProcessedMinute = currentlyProcessedMinute.Add(time.Minute * time.Duration(1))

		currMinusWindowSize := currentlyProcessedMinute.Add(-time.Minute * time.Duration(svc.windowSize))
		events, err := svc.db.ListTranslationDeliveredEvents(currMinusWindowSize, currentlyProcessedMinute)
		if err != nil {
			return nil, err
		}

		durationSum := float64(0)
		for _, ev := range events {
			tdEvent, err := NewTranslationDeliveredEventFromEvent(ev)
			if err != nil {
				return nil, err
			}

			durationSum += float64(tdEvent.Duration)
		}

		averages = append(averages, Average{
			Date:                currentlyProcessedMinute,
			AverageDeliveryTime: durationSum / float64(len(events)),
		})
	}

	return averages, nil
}
