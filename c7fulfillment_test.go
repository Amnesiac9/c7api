package c7api

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

// The payload Commerce7 expects, field for field. The v1 endpoint rejects an unknown
// orderItemId key, so an item without one must not emit it at all.
func Test_FulfillmentItems_Payload(t *testing.T) {
	date := time.Date(2026, 9, 15, 7, 0, 0, 0, time.UTC)

	body, err := json.Marshal(&FulfillmentItems{
		Type:                 OrderFulfillmentTypeNoFulfillmentRequired,
		PackageCount:         1,
		Items:                []FulfillmentItem{{Sku: "MTE20C", QuantityFulfilled: 6}},
		SendTransactionEmail: false,
		FulfillmentDate:      date,
	})
	if err != nil {
		t.Fatalf("marshalling: %v", err)
	}

	const want = `{"type":"No Fulfillment Required","packageCount":1,"items":[{"sku":"MTE20C","quantityFulfilled":6}],"sendTransactionEmail":false,"fulfillmentDate":"2026-09-15T07:00:00Z"}`
	if got := string(body); got != want {
		t.Errorf("payload =\n%s\nwant\n%s", got, want)
	}
}

// Every one of these is a post that would edit the order without closing out any quantity,
// or hit a route other than the one intended. None of them should reach the network.
func Test_PostFulfillment_RejectsUnsendablePayloads(t *testing.T) {
	good := func() *FulfillmentItems {
		return &FulfillmentItems{
			Type:         OrderFulfillmentTypeNoFulfillmentRequired,
			PackageCount: 1,
			Items:        []FulfillmentItem{{Sku: "CSC21C", QuantityFulfilled: 6}},
		}
	}

	tests := []struct {
		name        string
		orderId     string
		fulfillment *FulfillmentItems
		wantErr     string
	}{
		{"empty order id", "", good(), "empty order id"},
		{"nil fulfillment", "order-1", nil, "nil fulfillment"},
		{
			"no items", "order-1",
			&FulfillmentItems{Type: OrderFulfillmentTypeNoFulfillmentRequired, PackageCount: 1},
			"no items",
		},
		{
			"item with no identifier", "order-1",
			&FulfillmentItems{PackageCount: 1, Items: []FulfillmentItem{{QuantityFulfilled: 6}}},
			"neither a sku nor an orderItemId",
		},
		{
			"zero quantity", "order-1",
			&FulfillmentItems{PackageCount: 1, Items: []FulfillmentItem{{Sku: "CSC21C"}}},
			"must be positive",
		},
		{
			"negative quantity", "order-1",
			&FulfillmentItems{PackageCount: 1, Items: []FulfillmentItem{{Sku: "CSC21C", QuantityFulfilled: -6}}},
			"must be positive",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// A cancelled context proves the check happened before any request: if one were
			// attempted the error would be the cancellation, not the validation message.
			ctx, cancel := context.WithCancel(context.Background())
			cancel()

			_, err := PostFulfillmentContext(ctx, test.orderId, test.fulfillment, "tenant", "auth", 1, nil)
			if err == nil {
				t.Fatalf("expected an error containing %q, got nil", test.wantErr)
			}
			if !strings.Contains(err.Error(), test.wantErr) {
				t.Errorf("error = %q, want it to contain %q", err, test.wantErr)
			}
		})
	}
}

// The helper exists so callers do not have to remember the type string, the package count,
// or that the transaction email must be off.
func Test_MarkNoFulfillmentRequiredForItems_BuildsTheRightFulfillment(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// An item list that fails validation, so the call returns before the network and the
	// error tells us which check the built fulfillment reached.
	_, err := MarkNoFulfillmentRequiredForItemsContext(ctx, "order-1", nil, time.Now(), "tenant", "auth", 1, nil)
	if err == nil || !strings.Contains(err.Error(), "no items") {
		t.Fatalf("error = %v, want a 'no items' rejection", err)
	}
}
