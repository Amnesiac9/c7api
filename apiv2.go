package c7api

import (
	"encoding/json"
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
	return requestWithRetryAndRead(method, url, queries, reqBody, tenant, c7AppAuthEncoded, retryCount, rl, v2Headers(tenant))
}

// v2Headers are the extra headers the experimental v2 API requires on top of
// the standard v1 set.
func v2Headers(tenant string) map[string]string {
	return map[string]string{
		"tenantid":     tenant,
		"experimental": "Do not use if you are not Commerce7.  API likely to change",
	}
}
