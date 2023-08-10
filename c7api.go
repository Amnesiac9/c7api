package c7api

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Errors will return a custom error type called C7Error if there is an error directly from C7, calling err.Error() on this will return the error message from C7 and the status code.

// Takes in a full URL string and request JSON from C7 and return it as a byte array
func GetJsonFromC7(urlString *string, tenant *string, auth *string) (*[]byte, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", *urlString, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating GET request for C7: %v", err)
	}

	req.Header.Set("tenant", *tenant)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", *auth)

	// Make request to C7
	response, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making GET request to C7: %v", err)
	}

	defer response.Body.Close()

	// Read the body into variable
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body from C7: %v", err)
	}

	if response.StatusCode != 200 {
		return &body, C7Error{response.StatusCode, fmt.Errorf(string(body))}
	}

	// Response is 200, return body
	return &body, nil

}

// Function to send shipment data back to C7
// Recieve XML from ShipStation, and post shipment data back to C7 as JSON
// URL Example: [Your Web Endpoint]?action=shipnotify&order_number=[Order Number]&carrier=[Carrier]&service=&tracking_number=[Tracking Number]
// Need to decode the URL [URL END POINT]?action=shipnotify&order_number=ABC123&carrier=USPS&service=&tracking_number=9511343223432432432
func PostJsonToC7(urlString *string, tenant *string, body *[]byte, auth *string) (*[]byte, error) {

	if urlString == nil || tenant == nil || body == nil || auth == nil {
		return nil, fmt.Errorf("while posting JSON to C7, got: nil value in arguments")
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
		return nil, C7Error{http.StatusInternalServerError, fmt.Errorf("nil value in arguments")}
	}

	// prepare request
	client := &http.Client{}

	req, err := http.NewRequest("DELETE", *urlString, nil)
	if err != nil {
		return nil, C7Error{http.StatusInternalServerError, fmt.Errorf("while creating DELETE request to C7, got: %v", err)}
	}

	// Set headers
	req.Header.Set("tenant", *tenant)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", *auth) //AppAuthEncoded

	// Make request to C7
	response, err := client.Do(req)
	if err != nil {
		return nil, C7Error{http.StatusInternalServerError, fmt.Errorf("while making DELETE request to C7, got: %v", err)}
	}

	defer response.Body.Close()

	// Read the body into variable
	c7Body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, C7Error{http.StatusInternalServerError, fmt.Errorf("while reading response body from C7, got: %v", err)}
	}

	if response.StatusCode != 204 {
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
