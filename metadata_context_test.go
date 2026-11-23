package c7api

import (
	"context"
	"net/http"
	"testing"
)

func TestListMetaDataConfigs(t *testing.T) {
	reqs := fakeC7(t, func(r *http.Request) (int, string) {
		return http.StatusOK, `{"metaDataConfigs":[
			{"id":"m1","title":"Lifetime value","objectType":"customer","code":"wp-ltv","dataType":"String"},
			{"id":"m2","title":"Orders","objectType":"customer","code":"wp-total-orders","dataType":"Number"}
		],"total":2}`
	})

	fields, err := ListMetaDataConfigsContext(context.Background(), "customer", "my-tenant", "Basic abc", 0, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(*reqs) != 1 || (*reqs)[0].method != http.MethodGet || (*reqs)[0].path != "/meta-data-config/customer" {
		t.Fatalf("want one GET /meta-data-config/customer, got %+v", *reqs)
	}
	if len(fields) != 2 || fields[1].Code != "wp-total-orders" || fields[1].DataType != MetaDataTypeNumber {
		t.Fatalf("unexpected fields: %+v", fields)
	}
}

func TestListMetaDataConfigsRejectsUnknownObjects(t *testing.T) {
	if _, err := ListMetaDataConfigsContext(context.Background(), "spaceship", "t", "a", 0, nil); err == nil {
		t.Error("an unknown object type should fail")
	}
}

func TestSetMetaDataSendsOnlyTheFieldsGiven(t *testing.T) {
	reqs := fakeC7(t, func(r *http.Request) (int, string) {
		return http.StatusOK, `{"id":"x"}`
	})
	ctx := context.Background()

	if err := SetOrderMetaDataContext(ctx, "o1", map[string]any{"gift-note": "Happy birthday"}, "t", "a", 0, nil); err != nil {
		t.Fatalf("order: %v", err)
	}
	if err := SetCustomerMetaDataContext(ctx, "c1", map[string]any{"wp-total-orders": 3}, "t", "a", 0, nil); err != nil {
		t.Fatalf("customer: %v", err)
	}
	want := []recordedRequest{
		{method: http.MethodPut, path: "/order/o1", body: []byte(`{"metaData":{"gift-note":"Happy birthday"}}`)},
		{method: http.MethodPut, path: "/customer/c1", body: []byte(`{"metaData":{"wp-total-orders":3}}`)},
	}
	if len(*reqs) != len(want) {
		t.Fatalf("want %d requests, got %+v", len(want), *reqs)
	}
	for i, w := range want {
		got := (*reqs)[i]
		if got.method != w.method || got.path != w.path || string(got.body) != string(w.body) {
			t.Errorf("request %d = %s %s %s, want %s %s %s", i, got.method, got.path, got.body, w.method, w.path, w.body)
		}
	}
}

func TestSetMetaDataNeedsAnIdAndFields(t *testing.T) {
	ctx := context.Background()
	if SetOrderMetaDataContext(ctx, "", map[string]any{"a": 1}, "t", "a", 0, nil) == nil {
		t.Error("no order id should fail")
	}
	if SetCustomerMetaDataContext(ctx, "c1", nil, "t", "a", 0, nil) == nil {
		t.Error("no fields should fail")
	}
}
