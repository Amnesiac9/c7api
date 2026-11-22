package c7api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// fakeCustomerServer answers customer requests with body and records the last request.
func fakeCustomerServer(t *testing.T, status int, body string) (*http.Request, *[]byte) {
	t.Helper()
	var lastReq http.Request
	var lastBody []byte
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lastReq = *r
		lastBody, _ = io.ReadAll(r.Body)
		w.WriteHeader(status)
		w.Write([]byte(body))
	}))
	real := Endpoints
	Endpoints = GetEndpoints(srv.URL)
	t.Cleanup(func() {
		srv.Close()
		Endpoints = real
	})
	return &lastReq, &lastBody
}

func TestFindCustomerByEmail(t *testing.T) {
	req, _ := fakeCustomerServer(t, http.StatusOK, `{"customers":[
		{"id":"c1","emails":[{"email":"other@example.com"}]},
		{"id":"c2","emails":[{"email":"Jane@Example.com"}]}]}`)

	c, err := FindCustomerByEmailContext(context.Background(), "jane@example.com", "t", "Basic x", 0, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil || c.Id != "c2" {
		t.Fatalf("got %+v, want customer c2", c)
	}
	if req.URL.Path != "/customer" || req.URL.Query().Get("q") != "jane@example.com" {
		t.Errorf("got %s?%s", req.URL.Path, req.URL.RawQuery)
	}
}

func TestFindCustomerByEmail_None(t *testing.T) {
	fakeCustomerServer(t, http.StatusOK, `{"customers":[{"id":"c1","emails":[{"email":"jane.doe@example.com"}]}]}`)

	c, err := FindCustomerByEmailContext(context.Background(), "jane@example.com", "t", "Basic x", 0, nil)
	if err != nil || c != nil {
		t.Fatalf("got %+v, %v; want no customer and no error", c, err)
	}
}

func TestCreateCustomer(t *testing.T) {
	req, body := fakeCustomerServer(t, http.StatusOK, `{"id":"new-1","firstName":"Jane"}`)

	c, err := CreateCustomerContext(context.Background(), CustomerWrite{
		FirstName: "Jane", Emails: []CustomerEmail{{Email: "jane@example.com"}}, EmailMarketingStatus: EmailMarketingSubscribed,
	}, "t", "Basic x", 0, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Id != "new-1" {
		t.Errorf("id = %q, want new-1", c.Id)
	}
	if req.Method != http.MethodPost || req.URL.Path != "/customer" {
		t.Errorf("got %s %s, want POST /customer", req.Method, req.URL.Path)
	}
	var sent map[string]any
	json.Unmarshal(*body, &sent)
	if _, ok := sent["lastName"]; ok {
		t.Errorf("empty fields should be left out, sent %s", *body)
	}
	if sent["emailMarketingStatus"] != "Subscribed" {
		t.Errorf("sent %s", *body)
	}
}

func TestCreateCustomer_NeedsNameOrEmail(t *testing.T) {
	if _, err := CreateCustomerContext(context.Background(), CustomerWrite{}, "t", "Basic x", 0, nil); err == nil {
		t.Fatal("want an error for an empty customer")
	}
}

func TestUpdateCustomer(t *testing.T) {
	req, _ := fakeCustomerServer(t, http.StatusOK, `{"id":"c9"}`)

	if _, err := UpdateCustomerContext(context.Background(), "c9", CustomerWrite{BirthDate: "1990-02-03"}, "t", "Basic x", 0, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.Method != http.MethodPut || req.URL.Path != "/customer/c9" {
		t.Errorf("got %s %s, want PUT /customer/c9", req.Method, req.URL.Path)
	}
}

func TestUpdateCustomer_KeepsEmailIds(t *testing.T) {
	_, body := fakeCustomerServer(t, http.StatusOK, `{"id":"c9"}`)

	w := CustomerWrite{Emails: []CustomerEmail{{Id: "e1", Email: "a@example.com", Status: "Valid"}}, Phones: []CustomerPhone{{Phone: "555"}}}
	if _, err := UpdateCustomerContext(context.Background(), "c9", w, "t", "Basic x", 0, nil); err != nil {
		t.Fatal(err)
	}
	want := `"emails":[{"id":"e1","email":"a@example.com","status":"Valid"}],"phones":[{"phone":"555"}]`
	if !strings.Contains(string(*body), want) {
		t.Errorf("sent %s", *body)
	}
}
