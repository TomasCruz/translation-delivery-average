package service

import "time"

type Average struct {
	Date                time.Time `json:"date" example:"2018-12-26 18:12:19"`
	AverageDeliveryTime int       `json:"average_delivery_time" example:"20"`
}

// TODO
// CalculateAverages extract times of first and last event from the DB,
// then for each minute between first and last time calculates average for the particular minute, and returns slice of Average
func (svc Service) CalculateAverages() []Average {
	return nil
}
