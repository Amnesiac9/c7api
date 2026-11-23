package c7api

import (
	"context"
	"net/http"
	"testing"
)

func TestGetByIdContext(t *testing.T) {
	reqs := fakeC7(t, func(r *http.Request) (int, string) {
		return http.StatusOK, `{"id":"x1","orderInformation":{"lifetimeValue":250000}}`
	})
	ctx := context.Background()

	customer, err := GetCustomerByIdContext[map[string]any](ctx, "c1", "my-tenant", "Basic abc", 0, nil)
	if err != nil {
		t.Fatalf("customer: %v", err)
	}
	if info, _ := (*customer)["orderInformation"].(map[string]any); info["lifetimeValue"] != 250000.0 {
		t.Errorf("customer fields lost: %v", *customer)
	}
	if _, err := GetOrderByIdContext[map[string]any](ctx, "o1", "my-tenant", "Basic abc", 0, nil); err != nil {
		t.Fatalf("order: %v", err)
	}
	if _, err := GetClubMembershipByIdContext[map[string]any](ctx, "m/1", "my-tenant", "Basic abc", 0, nil); err != nil {
		t.Fatalf("membership: %v", err)
	}

	paths := []string{}
	for _, r := range *reqs {
		if r.method != http.MethodGet {
			t.Errorf("method %s", r.method)
		}
		paths = append(paths, r.path)
	}
	want := []string{"/customer/c1", "/order/o1", "/club-membership/m/1"}
	if len(paths) != 3 || paths[0] != want[0] || paths[1] != want[1] || paths[2] != want[2] {
		t.Errorf("paths = %v, want %v", paths, want)
	}
}

func TestGetByIdContextNotFoundAndNoId(t *testing.T) {
	fakeC7(t, func(r *http.Request) (int, string) {
		return http.StatusNotFound, `{"statusCode":404,"type":"notFound","message":"not found"}`
	})
	_, err := GetCustomerByIdContext[map[string]any](context.Background(), "gone", "t", "a", 0, nil)
	if !IsNotFound(err) {
		t.Errorf("want a not found error, got %v", err)
	}
	if _, err := GetOrderByIdContext[map[string]any](context.Background(), "", "t", "a", 0, nil); err == nil {
		t.Error("no id should fail")
	}
}
