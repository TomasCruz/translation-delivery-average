package presenter

import "github.com/TomasCruz/translation-delivery-average/service"

// TODO
// PresentAverages marshalls each of Average to JSON, and prints them.
//
// For great number of minutes involved, presenting averages (or sending them to another service for further processing or whatever),
// could be done piecemeal, e.g. 100 rows of calculated averages at a time, retrieved from DB by use of offset and limit in querying
// calculated_averages table, much like pagination is usually done.
func PresentAverages(averages []service.Average) {
}
