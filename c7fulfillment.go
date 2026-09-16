package c7api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// FulfillmentItem is one line of a fulfillment that covers specific quantities.
//
// The v1 endpoint matches on Sku; OrderItemId is accepted only by v2. Both are omitted when
// empty so one struct serves either, and so a v1 payload does not carry a key v1 rejects.
type FulfillmentItem struct {
	OrderItemId       string `json:"orderItemId,omitempty"`
	Sku               string `json:"sku,omitempty"`
	QuantityFulfilled int    `json:"quantityFulfilled"`
}

// FulfillmentShipped is the tracking half of a "Shipped" fulfillment.
type FulfillmentShipped struct {
	TrackingNumbers []string `json:"trackingNumbers"`
	Carrier         string   `json:"carrier"`
}

// FulfillmentItems is a fulfillment covering named line items and quantities, posted to
// /order/{:id}/fulfillment.
//
// The difference from FulfillmentAllItems is the endpoint and the granularity.
// FulfillmentAllItems goes to /fulfillment/all and closes out everything still unfulfilled
// on the order — all or nothing. This one closes out only the quantities listed, leaving the
// rest of the order fulfillable, which is what a partial refund or a partial exchange needs.
type FulfillmentItems struct {
	Type                 string              `json:"type"`
	PackageCount         int                 `json:"packageCount"`
	Items                []FulfillmentItem   `json:"items"`
	SendTransactionEmail bool                `json:"sendTransactionEmail"`
	FulfillmentDate      time.Time           `json:"fulfillmentDate"`
	Shipped              *FulfillmentShipped `json:"shipped,omitempty"`
}

// MarkNoFulfillmentRequiredForItems closes out just the given quantities on an order as
// "No Fulfillment Required", leaving everything else on the order still fulfillable.
//
// Use this for a partial refund or partial exchange, where only the returned quantity should
// stop appearing as awaiting fulfillment. MarkNoFulfillmentRequired is the all-or-nothing
// sibling, for an order refunded in its entirety.
//
// Commerce7 answers 422 when a line is already fulfilled, or when the quantity asked for
// exceeds what is left. Callers that treat "already closed out" as success should inspect
// the returned error with errors.As on *C7Error rather than matching on its text.
func MarkNoFulfillmentRequiredForItems(orderId string, items []FulfillmentItem, fulfillmentDate time.Time, tenant string, auth string, attempts int, rl genericRateLimiter) (*[]byte, error) {
	return MarkNoFulfillmentRequiredForItemsContext(context.Background(), orderId, items, fulfillmentDate, tenant, auth, attempts, rl)
}

// MarkNoFulfillmentRequiredForItemsContext is MarkNoFulfillmentRequiredForItems with a
// caller-supplied context. Cancelling ctx aborts the in-flight request and any pending retry
// backoff.
func MarkNoFulfillmentRequiredForItemsContext(ctx context.Context, orderId string, items []FulfillmentItem, fulfillmentDate time.Time, tenant string, auth string, attempts int, rl genericRateLimiter) (*[]byte, error) {
	return PostFulfillmentContext(ctx, orderId, &FulfillmentItems{
		Type:                 OrderFulfillmentTypeNoFulfillmentRequired,
		PackageCount:         1,
		Items:                items,
		SendTransactionEmail: false,
		FulfillmentDate:      fulfillmentDate,
	}, tenant, auth, attempts, rl)
}

// PostFulfillment posts a per-item fulfillment to /order/{:id}/fulfillment.
func PostFulfillment(orderId string, fulfillment *FulfillmentItems, tenant string, auth string, attempts int, rl genericRateLimiter) (*[]byte, error) {
	return PostFulfillmentContext(context.Background(), orderId, fulfillment, tenant, auth, attempts, rl)
}

// PostFulfillmentContext is PostFulfillment with a caller-supplied context.
//
// The error from Commerce7 is passed through with %w rather than flattened into a new
// message: the status code is the only thing that separates "this line is already fulfilled"
// from a real failure, and wrapping is what keeps errors.As able to reach it.
func PostFulfillmentContext(ctx context.Context, orderId string, fulfillment *FulfillmentItems, tenant string, auth string, attempts int, rl genericRateLimiter) (*[]byte, error) {
	// An empty id builds ".../order//fulfillment", which is a different route entirely.
	// Fail here rather than let Commerce7 answer something that looks unrelated.
	if orderId == "" {
		return nil, errors.New("c7api: cannot post a fulfillment with an empty order id")
	}
	if fulfillment == nil {
		return nil, errors.New("c7api: cannot post a nil fulfillment")
	}
	// Commerce7 accepts an empty items array and treats it as a fulfillment of nothing,
	// which closes out no quantity but still edits the order — firing an order update
	// webhook for a post that achieved nothing.
	if len(fulfillment.Items) == 0 {
		return nil, errors.New("c7api: cannot post a fulfillment with no items")
	}
	for i, item := range fulfillment.Items {
		if item.Sku == "" && item.OrderItemId == "" {
			return nil, fmt.Errorf("c7api: fulfillment item %d has neither a sku nor an orderItemId", i)
		}
		if item.QuantityFulfilled <= 0 {
			return nil, fmt.Errorf("c7api: fulfillment item %d (sku %q) has quantityFulfilled %d, must be positive", i, item.Sku, item.QuantityFulfilled)
		}
	}

	body, err := json.Marshal(fulfillment)
	if err != nil {
		return nil, fmt.Errorf("c7api: marshalling fulfillment for order %s: %w", orderId, err)
	}

	url := Endpoints.Order + "/" + orderId + "/fulfillment"
	respBody, err := RequestWithRetryAndReadContext(ctx, http.MethodPost, url, nil, &body, tenant, auth, attempts, rl)
	if err != nil {
		return nil, fmt.Errorf("c7api: posting fulfillment to order %s: %w", orderId, err)
	}

	return respBody, nil
}
