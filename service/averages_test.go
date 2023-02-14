package service

import (
	"testing"
	"time"

	"github.com/TomasCruz/translation-delivery-average/entities"
	"github.com/TomasCruz/translation-delivery-average/tests/mocks"
	"github.com/stretchr/testify/assert"
)

// run from terminal (needed for unit tests)
// "go install github.com/vektra/mockery/v2@latest"
// "sudo mv ~/go/bin/mockery /usr/local/go/bin/"
func Test_CalculateAverages(t *testing.T) {
	firstMin, _ := time.Parse(entities.MsLayout, "2018-12-26 18:11:08.509654")
	midMin, _ := time.Parse(entities.MsLayout, "2018-12-26 18:15:19.903159")
	lastMin, _ := time.Parse(entities.MsLayout, "2018-12-26 18:23:19.903159")

	startMinutes := make([]time.Time, 14)
	endMinutes := make([]time.Time, 14)
	startMinutes[0], _ = time.Parse(entities.MsLayout, "2018-12-26 18:02:00.000000")
	endMinutes[0] = startMinutes[0].Add(time.Minute * time.Duration(10))

	for i := 1; i < 14; i++ {
		startMinutes[i] = startMinutes[i-1].Add(time.Minute)
		endMinutes[i] = endMinutes[i-1].Add(time.Minute)
	}

	events := []entities.Event{
		{
			EventID:   "5aa5b2f39f7254a75aa5",
			EventName: entities.TranslationDeliveredEventName,
			EventTS:   entities.MicrosecondTime{T: firstMin},
			Payload:   `{"timestamp": "2018-12-26 18:11:08.509654","translation_id": "5aa5b2f39f7254a75aa5","source_language": "en","target_language": "fr","client_name": "airliberty","event_name": "translation_delivered","nr_words": 30, "duration": 20}`,
		},
		{
			EventID:   "5aa5b2f39f7254a75aa4",
			EventName: entities.TranslationDeliveredEventName,
			EventTS:   entities.MicrosecondTime{T: midMin},
			Payload:   `{"timestamp": "2018-12-26 18:15:19.903159","translation_id": "5aa5b2f39f7254a75aa4","source_language": "en","target_language": "fr","client_name": "airliberty","event_name": "translation_delivered","nr_words": 30, "duration": 31}`,
		},
		{
			EventID:   "5aa5b2f39f7254a75bb3",
			EventName: entities.TranslationDeliveredEventName,
			EventTS:   entities.MicrosecondTime{T: lastMin},
			Payload:   `{"timestamp": "2018-12-26 18:23:19.903159","translation_id": "5aa5b2f39f7254a75bb3","source_language": "en","target_language": "fr","client_name": "taxi-eats","event_name": "translation_delivered","nr_words": 100, "duration": 54}`,
		},
	}

	eventReturns := [][]entities.Event{
		{events[0]},
		{events[0]},
		{events[0]},
		{events[0]},
		{events[0], events[1]},
		{events[0], events[1]},
		{events[0], events[1]},
		{events[0], events[1]},
		{events[0], events[1]},
		{events[0], events[1]},
		{events[1]},
		{events[1]},
		{events[1], events[2]},
	}

	expected := []Average{
		{Date: startMinutes[0].Add(time.Minute * time.Duration(9)), AverageDeliveryTime: 0},
		{Date: startMinutes[0].Add(time.Minute * time.Duration(10)), AverageDeliveryTime: 20},
		{Date: startMinutes[0].Add(time.Minute * time.Duration(11)), AverageDeliveryTime: 20},
		{Date: startMinutes[0].Add(time.Minute * time.Duration(12)), AverageDeliveryTime: 20},
		{Date: startMinutes[0].Add(time.Minute * time.Duration(13)), AverageDeliveryTime: 20},
		{Date: startMinutes[0].Add(time.Minute * time.Duration(14)), AverageDeliveryTime: 25.5},
		{Date: startMinutes[0].Add(time.Minute * time.Duration(15)), AverageDeliveryTime: 25.5},
		{Date: startMinutes[0].Add(time.Minute * time.Duration(16)), AverageDeliveryTime: 25.5},
		{Date: startMinutes[0].Add(time.Minute * time.Duration(17)), AverageDeliveryTime: 25.5},
		{Date: startMinutes[0].Add(time.Minute * time.Duration(18)), AverageDeliveryTime: 25.5},
		{Date: startMinutes[0].Add(time.Minute * time.Duration(19)), AverageDeliveryTime: 25.5},
		{Date: startMinutes[0].Add(time.Minute * time.Duration(20)), AverageDeliveryTime: 31},
		{Date: startMinutes[0].Add(time.Minute * time.Duration(21)), AverageDeliveryTime: 31},
		{Date: startMinutes[0].Add(time.Minute * time.Duration(22)), AverageDeliveryTime: 42.5},
	}

	db := &mocks.Database{}
	svc := NewService(db, 10)

	db.On("GetFirstTranslationDeliveredEventTime").
		Return(firstMin, nil)

	db.On("GetLastTranslationDeliveredEventTime").
		Return(lastMin, nil)

	for i := 0; i < 13; i++ {
		db.On("ListTranslationDeliveredEvents", startMinutes[i], endMinutes[i]).
			Return(eventReturns[i], nil)
	}

	averages, err := svc.CalculateAverages()
	assert.NoError(t, err, "unexpected CalculateAverages error")

	for i := 0; i < 14; i++ {
		assert.Equal(t, expected[i], averages[i])
	}
}
