package c7api

import (
	"context"
	"errors"
	"net/http"
	"testing"
)

func TestListSalesAttributes(t *testing.T) {
	reqs := fakeC7(t, func(r *http.Request) (int, string) {
		return http.StatusOK, `{"salesAttributes":[
			{"id":"a1","title":"Web","code":"Web","isDeletable":false},
			{"id":"a2","title":"Wine Club","code":"CLUB"}
		],"total":2}`
	})

	attrs, err := ListSalesAttributesContext(context.Background(), "my-tenant", "Basic abc", 0, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(*reqs) != 1 || (*reqs)[0].path != "/sales-attribute" {
		t.Fatalf("want one GET /sales-attribute, got %+v", *reqs)
	}
	if len(attrs) != 2 || attrs[1].Title != "Wine Club" || attrs[1].Code != "CLUB" {
		t.Fatalf("unexpected attributes: %+v", attrs)
	}
}

func TestSetOrderSalesAttributeSendsOnlyTheCode(t *testing.T) {
	reqs := fakeC7(t, func(r *http.Request) (int, string) {
		return http.StatusOK, `{"id":"o1","salesAttributeCode":"CLUB"}`
	})

	if err := SetOrderSalesAttributeContext(context.Background(), "o1", "CLUB", "my-tenant", "Basic abc", 0, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(*reqs) != 1 || (*reqs)[0].method != http.MethodPut || (*reqs)[0].path != "/order/o1" {
		t.Fatalf("want one PUT /order/o1, got %+v", *reqs)
	}
	if got := string((*reqs)[0].body); got != `{"salesAttributeCode":"CLUB"}` {
		t.Errorf("body = %s", got)
	}
}

func TestSetOrderSalesAttributeNeedsIdAndCode(t *testing.T) {
	ctx := context.Background()
	if SetOrderSalesAttributeContext(ctx, "", "CLUB", "t", "a", 0, nil) == nil {
		t.Error("no order id should fail")
	}
	if SetOrderSalesAttributeContext(ctx, "o1", "", "t", "a", 0, nil) == nil {
		t.Error("no code should fail")
	}
}

func TestSetOrderSalesAttributeReportsCommerce7Errors(t *testing.T) {
	fakeC7(t, func(r *http.Request) (int, string) {
		return http.StatusUnauthorized, `{"statusCode":401,"type":"unauthorized","message":"Unauthorized User for Route"}`
	})
	err := SetOrderSalesAttributeContext(context.Background(), "o1", "CLUB", "t", "a", 0, nil)
	var c7err *C7Error
	if !errors.As(err, &c7err) || c7err.StatusCode != http.StatusUnauthorized {
		t.Errorf("want a 401 C7Error, got %T: %v", err, err)
	}
}
