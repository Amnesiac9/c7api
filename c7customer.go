package c7api

import (
	"encoding/json"
	"fmt"
	"strings"
)

type HasEmails interface {
	GetEmails() []Email
}

func GetCustomerByEmail[T HasEmails](email string, tenant string, c7AppAuthEncoded string, retryCount int, rl genericRateLimiter) (*T, error) {

	type Customers struct {
		Customers []T `json:"customers"`
		// Total     int `json:"total"`
	}

	reqUrl := Endpoints.Customer

	email = strings.ToLower(email)

	// Query for the email
	quieries := map[string]string{
		"q": email,
	}

	resp, err := RequestWithRetryAndRead("GET", reqUrl, quieries, nil, tenant, c7AppAuthEncoded, retryCount, rl)
	if err != nil {
		return nil, fmt.Errorf("failed to get customer from tenant: %w", err)
	}

	var c7Customers Customers
	if err := json.Unmarshal(*resp, &c7Customers); err != nil {
		return nil, fmt.Errorf("failed to unmarshal customer: %w", err)
	}

	for _, customer := range c7Customers.Customers {
		for _, e := range customer.GetEmails() {
			if strings.ToLower(e.Email) == email {
				return &customer, nil
			}
		}
	}

	return nil, fmt.Errorf("no customer found")
}

func GetCustomerById[T HasEmails](customerId string, tenant string, c7AppAuthEncoded string, retryCount int, rl genericRateLimiter) (*T, error) {

	reqUrl := Endpoints.Customer + "/" + customerId

	resp, err := RequestWithRetryAndRead("GET", reqUrl, nil, nil, tenant, c7AppAuthEncoded, retryCount, rl)
	if err != nil {
		return nil, fmt.Errorf("failed to get customer from tenant: %w", err)
	}

	var c7Customer T
	if err := json.Unmarshal(*resp, &c7Customer); err != nil {
		return nil, fmt.Errorf("failed to unmarshal customer: %w", err)
	}

	return &c7Customer, nil
}

func GetCustomersWithCursor[T HasEmails](tenant string, queries map[string]string, c7AppAuthEncoded string, retryCount int, rl genericRateLimiter) (*[]T, error) {

	type CustomersCursor struct {
		Customers []T
		Cursor    string `json:"cursor"`
	}

	var c7Customers []T

	cursor := "start"
	reqUrl := Endpoints.Customer
	if queries != nil {
		queries["cursor"] = cursor
	}

	// Loop until cursor ends
	for {
		if cursor == "" {
			break
		}
		resp, err := RequestWithRetryAndRead("GET", reqUrl, queries, nil, tenant, c7AppAuthEncoded, retryCount, rl)
		if err != nil {
			return nil, fmt.Errorf("failed to get customer from tenant: %w", err)
		}

		var c7CustomerBatch CustomersCursor
		if err := json.Unmarshal(*resp, &c7CustomerBatch); err != nil {
			return nil, fmt.Errorf("failed to unmarshal customer: %w", err)
		}

		c7Customers = append(c7Customers, c7CustomerBatch.Customers...)
		cursor = c7CustomerBatch.Cursor
	}

	return &c7Customers, nil
}
