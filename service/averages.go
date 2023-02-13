package service

import (
	"fmt"
	"strings"
	"time"
)

const (
	minLayout = "2006-01-02 15:04:05"
)

type MinuteTime struct {
	T time.Time
}

func (ct *MinuteTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		ct.T = time.Time{}
		return
	}

	ct.T, err = time.Parse(minLayout, s)
	return
}

func (ct *MinuteTime) MarshalJSON() ([]byte, error) {
	if ct.T.UnixNano() == int64(0) {
		return []byte("null"), nil
	}

	return []byte(fmt.Sprintf("\"%s\"", ct.T.Format(minLayout))), nil
}

type Average struct {
	Date                MinuteTime `json:"date" example:"2018-12-26 18:12:19"`
	AverageDeliveryTime int        `json:"average_delivery_time" example:"20"`
}

// TODO
// CalculateAverages extract times of first and last event from the DB,
// then for each minute between first and last time calculates average for the particular minute, and returns slice of Average
//
// The sequential approach used here only works for a limited number of minutes involved. For great number,
// the results of the calculation would be stored in it's own DB table (call it calculated_averages),
// one row per minute, instead of being returned from the function.
func (svc Service) CalculateAverages() []Average {
	return nil
}
