package c7api

import (
	"context"
	"errors"
	"fmt"
)

// SalesAttribute is a value an order can carry in salesAttributeCode, like "Web" or "POS".
type SalesAttribute struct {
	Id    string `json:"id"`
	Title string `json:"title"`
	Code  string `json:"code"`
}

// ListSalesAttributesContext returns the tenant's sales attributes.
func ListSalesAttributesContext(ctx context.Context, tenant, auth string, retryCount int, rl genericRateLimiter) ([]SalesAttribute, error) {
	items, err := GetAllContext[SalesAttribute, ListPage[SalesAttribute]](ctx, Endpoints.SalesAttribute, nil, nil, tenant, auth, retryCount, rl)
	if err != nil {
		return nil, fmt.Errorf("listing sales attributes: %w", err)
	}
	return *items, nil
}

// SetOrderSalesAttributeContext sets an order's sales attribute by its code. Only that field is
// sent, so the rest of the order is left as it is.
func SetOrderSalesAttributeContext(ctx context.Context, orderId, code, tenant, auth string, retryCount int, rl genericRateLimiter) error {
	if orderId == "" {
		return errors.New("no order id provided")
	}
	if code == "" {
		return errors.New("no sales attribute code provided")
	}
	body := struct {
		SalesAttributeCode string `json:"salesAttributeCode"`
	}{code}
	if _, err := PutContext(ctx, &body, Endpoints.Order+"/"+orderId, tenant, auth, retryCount, rl); err != nil {
		return fmt.Errorf("setting order sales attribute: %w", err)
	}
	return nil
}
