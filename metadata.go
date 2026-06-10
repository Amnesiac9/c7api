package c7api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Customer metadata

// CustomerMetaDataPost is the payload for updating a customer's metadata (custom fields).
//
// MetaData values can be many different types in the system (string, bool, etc.),
// so they are always keyed by string with an arbitrary value.
//
// {"metaData":{"test3":"blah1"}}
//
// #NOTE
// Does NOT require the full metadata map. Single updates are fine, the rest of the values will be untouched.
type CustomerMetaDataPost struct {
	MetaData map[string]any `json:"metaData"`
}

// PutCustomerMetaData updates the metadata (custom fields) for a single customer
// via PUT on the customer/{id} endpoint.
func PutCustomerMetaData(metaData map[string]any, customerId, tenant, c7AppAuthEncoded string, retryCount int, rl genericRateLimiter) (*Customer, error) {

	payload := CustomerMetaDataPost{MetaData: metaData}

	objectBytes, err := json.Marshal(&payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal customer metadata payload: %w", err)
	}

	reqUrl := Endpoints.Customer + "/" + customerId

	resp, err := RequestWithRetryAndRead(http.MethodPut, reqUrl, nil, &objectBytes, tenant, c7AppAuthEncoded, retryCount, rl)
	if err != nil {
		return nil, fmt.Errorf("failed to post customer metadata: %w", err)
	}

	var c7Customer Customer
	if err := json.Unmarshal(*resp, &c7Customer); err != nil {
		return nil, fmt.Errorf("failed to unmarshal customer after attempted metadata post: %w", err)
	}

	return &c7Customer, nil
}
