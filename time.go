package c7api

import (
	"errors"
	"time"
)

const (
	TimeFormat = "2006-01-02T15:04:05.000Z"
)

// Takes a date string and formats using time.Parse(layout, date)
//
// Example layout to pass in: layout := "01/02/2006 15:04"
//
// Returns the required format for the API: "2006-01-02T15:04:05.000Z"
func FormatDatesForC7(layout string, date string) (string, error) {
	if date == "" {
		return date, errors.New("date is empty")
	}

	dateFormatted, err := time.Parse(layout, date)
	if err != nil {
		return date, err
	}

	return dateFormatted.Format("2006-01-02T15:04:05.000Z"), err
}
