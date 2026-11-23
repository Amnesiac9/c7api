package c7api

import (
	"context"
	"errors"
	"fmt"
	"net/url"
)

// ListMetaDataConfigsContext returns the custom fields set up for an object type, like
// "customer" or "order" (see the MetaDataConfigObject constants).
func ListMetaDataConfigsContext(ctx context.Context, objectType, tenant, auth string, retryCount int, rl genericRateLimiter) ([]MetaDataConfig, error) {
	if !IsValidMetaDataConfigObjectType(objectType) {
		return nil, fmt.Errorf("%q is not an object type with custom fields", objectType)
	}
	page, err := GetContext[MetaDataConfigPayload](ctx, Endpoints.MetaDataConfig+url.PathEscape(objectType), nil, nil, tenant, auth, retryCount, rl)
	if err != nil {
		return nil, fmt.Errorf("listing %s custom fields: %w", objectType, err)
	}
	return page.MetaDataConfigs, nil
}

// SetOrderMetaDataContext sets custom fields on an order, keyed by each field's code. Only the
// fields given are sent; the order's other custom fields are left as they are.
func SetOrderMetaDataContext(ctx context.Context, orderId string, metaData map[string]any, tenant, auth string, retryCount int, rl genericRateLimiter) error {
	return putMetaData(ctx, Endpoints.Order, "order", orderId, metaData, tenant, auth, retryCount, rl)
}

// SetCustomerMetaDataContext sets custom fields on a customer, keyed by each field's code. Only
// the fields given are sent; the customer's other custom fields are left as they are.
func SetCustomerMetaDataContext(ctx context.Context, customerId string, metaData map[string]any, tenant, auth string, retryCount int, rl genericRateLimiter) error {
	return putMetaData(ctx, Endpoints.Customer, "customer", customerId, metaData, tenant, auth, retryCount, rl)
}

func putMetaData(ctx context.Context, endpoint, object, id string, metaData map[string]any, tenant, auth string, retryCount int, rl genericRateLimiter) error {
	if id == "" {
		return fmt.Errorf("no %s id provided", object)
	}
	if len(metaData) == 0 {
		return errors.New("no custom fields provided")
	}
	body := struct {
		MetaData map[string]any `json:"metaData"`
	}{metaData}
	if _, err := PutContext(ctx, &body, endpoint+"/"+id, tenant, auth, retryCount, rl); err != nil {
		return fmt.Errorf("setting %s custom fields: %w", object, err)
	}
	return nil
}
