package presenter

import (
	"fmt"

	"github.com/TomasCruz/translation-delivery-average/service"
)

// PresentAverages marshalls each of Average to JSON, and prints them.
//
// For great number of minutes involved, presenting averages (or sending them to another service for further processing or whatever),
// could be done piecemeal, e.g. 100 rows of calculated averages at a time, retrieved from DB by use of offset and limit in querying
// calculated_averages table, much like pagination is usually done.
func PresentAverages(averages []service.Average) error {
	for _, a := range averages {
		s := fmt.Sprintf(`{"date":"%s","average_delivery_time":%f}`, a.Date.Format(service.MinLayout), a.AverageDeliveryTime)
		fmt.Println(s)
	}

	return nil
}
