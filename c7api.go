package c7api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const SLEEP_TIME = 500 * time.Millisecond

// Basic request. Will return the response or error if any.
func Request(method string, url string, reqBody *[]byte, tenant string, c7AppAuthEncoded string, errorOnNotOK bool) (*http.Response, error) {
	//
	if url == "" || tenant == "" || c7AppAuthEncoded == "" {
		return nil, fmt.Errorf("error getting JSON from C7: nil or blank value in arguments")
	}

	if reqBody == nil {
		reqBody = &[]byte{}
	}

	client := &http.Client{}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(*reqBody))
	if err != nil {
		return nil, fmt.Errorf("error creating GET request for C7: %v", err)
	}

	req.Header.Set("tenant", tenant)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", c7AppAuthEncoded)

	response, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making GET request to C7: %v", err)
	}

	if errorOnNotOK && !ResponseIsOK(response.StatusCode) {
		return response, errors.New("reponse status not within 200-299: " + response.Status)
	}

	return response, nil

}

// Basic requests to C7 endpoint wrapped in retry logic with exponential backoff for TooManyRequest responses.
//
// Reads out the response body and returns the bytes.
//
// Min Retry Count: 0 | Max Retry Count: 10
func RequestWithRetryAndRead(method string, url string, reqBody *[]byte, tenant string, c7AppAuthEncoded string, retryCount int, rl genericRateLimiter) (*[]byte, error) {
	//
	if url == "" || tenant == "" || c7AppAuthEncoded == "" {
		return nil, fmt.Errorf("error getting JSON from C7: nil or blank value in arguments")
	}

	if reqBody == nil {
		reqBody = &[]byte{}
	}

	minRetryCount := 0
	maxRetryCount := 10

	if retryCount < minRetryCount {
		retryCount = minRetryCount
	} else if retryCount > maxRetryCount {
		retryCount = maxRetryCount
	}

	client := &http.Client{}
	response := &http.Response{StatusCode: 0}
	body := []byte{}

	for i := 0; i <= retryCount; i++ {
		if rl != nil && !reflect.ValueOf(rl).IsNil() {
			rl.Wait()
		}
		req, err := http.NewRequest(method, url, bytes.NewBuffer(*reqBody))
		if err != nil {
			return nil, fmt.Errorf("error creating GET request for C7: %v", err)
		}

		req.Header.Set("tenant", tenant)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("Authorization", c7AppAuthEncoded)

		response, err = client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("error making GET request to C7: %v", err)
		}

		body, err = io.ReadAll(response.Body)
		response.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("error reading response body from C7: %v", err)
		}

		// 200-299 is success, return body and nil error
		if ResponseIsOK(response.StatusCode) {
			return &body, nil
		} else {
			// Exponential backoff based on retry count
			if response.StatusCode == http.StatusTooManyRequests {
				exponSleepTime := SLEEP_TIME * time.Duration(i)
				time.Sleep(exponSleepTime)
			} else {
				time.Sleep(SLEEP_TIME)
			}
		}
	}

	// Read the C7 Error if present
	// Always return as C7Error after this point, since this means C7 sent an error message.
	// If we have trouble reading it for some reason, handle that here.
	c7Error := C7Error{}
	err := json.Unmarshal(body, &c7Error)
	if err != nil {
		c7Error.StatusCode = response.StatusCode
		c7Error.Err = errors.New("error unmarshalling Commerce7 Error Message: " + err.Error() + "json: " + string(body))
		return &body, c7Error
	}

	c7Error.Err = errors.New(string(body))
	return &body, c7Error

}

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

// Returns the fulfillment ids if there is any fulfillments on a C7 order.
//
// Usually this will return just one, but can return multiple if there are partial fulfillments or errors with C7.
func GetFulfillmentIds(OrderNumber int, tenant string, auth string, attempts int, rl genericRateLimiter) ([]string, error) {

	orderUrl := Endpoints.Order + "?q=" + strconv.Itoa(OrderNumber)
	fulfillments := []string{}
	// Get the order from C7
	ordersBytes, err := RequestWithRetryAndRead("GET", orderUrl, nil, tenant, auth, attempts, rl)
	if err != nil {
		return fulfillments, err
	}

	// Unmarshal the order
	var orders C7OrdersFulfillmentsOnly
	err = json.Unmarshal(*ordersBytes, &orders)
	if err != nil {
		return fulfillments, err
	}

	// Get the fulfillment ID
	if len(orders.Orders) == 0 {
		return fulfillments, errors.New("no orders found")
	}
	if len(orders.Orders[0].Fulfillments) == 0 {
		return fulfillments, errors.New("no fulfillments found")
	}
	for _, order := range orders.Orders {
		if order.OrderNumber == OrderNumber {
			// fulfillments are always an array of len 1 in C7 unless there were multiple fulfillments that are still valid.
			for _, fulfillment := range order.Fulfillments {
				fulfillments = append(fulfillments, fulfillment.ID)
			}
			return fulfillments, nil
		}
	}
	return fulfillments, errors.New("no matching order found")

}

func GetFulfillments(OrderNumber int, tenant string, auth string, attempts int, rl genericRateLimiter) (*[]C7OrderFulfillment, error) {

	orderUrl := Endpoints.Order + "?q=" + strconv.Itoa(OrderNumber)
	// Get the order from C7
	ordersBytes, err := RequestWithRetryAndRead("GET", orderUrl, nil, tenant, auth, attempts, rl)
	if err != nil {
		return nil, err
	}

	// Unmarshal the order
	var orders C7Orders
	err = json.Unmarshal(*ordersBytes, &orders)
	if err != nil {
		return nil, err
	}

	// Get the fulfillment ID
	if len(orders.Orders) == 0 {
		return nil, errors.New("no orders found")
	}
	if len(orders.Orders[0].Fulfillments) == 0 {
		return nil, errors.New("no fulfillments found")
	}
	for _, order := range orders.Orders {
		if order.OrderNumber == OrderNumber {
			return &order.Fulfillments, nil
		}
	}
	return nil, errors.New("no matching order found")

}

func DeleteFulfillment(orderId string, fulfillmentId string, tenant string, auth string, attempts int, rl genericRateLimiter) (*[]byte, error) {

	deleteUrl := Endpoints.Order + "/" + orderId + "/fulfillment/" + fulfillmentId
	// DELETE /order/{:id}/fulfillment/{:id}
	return RequestWithRetryAndRead("DELETE", deleteUrl, nil, tenant, auth, attempts, rl)

}

func MarkNoFulfillmentRequired(orderId string, shipTime time.Time, tenant string, auth string, attempts int, rl genericRateLimiter) error {
	// POST // https://api.commerce7.com/v1/order/b9f10447-4285-4dc2-add2-b38798dba8f9/fulfillment

	// Create new Fulfillment from struct
	var fulfillment FulfillmentAllItems
	fulfillment.SendTransactionEmail = false
	fulfillment.Type = "No Fulfillment Required"
	fulfillment.FulfillmentDate = shipTime

	url := "https://api.commerce7.com/v1/order/" + orderId + "/fulfillment/all"

	// Convert Fulfillment struct to JSON
	fulfillmentJSON, err := json.Marshal(fulfillment)
	if err != nil {
		return errors.New("error marshaling NFR fulfillment into JSON: " + err.Error())
	}

	// Post the fulfillment to C7
	_, err = RequestWithRetryAndRead("POST", url, &fulfillmentJSON, tenant, auth, attempts, rl)
	if err != nil {
		return errors.New("error posting NFR fulfillment to C7: " + err.Error())
	}

	return nil
}

func IsCarrierSupported(carrier string) bool {
	switch strings.ToUpper(carrier) {
	case "UPS":
		return true
	case "FEDEX":
		return true
	case "GSO":
		return true
	case "ATS HEALTHCARE":
		return true
	case "AUSTRALIA POST":
		return true
	default:
		return false
	}
}

func ResponseIsOK(statusCode int) bool {
	return statusCode >= 200 && statusCode <= 299
}
