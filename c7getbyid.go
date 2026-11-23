package c7api

import (
	"context"
	"errors"
	"net/url"
)

// Get one record by its id. T is the struct to read it into, or map[string]any to keep every
// field. A record that doesn't exist returns an error that IsNotFound reports.

// GetCustomerByIdContext gets one customer.
func GetCustomerByIdContext[T any](ctx context.Context, customerId, tenant, c7AppAuthEncoded string, retryCount int, rl genericRateLimiter) (*T, error) {
	return getByIdContext[T](ctx, Endpoints.Customer, customerId, tenant, c7AppAuthEncoded, retryCount, rl)
}

// GetOrderByIdContext gets one order.
func GetOrderByIdContext[T any](ctx context.Context, orderId, tenant, c7AppAuthEncoded string, retryCount int, rl genericRateLimiter) (*T, error) {
	return getByIdContext[T](ctx, Endpoints.Order, orderId, tenant, c7AppAuthEncoded, retryCount, rl)
}

// GetClubMembershipByIdContext gets one club membership.
func GetClubMembershipByIdContext[T any](ctx context.Context, membershipId, tenant, c7AppAuthEncoded string, retryCount int, rl genericRateLimiter) (*T, error) {
	return getByIdContext[T](ctx, Endpoints.ClubMembership, membershipId, tenant, c7AppAuthEncoded, retryCount, rl)
}

func getByIdContext[T any](ctx context.Context, endpoint, id, tenant, c7AppAuthEncoded string, retryCount int, rl genericRateLimiter) (*T, error) {
	if id == "" {
		return nil, errors.New("no id provided")
	}
	return GetContext[T](ctx, endpoint+"/"+url.PathEscape(id), nil, nil, tenant, c7AppAuthEncoded, retryCount, rl)
}
