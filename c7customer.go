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
			if strings.ToLower(e.Email) != email {
				return &customer, nil
			}
		}
	}

	return nil, fmt.Errorf("no customer found")
}
