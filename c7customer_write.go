package c7api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// CustomerWrite is a customer to create or update. Empty fields aren't sent.
// A PUT replaces lists like Emails and Phones, so send the full list when updating them.
type CustomerWrite struct {
	FirstName              string          `json:"firstName,omitempty"`
	LastName               string          `json:"lastName,omitempty"`
	BirthDate              string          `json:"birthDate,omitempty"`   // YYYY-MM-DD
	CountryCode            string          `json:"countryCode,omitempty"` // needed when sending phones
	Emails                 []CustomerEmail `json:"emails,omitempty"`
	Phones                 []CustomerPhone `json:"phones,omitempty"`
	EmailMarketingStatus   string          `json:"emailMarketingStatus,omitempty"` // EmailMarketingSubscribed or EmailMarketingUnsubscribed
	IsSendTransactionEmail *bool           `json:"isSendTransactionEmail,omitempty"`
}

type CustomerEmail struct {
	Email string `json:"email"`
}

type CustomerPhone struct {
	Phone string `json:"phone"`
}

const (
	EmailMarketingSubscribed   = "Subscribed"
	EmailMarketingUnsubscribed = "Unsubscribed"
)

type customerPage struct {
	Customers []CustomerFull `json:"customers"`
}

// FindCustomerByEmailContext returns the customer with this email, or nil when there isn't one.
func FindCustomerByEmailContext(ctx context.Context, email, tenant, auth string, retryCount int, rl genericRateLimiter) (*CustomerFull, error) {
	email = strings.TrimSpace(email)
	if email == "" {
		return nil, errors.New("an email is needed to find a customer")
	}
	page, err := GetContext[customerPage](ctx, Endpoints.Customer, map[string]string{"q": email}, nil, tenant, auth, retryCount, rl)
	if err != nil {
		return nil, fmt.Errorf("finding customer: %w", err)
	}
	for i, c := range page.Customers {
		for _, e := range c.Emails {
			if strings.EqualFold(e.Email, email) {
				return &page.Customers[i], nil
			}
		}
	}
	return nil, nil
}

// CreateCustomerContext creates a customer and returns it with its new id.
// Creating isn't idempotent: a retry after a timeout can create the customer twice.
func CreateCustomerContext(ctx context.Context, c CustomerWrite, tenant, auth string, retryCount int, rl genericRateLimiter) (*CustomerFull, error) {
	if len(c.Emails) == 0 && c.FirstName == "" && c.LastName == "" {
		return nil, errors.New("a new customer needs a name or an email")
	}
	body, err := PostContext(ctx, &c, Endpoints.Customer, tenant, auth, retryCount, rl)
	if err != nil {
		return nil, fmt.Errorf("creating customer: %w", err)
	}
	return decodeCustomer(body)
}

// UpdateCustomerContext changes the fields set in c on the customer with this id.
func UpdateCustomerContext(ctx context.Context, id string, c CustomerWrite, tenant, auth string, retryCount int, rl genericRateLimiter) (*CustomerFull, error) {
	if id == "" {
		return nil, errors.New("a customer id is needed to update a customer")
	}
	body, err := PutContext(ctx, &c, Endpoints.Customer+"/"+id, tenant, auth, retryCount, rl)
	if err != nil {
		return nil, fmt.Errorf("updating customer: %w", err)
	}
	return decodeCustomer(body)
}

func decodeCustomer(body *[]byte) (*CustomerFull, error) {
	var out CustomerFull
	if err := json.Unmarshal(*body, &out); err != nil {
		return nil, fmt.Errorf("decoding customer: %w", err)
	}
	return &out, nil
}
