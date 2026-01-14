package c7api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"time"
)

// * V2 API - Currently in Experimental mode per Commerce7. Routes and headers required subject to change. *//

func RequestV2[T any](method, url string, queries map[string]string, reqBody *[]byte, tenant string, c7AppAuthEncoded string, retryCount int, rl genericRateLimiter) (*T, error) {
	data, err := RequestWithRetryAndReadV2(method, url, queries, reqBody, tenant, c7AppAuthEncoded, retryCount, rl)
	if err != nil {
		return nil, err
	}
	var v T
	err = json.Unmarshal(*data, &v)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

// Basic requests to C7 endpoint wrapped in retry logic with exponential backoff for TooManyRequest responses.
//
// Reads out the response body and returns the bytes.
//
// Min Retry Count: 0 | Max Retry Count: 10
func RequestWithRetryAndReadV2(method string, url string, queries map[string]string, reqBody *[]byte, tenant string, c7AppAuthEncoded string, retryCount int, rl genericRateLimiter) (*[]byte, error) {
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

		if queries != nil {
			query := req.URL.Query()
			for k, v := range queries {
				query.Add(k, v)
			}
			req.URL.RawQuery = query.Encode()
		}

		req.Header.Set("tenant", tenant)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("Authorization", c7AppAuthEncoded)

		//v2
		req.Header.Set("tenantid", tenant)
		req.Header.Set("experimental", "Do not use if you are not Commerce7.  API likely to change")

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
	err := c7Error.UnmarshalJSON(body)
	if err != nil {
		c7Error.StatusCode = response.StatusCode
		c7Error.Err = errors.New("error unmarshalling Commerce7 Error Message: " + err.Error() + "json: " + string(body))
		return &body, &c7Error
	}

	c7Error.Err = errors.New(string(body))
	return &body, &c7Error
}
