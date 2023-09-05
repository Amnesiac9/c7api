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

// type attemptCounter struct {
// 	attempts int
// }

// func (r *attemptCounter) add1() {
// 	r.attempts++
// }

// func (r *attemptCounter) reset() {
// 	r.attempts = 0
// }

// var attemptCount attemptCounter

// Errors will return a custom error type called C7Error if there is an error directly from C7, calling err.Error() on this will return the error message from C7 and the status code.

// Takes in a full URL string and request JSON from C7 and return it as a byte array
//
// Attempts to get JSON from C7, if it fails, it will retry the request up to the number of attempts specified. Min 1, Max 10.
// Will wait 500ms between attempts.
func GetJsonFromC7(urlString *string, tenant *string, auth *string, attempts int) (*[]byte, error) {

	const SLEEP_TIME = 500 * time.Millisecond

	if urlString == nil || tenant == nil || auth == nil {
		return nil, fmt.Errorf("error getting JSON from C7: nil value in arguments")
	}

	if attempts < 1 {
		attempts = 1
	} else if attempts > 10 {
		attempts = 10
	}

	req, err := http.NewRequest("GET", *urlString, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating GET request for C7: %v", err)
	}

	req.Header.Set("tenant", *tenant)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", *auth)

	// Make request to C7
	client := &http.Client{}
	response := &http.Response{StatusCode: 0}
	body := []byte{}
	var i int

	for i = 0; i < attempts; i++ {

		response, err = client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("error making GET request to C7: %v", err)
		}

		defer response.Body.Close()

		// Read the body into variable
		body, err = io.ReadAll(response.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading response body from C7: %v", err)
		}

		if response.StatusCode == 200 {
			return &body, nil
		} else {
			//fmt.Println("Attempt: ", i+1, " of ", attempts, " failed. Status Code: ", response.StatusCode, " Error: ", string(body))
			time.Sleep(SLEEP_TIME)
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
func PostJsonToC7(urlString *string, tenant *string, body *[]byte, auth *string) (*[]byte, error) {

	if urlString == nil || tenant == nil || body == nil || auth == nil {
		return nil, fmt.Errorf("error posting JSON to C7: nil value in arguments")
	}

	// prepare request
	client := &http.Client{}

	req, err := http.NewRequest("POST", *urlString, bytes.NewBuffer(*body))
	if err != nil {
		return nil, fmt.Errorf("error creating POST request to C7: %v", err)
	}

	// Set headers
	req.Header.Set("tenant", *tenant)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", *auth) //AppAuthEncoded

	// Make request to C7
	response, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making POST request to C7: %v", err)
	}

	defer response.Body.Close()

	// Read the body into variable
	c7Body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body from C7: %v", err)
	}

	if response.StatusCode != 200 {
		return &c7Body, C7Error{response.StatusCode, fmt.Errorf(string(c7Body))}
	}

	return &c7Body, nil
}

func DeleteFromC7(urlString *string, tenant *string, auth *string) (*[]byte, error) {
	if urlString == nil || tenant == nil || auth == nil {
		return nil, fmt.Errorf("nil value in arguments")
	}

	req, err := http.NewRequest("DELETE", *urlString, nil)
	if err != nil {
		return nil, fmt.Errorf("while creating DELETE request to C7, got: %v", err)
	}

	// Set headers
	req.Header.Set("tenant", *tenant)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", *auth) //AppAuthEncoded

	// Make request to C7
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("while making DELETE request to C7, got: %v", err)
	}

	defer response.Body.Close()

	// Read the body into variable
	c7Body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("while reading response body from C7, got: %v", err)
	}

	if response.StatusCode != 200 { // C7 docs are lying, they return 200 on success along with the full order object.
		return &c7Body, C7Error{response.StatusCode, fmt.Errorf(string(c7Body))}
	}

	return &c7Body, nil

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

func GetFulfillmentId(OrderNumber int, tenant string, auth string) (string, error) {

	orderUrl := "https://api.commerce7.com/v1/order?q=" + strconv.Itoa(OrderNumber)
	// Get the order from C7
	ordersBytes, err := GetJsonFromC7(&orderUrl, &tenant, &auth, 1)
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

func DeleteC7Fulfillment(orderId string, fulfillmentId string, tenant string, auth string) error {

	deleteUrl := "https://api.commerce7.com/v1/order/" + orderId + "/fulfillment/" + fulfillmentId
	// DELETE /order/{:id}/fulfillment/{:id}
	_, err := DeleteFromC7(&deleteUrl, &tenant, &auth)
	if err != nil {
		return err
	}

	//fmt.Println("Delete Fulfillment Response: ", string(*bytes))

	return nil

}
