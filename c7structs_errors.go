package c7api

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// C7Error represents a structured error response from the Commerce7 API.
// It implements the error interface and contains detailed information about API failures.
//
// Basic usage:
//
//	body, err := c7api.RequestWithRetryAndReadV2(...)
//	if err != nil {
//	    log.Fatal(err) // Prints the error message
//	}
//
// To access structured error details, use errors.As:
//
//	var c7Err *c7api.C7Error
//	if errors.As(err, &c7Err) {
//	    fmt.Printf("Status: %d\n", c7Err.StatusCode)
//	    fmt.Printf("Type: %s\n", c7Err.Type)
//	    fmt.Printf("Message: %s\n", c7Err.Message)
//	    for _, errDetail := range c7Err.Errors {
//	        fmt.Printf("Error details: %v\n", errDetail)
//	    }
//	}
type C7Error struct {
	StatusCode int              `json:"statusCode"` // HTTP status code from the API response
	Type       string           `json:"type"`       // Error type classification from Commerce7
	Message    string           `json:"message"`    // Human-readable error message from Commerce7
	Errors     []map[string]any `json:"errors"`     // Additional error details and validation errors
	Err        error            // Internal error containing the full response body or parsing errors
}

// Error implements the error interface.
// Returns the internal error message, which typically includes the full response body from Commerce7.
func (e C7Error) Error() string {
	return e.Err.Error()
}

// ErrorFull returns a comprehensive single-line error string containing all error details.
// Includes status code, type, message, and all nested errors with their fields.
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

// ErrorReadable returns a multi-line formatted error string optimized for human readability.
// Each error detail is printed on separate lines with proper indentation.
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

// ErrorSimple returns a concise single-line summary of the error.
// Includes only the status code, type, and message without nested error details.
func (e *C7Error) ErrorSimple() string {
	return fmt.Sprintf("status code: %d, type: %s, message: %s", e.StatusCode, e.Type, e.Message)
}

// UnmarshalJSON implements custom JSON unmarshaling for C7Error.
// Handles statusCode as either an integer or string type, converting to int.
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

	// Convert statusCode to int
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
