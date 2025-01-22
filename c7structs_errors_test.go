package c7api

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
)

func Test_PrintC7Error(t *testing.T) {
	jsonData := `{
		"statusCode": 422,
		"type": "validationError",
		"message": "One or more elements is missing or invalid",
		"errors": [
			{
				"field": "csdom",
				"message": "invalid additional property"
			},
			{
				"field": "22d2d",
				"message": "invalid additional property"
			}
		]
	}`

	c7Error := C7Error{}
	if err := json.Unmarshal([]byte(jsonData), &c7Error); err != nil {
		t.Error("Error parsing JSON:", err)
		return
	}

	c7Error.Err = errors.New(jsonData)
	fmt.Println(c7Error.Error())

}

func Test_PrintC7ErrorFull(t *testing.T) {
	jsonData := `{
		"statusCode": 422,
		"type": "validationError",
		"message": "One or more elements is missing or invalid",
		"errors": [
			{
				"field": "csdom",
				"message": "invalid additional property"
			},
			{
				"field": "22d2d",
				"message": "invalid additional property"
			}
		]
	}`

	c7Error := C7Error{}
	if err := json.Unmarshal([]byte(jsonData), &c7Error); err != nil {
		t.Error("Error parsing JSON:", err)
		return
	}

	c7Error.Err = errors.New(jsonData)
	fmt.Println(c7Error.ErrorFull())

}

func Test_PrintC7ErrorReadable(t *testing.T) {
	jsonData := `{
		"statusCode": 422,
		"type": "validationError",
		"message": "One or more elements is missing or invalid",
		"errors": [
			{
				"field": "csdom",
				"message": "invalid additional property"
			},
			{
				"field": "22d2d",
				"message": "invalid additional property"
			}
		]
	}`

	c7Error := C7Error{}
	if err := json.Unmarshal([]byte(jsonData), &c7Error); err != nil {
		t.Error("Error parsing JSON:", err)
		return
	}

	c7Error.Err = errors.New(jsonData)
	fmt.Println(c7Error.ErrorReadable())

}

func Test_PrintC7ErrorSimple(t *testing.T) {
	jsonData := `{
		"statusCode": 422,
		"type": "validationError",
		"message": "One or more elements is missing or invalid",
		"errors": [
			{
				"field": "csdom",
				"message": "invalid additional property"
			},
			{
				"field": "22d2d",
				"message": "invalid additional property"
			}
		]
	}`

	c7Error := C7Error{}
	if err := json.Unmarshal([]byte(jsonData), &c7Error); err != nil {
		t.Error("Error parsing JSON:", err)
		return
	}

	c7Error.Err = errors.New(jsonData)
	fmt.Println(c7Error.ErrorSimple())

}
