package c7api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

const SLEEP_TIME = 500 * time.Millisecond

// Basic requests to C7 endpoint wrapped in retry logic with exponential backoff
func NewRequest(method string, url *string, reqBody *[]byte, tenant *string, c7AppAuthEncoded *string, retryCount int) (*[]byte, error) {
	//
	if url == nil || tenant == nil || c7AppAuthEncoded == nil {
		return nil, fmt.Errorf("error getting JSON from C7: nil value in arguments")
	}

	if reqBody == nil {
		fmt.Println("reqBody is nil")
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
		req, err := http.NewRequest(method, *url, bytes.NewBuffer(*reqBody))
		if err != nil {
			return nil, fmt.Errorf("error creating GET request for C7: %v", err)
		}

		req.Header.Set("tenant", *tenant)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("Authorization", *c7AppAuthEncoded)

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

	return &body, C7Error{response.StatusCode, fmt.Errorf(string(body))}

}

// Errors will return a custom error type called C7Error if there is an error directly from C7, calling err.Error() on this will return the error message from C7 and the status code.

// Takes in a full URL string and request JSON from C7 and return it as a byte array
//
// Attempts to get JSON from C7, if it fails, it will retry the request up to the number of attempts specified. Min 1, Max 10.
// Will wait 500ms between attempts.
func GetJsonFromC7(urlString *string, tenant *string, auth *string, attempts int) (*[]byte, error) {

	if urlString == nil || tenant == nil || auth == nil {
		return nil, fmt.Errorf("error getting JSON from C7: nil value in arguments")
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

		req.Header.Set("tenant", *tenant)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("Authorization", *auth)

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

		if response.StatusCode == 200 || response.StatusCode == 201 {
			return &body, nil
		} else {
			if response.StatusCode == http.StatusTooManyRequests {
				exponSleepTime := SLEEP_TIME * time.Duration(i)
				time.Sleep(exponSleepTime)
			} else {
				time.Sleep(SLEEP_TIME)
			}
		}
	}

	// Response is not 200, return error
	return &body, C7Error{response.StatusCode, fmt.Errorf(string(body))}

}

// Function to send shipment data back to C7
// Recieve XML from ShipStation, and post shipment data back to C7 as JSON
//
// URL Example: [Your Web Endpoint]?action=shipnotify&order_number=[Order Number]&carrier=[Carrier]&service=&tracking_number=[Tracking Number]
// URL END POINT]?action=shipnotify&order_number=ABC123&carrier=USPS&service=&tracking_number=9511343223432432432
func PostJsonToC7(urlString *string, tenant *string, reqBody *[]byte, auth *string, attempts int) (*[]byte, error) {

	if urlString == nil || tenant == nil || reqBody == nil || auth == nil {
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

		req.Header.Set("tenant", *tenant)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("Authorization", *auth) //AppAuthEncoded

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

		if response.StatusCode == 200 || response.StatusCode == 201 {
			return &body, nil
		} else {
			if response.StatusCode == http.StatusTooManyRequests {
				exponSleepTime := SLEEP_TIME * time.Duration(i)
				time.Sleep(exponSleepTime)
			} else {
				time.Sleep(SLEEP_TIME)
			}
		}
	}

	return &body, C7Error{response.StatusCode, fmt.Errorf(string(body))}
}

func PutJsonToC7(urlString *string, tenant *string, reqBody *[]byte, auth *string, attempts int) (*[]byte, error) {
	if urlString == nil || tenant == nil || reqBody == nil || auth == nil {
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
		req, err := http.NewRequest("PUT", *urlString, bytes.NewBuffer(*reqBody))
		if err != nil {
			return nil, fmt.Errorf("error creating PUT request to C7: %v", err)
		}

		req.Header.Set("tenant", *tenant)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("Authorization", *auth) //AppAuthEncoded

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

		if response.StatusCode == 200 || response.StatusCode == 201 {
			return &body, nil
		} else {
			//fmt.Println("Attempt: ", i+1, " of ", attempts, " failed. Status Code: ", response.StatusCode, " Error: ", string(body))
			if response.StatusCode == http.StatusTooManyRequests {
				exponSleepTime := SLEEP_TIME * time.Duration(i)
				time.Sleep(exponSleepTime)
			} else {
				time.Sleep(SLEEP_TIME)
			}
		}
	}

	return &body, C7Error{response.StatusCode, fmt.Errorf(string(body))}
}

func DeleteFromC7(urlString *string, tenant *string, auth *string, attempts int) (*[]byte, error) {
	if urlString == nil || tenant == nil || auth == nil {
		return nil, fmt.Errorf("nil value in arguments")
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
		req.Header.Set("tenant", *tenant)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("Authorization", *auth) //AppAuthEncoded

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

		if response.StatusCode == 200 || response.StatusCode == 201 { // C7 docs are lying, they return 200 on success along with the full order object.
			return &body, nil
		} else {
			if response.StatusCode == http.StatusTooManyRequests {
				exponSleepTime := SLEEP_TIME * time.Duration(i)
				time.Sleep(exponSleepTime)
			} else {
				time.Sleep(SLEEP_TIME)
			}
		}
	}

	// Response is not 200 or 201, return error
	return &body, C7Error{response.StatusCode, fmt.Errorf(string(body))}

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
	ordersBytes, err := GetJsonFromC7(&orderUrl, &tenant, &auth, attempts)
	if err != nil {
		return "", err
	}

	// Unmarshal the order
	var orders C7Orders
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

	return orders.Orders[0].Fulfillments[0].ID, nil

}

func DeleteC7Fulfillment(orderId string, fulfillmentId string, tenant string, auth string, attempts int) error {

	deleteUrl := "https://api.commerce7.com/v1/order/" + orderId + "/fulfillment/" + fulfillmentId
	// DELETE /order/{:id}/fulfillment/{:id}
	_, err := DeleteFromC7(&deleteUrl, &tenant, &auth, attempts)
	if err != nil {
		return err
	}

	//fmt.Println("Delete Fulfillment Response: ", string(*bytes))

	return nil

}

func IsCarrierSupported(carrier string) bool {
	switch carrier {
	case "UPS":
		return true
	case "FedEx":
		return true
	case "GSO":
		return true
	case "ATS Healthcare":
		return true
	case "Australia Post":
		return true
	default:
		return false
	}
}
