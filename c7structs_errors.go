package c7api

import (
	"fmt"
)

type C7Error struct {
	StatusCode string           `json:"statusCode"`
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
	errorString := fmt.Sprintf("status code: %s, type: %s, message: %s, errors:", e.StatusCode, e.Type, e.Message)

	for i, err := range e.Errors {
		errorString += fmt.Sprintf(" (%d):", i+1)
		for key, value := range err {
			errorString += fmt.Sprintf("{ %s: %v }", key, value)
		}
	}
	return errorString
}

func (e *C7Error) ErrorReadable() string {
	errorString := fmt.Sprintf("status code: %s\ntype: %s\n message: %s\n", e.StatusCode, e.Type, e.Message)

	for i, err := range e.Errors {
		errorString += fmt.Sprintf("  Error %d:\n", i+1)
		for key, value := range err {
			errorString += fmt.Sprintf("    %s: %v\n", key, value)
		}
	}
	return errorString
}

func (e *C7Error) ErrorSimple() string {
	return fmt.Sprintf("status code: %s, type: %s, message: %s", e.StatusCode, e.Type, e.Message)
}
