package c7api

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"
)

// Takes in a full URL string and request JSON from C7 and return it as a byte array
func GetJsonFromC7(urlString *string, tenant *string, auth *string) (*[]byte, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", *urlString, nil)
	if err != nil {
		return nil, fmt.Errorf("while creating request to C7, got: %v", err)
	}

	req.Header.Set("tenant", *tenant)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", *auth)

	// Make request to C7
	response, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("while making request to C7, got: %v", err)
	}

	defer response.Body.Close()

	if response.StatusCode == 401 {
		return nil, errors.New("while making request to C7, got status code: 401 unauthorized, please contact marsbytes support. support@marsbytesapps.com")
	}

	if response.StatusCode != 200 {
		attemptCount := 0
		for response.StatusCode != 200 && attemptCount < 3 {
			time.Sleep(time.Duration(rand.Intn(3)) * time.Second)
			//time.Sleep(10 * time.Second)
			response, err = client.Do(req)
			if err != nil {
				return nil, fmt.Errorf("while making request to C7, got: %v", err)
			}
			attemptCount++
		}
		if attemptCount >= 3 {
			return nil, fmt.Errorf("while making request to C7, got status code: %v, please contact marsbytes support. Marsbytes.dev/shipstationapi", response.StatusCode)
		}
	}

	// Read the body into variable
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("while reading response body from C7, got: %v", err)
	}

	return &body, nil

}

// Function to send shipment data back to C7
// Recieve XML from ShipStation, and post shipment data back to C7 as JSON
// URL Example: [Your Web Endpoint]?action=shipnotify&order_number=[Order Number]&carrier=[Carrier]&service=&tracking_number=[Tracking Number]
// Need to decode the URL [URL END POINT]?action=shipnotify&order_number=ABC123&carrier=USPS&service=&tracking_number=9511343223432432432
func PostJsonToC7(urlString *string, tenant *string, body *[]byte, auth *string) (*[]byte, int, error) {

	if urlString == nil || tenant == nil || body == nil || auth == nil {
		return nil, 0, fmt.Errorf("while posting JSON to C7, got: nil value in arguments")
	}

	// prepare request
	client := &http.Client{}

	req, err := http.NewRequest("POST", *urlString, bytes.NewBuffer(*body))
	if err != nil {
		return nil, 0, fmt.Errorf("while creating POST request to C7, got: %v", err)
	}

	// Set headers
	req.Header.Set("tenant", *tenant)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", *auth) //AppAuthEncoded

	//req.Body = io.NopCloser(bytes.NewReader(*body))

	// Make request to C7
	response, err := client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("while making POST request to C7, got: %v", err)
	}

	defer response.Body.Close()

	// if response.StatusCode == 401 {
	// 	return nil, response.StatusCode, errors.New("while making request to C7, got status code: 401 unauthorized, please contact marsbytes support. support@marsbytesapps.com")
	// }

	// if response.StatusCode != 200 {
	// 	attemptCount := 0
	// 	for response.StatusCode != 200 && attemptCount < 3 {
	// 		time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
	// 		//time.Sleep(10 * time.Second)
	// 		response, err = client.Do(req)
	// 		if err != nil {
	// 			return nil, 0, fmt.Errorf("while making request to C7, got: %v", err)
	// 		}
	// 		attemptCount++
	// 	}
	// 	if attemptCount >= 3 {
	// 		*body, _ = io.ReadAll(response.Body)
	// 		return body, response.StatusCode, fmt.Errorf("while making request to C7, got status code: %v, please contact marsbytes support. Marsbytes.dev/shipstationapi", response.StatusCode)
	// 	}
	// }

	// Read the body into variable
	*body, err = io.ReadAll(response.Body)
	if err != nil {
		return nil, 0, fmt.Errorf("while reading response body from C7, got: %v", err)
	}

	return body, response.StatusCode, nil
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
