package c7api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
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

// Basic requests to C7 endpoint wrapped in retry logic with exponential backoff.
// Reads out the response body and returns the bytes.
func RequestWithRetryAndRead(method string, url string, reqBody *[]byte, tenant string, c7AppAuthEncoded string, retryCount int) (*[]byte, error) {
	//
	if url == "" || tenant == "" || c7AppAuthEncoded == "" {
		return nil, fmt.Errorf("error getting JSON from C7: nil or blank value in arguments")
	}

	if reqBody == nil {
		reqBody = &[]byte{}
	}

	if retryCount < 0 {
		retryCount = 0
	} else if retryCount > 10 {
		retryCount = 10
	}

	client := &http.Client{}
	response := &http.Response{StatusCode: 0}
	body := []byte{}

	for i := 0; i <= retryCount; i++ {
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
		if response.StatusCode >= 200 && response.StatusCode <= 299 {
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

	return &body, C7Error{response.StatusCode, errors.New(string(body))}

}

// Errors will return a custom error type called C7Error if there is an error directly from C7, calling err.Error() on this will return the error message from C7 and the status code.
//
// # Takes in a full URL string and request JSON from C7 and return it as a byte array
//
// Attempts to get JSON from C7, if it fails, it will retry the request up to the number of attempts specified. Min 1, Max 10.
// Will wait 500ms between attempts.
func GetReq(urlString *string, tenant string, auth string, attempts int) (*[]byte, error) {

	if urlString == nil || tenant == "" || auth == "" {
		return nil, fmt.Errorf("error getting JSON from C7: nil or blank value in arguments")
	}

	if attempts < 1 {
		attempts = 1
	} else if attempts > 10 {
		attempts = 10
	}

	// Make request to C7
	client := &http.Client{}
	response := &http.Response{StatusCode: 0}
	body := []byte{}
	var i int

	for i = 0; i < attempts; i++ {
		req, err := http.NewRequest("GET", *urlString, nil)
		if err != nil {
			return nil, fmt.Errorf("error creating GET request for C7: %v", err)
		}

		req.Header.Set("tenant", tenant)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("Authorization", auth)

		response, err = client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("error making GET request to C7: %v", err)
		}

		// Read the body into variable
		body, err = io.ReadAll(response.Body)
		response.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("error reading response body from C7: %v", err)
		}

		// 200-299 is success, return body and nil error
		if response.StatusCode >= 200 && response.StatusCode <= 299 {
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

	// Response is not 200, return error
	return &body, C7Error{response.StatusCode, errors.New(string(body))}

}

// Function to send shipment data back to C7
// Recieve XML from ShipStation, and post shipment data back to C7 as JSON
//
// URL Example: [Your Web Endpoint]?action=shipnotify&order_number=[Order Number]&carrier=[Carrier]&service=&tracking_number=[Tracking Number]
// URL END POINT]?action=shipnotify&order_number=ABC123&carrier=USPS&service=&tracking_number=9511343223432432432
func PostReq(urlString *string, reqBody *[]byte, tenant string, auth string, attempts int) (*[]byte, error) {

	if urlString == nil || tenant == "" || reqBody == nil || auth == "" {
		return nil, fmt.Errorf("error posting JSON to C7: nil value in arguments")
	}

	if attempts < 1 {
		attempts = 1
	} else if attempts > 10 {
		attempts = 10
	}

	// Make request to C7
	client := &http.Client{}
	response := &http.Response{StatusCode: 0}
	body := []byte{}
	var i int

	for i = 0; i < attempts; i++ {
		// Cannot reuse the same request, need to create a new one each time. (Not sure why, but causes cloudflare issues on C7's end)
		req, err := http.NewRequest("POST", *urlString, bytes.NewBuffer(*reqBody))
		if err != nil {
			return nil, fmt.Errorf("error creating POST request to C7: %v", err)
		}

		req.Header.Set("tenant", tenant)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("Authorization", auth) //AppAuthEncoded

		response, err = client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("error making POST request to C7: %v", err)
		}

		// Read the body into variable
		body, err = io.ReadAll(response.Body)
		response.Body.Close() // Remove defer when using for loop, close the body each time it is read.
		if err != nil {
			return nil, fmt.Errorf("error reading response body from C7: %v", err)
		}

		// 200-299 is success, return body and nil error
		if response.StatusCode >= 200 && response.StatusCode <= 299 {
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

	return &body, C7Error{response.StatusCode, errors.New(string(body))}
}

func PutReq(urlString *string, reqBody *[]byte, tenant string, auth string, attempts int) (*[]byte, error) {
	if urlString == nil || reqBody == nil || auth == "" {
		return nil, fmt.Errorf("error posting JSON to C7: nil or blank value in arguments")
	}

	if attempts < 1 {
		attempts = 1
	} else if attempts > 10 {
		attempts = 10
	}

	// Make request to C7
	client := &http.Client{}
	response := &http.Response{StatusCode: 0}
	body := []byte{}
	var i int

	for i = 0; i < attempts; i++ {
		// Cannot reuse the same request, need to create a new one each time. (Not sure why, but causes cloudflare issues on C7's end)
		req, err := http.NewRequest("PUT", *urlString, bytes.NewBuffer(*reqBody))
		if err != nil {
			return nil, fmt.Errorf("error creating PUT request to C7: %v", err)
		}

		req.Header.Set("tenant", tenant)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("Authorization", auth) //AppAuthEncoded

		response, err = client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("error making PUT request to C7: %v", err)
		}

		// Read the body into variable
		body, err = io.ReadAll(response.Body)
		response.Body.Close() // Remove defer when using for loop, close the body each time it is read.
		if err != nil {
			return nil, fmt.Errorf("error reading response body from C7: %v", err)
		}

		// 200-299 is success, return body and nil error
		if response.StatusCode >= 200 && response.StatusCode <= 299 {
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

	return &body, C7Error{response.StatusCode, errors.New(string(body))}
}

func DeleteReq(urlString *string, tenant string, auth string, attempts int) (*[]byte, error) {
	if urlString == nil || tenant == "" || auth == "" {
		return nil, fmt.Errorf("nil or blank value in arguments")
	}

	if attempts < 1 {
		attempts = 1
	} else if attempts > 10 {
		attempts = 10
	}

	// Make request to C7
	client := &http.Client{}
	response := &http.Response{StatusCode: 0}
	body := []byte{}
	var i int

	for i = 0; i < attempts; i++ {

		req, err := http.NewRequest("DELETE", *urlString, nil)
		if err != nil {
			return nil, fmt.Errorf("while creating DELETE request to C7, got: %v", err)
		}

		// Set headers
		req.Header.Set("tenant", tenant)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("Authorization", auth) //AppAuthEncoded

		response, err = client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("while making DELETE request to C7, got: %v", err)
		}

		// Read the body into variable
		body, err = io.ReadAll(response.Body)
		response.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("while reading response body from C7, got: %v", err)
		}

		// C7 docs are lying, they return 200 on success along with the full order object.
		// 200-299 is success, return body and nil error
		if response.StatusCode >= 200 && response.StatusCode <= 299 {
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

	// Response is not 200 or 201, return error
	return &body, C7Error{response.StatusCode, errors.New(string(body))}

}

func FormatDatesForC7(date string) (string, error) {
	if date == "" {
		return date, errors.New("date is empty")
	}

	layout := "01/02/2006 15:04"

	dateFormatted, err := time.Parse(layout, date)
	if err != nil {
		return date, err
	}

	return dateFormatted.Format("2006-01-02T15:04:05.000Z"), err
}

func GetFulfillmentId(OrderNumber int, tenant string, auth string, attempts int) (string, error) {

	orderUrl := "https://api.commerce7.com/v1/order?q=" + strconv.Itoa(OrderNumber)
	// Get the order from C7
	ordersBytes, err := GetReq(&orderUrl, tenant, auth, attempts)
	if err != nil {
		return "", err
	}

	// Unmarshal the order
	var orders C7OrdersFulfillmentsOnly
	err = json.Unmarshal(*ordersBytes, &orders)
	if err != nil {
		return "", err
	}

	// Get the fulfillment ID
	if len(orders.Orders) == 0 {
		return "", errors.New("no orders found")
	}
	if len(orders.Orders[0].Fulfillments) == 0 {
		return "", errors.New("no fulfillments found")
	}
	for _, order := range orders.Orders {
		if order.OrderNumber == OrderNumber {
			// fulfillments are always an array of len 1 in C7.
			return order.Fulfillments[0].ID, nil
		}
	}
	return "", errors.New("no matching order found")

}

func GetFulfillments(OrderNumber int, tenant string, auth string, attempts int) (*[]C7OrderFulfillment, error) {

	orderUrl := "https://api.commerce7.com/v1/order?q=" + strconv.Itoa(OrderNumber)
	// Get the order from C7
	ordersBytes, err := GetReq(&orderUrl, tenant, auth, attempts)
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

func DeleteFulfillment(orderId string, fulfillmentId string, tenant string, auth string, attempts int) (*[]byte, error) {

	deleteUrl := "https://api.commerce7.com/v1/order/" + orderId + "/fulfillment/" + fulfillmentId
	// DELETE /order/{:id}/fulfillment/{:id}
	return DeleteReq(&deleteUrl, tenant, auth, attempts)

}

func MarkNoFulfillmentRequired(orderId string, shipTime time.Time, tenant string, auth string, attempts int) error {
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
	_, err = RequestWithRetryAndRead("POST", url, &fulfillmentJSON, tenant, auth, attempts)
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
