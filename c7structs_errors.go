package c7api

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type C7Error struct {
	StatusCode int              `json:"statusCode"`
	Type       string           `json:"type"`
	Message    string           `json:"message"`
	Errors     []map[string]any `json:"errors"`
	// Body       []byte
	Err error
}

// Prints the value of Err, which is an internal error message, usually including the body of the returned json from C7
func (e C7Error) Error() string {
	return e.Err.Error()
}

func (e *C7Error) ErrorFull() string {
	errorString := fmt.Sprintf("status code: %d, type: %s, message: %s, errors:", e.StatusCode, e.Type, e.Message)

	for i, err := range e.Errors {
		errorString += fmt.Sprintf(" (%d):", i+1)
		for key, value := range err {
			errorString += fmt.Sprintf("{ %s: %v }", key, value)
		}
	}
	return errorString
}

func (e *C7Error) ErrorReadable() string {
	errorString := fmt.Sprintf("status code: %d\ntype: %s\n message: %s\n", e.StatusCode, e.Type, e.Message)

	for i, err := range e.Errors {
		errorString += fmt.Sprintf("  Error %d:\n", i+1)
		for key, value := range err {
			errorString += fmt.Sprintf("    %s: %v\n", key, value)
		}
	}
	return errorString
}

func (e *C7Error) ErrorSimple() string {
	return fmt.Sprintf("status code: %d, type: %s, message: %s", e.StatusCode, e.Type, e.Message)
}

func (e *C7Error) UnmarshalJSON(data []byte) error {
	type Alias C7Error // Prevent recursion
	aux := &struct {
		StatusCode any `json:"statusCode"` // Accepts int or string
		*Alias
	}{
		Alias: (*Alias)(e),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Convert statusCode to string
	switch v := aux.StatusCode.(type) {
	case int:
		e.StatusCode = v
	case float64:
		e.StatusCode = int(v)
	case float32:
		e.StatusCode = int(v)
	case string:
		var err error
		if e.StatusCode, err = strconv.Atoi(v); err != nil {
			e.StatusCode = 0
		}
	default:
		e.StatusCode = 0
	}
	return nil
}
