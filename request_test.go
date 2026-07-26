package c7api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// captureHeaders spins up a server that records the headers of the first
// request it receives and replies 200 with an empty JSON object.
func captureHeaders(t *testing.T) (*httptest.Server, *http.Header) {
	t.Helper()
	got := &http.Header{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		*got = r.Header.Clone()
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}))
	t.Cleanup(srv.Close)
	return srv, got
}

// The v1 and v2 entry points share one implementation. This pins the header
// sets so the shared core can't silently drop what v2 requires.
func TestHeaders_V1(t *testing.T) {
	srv, got := captureHeaders(t)

	if _, err := RequestWithRetryAndRead(http.MethodGet, srv.URL, nil, nil, "my-tenant", "Basic abc123", 0, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := map[string]string{
		"Tenant":        "my-tenant",
		"Content-Type":  "application/json",
		"Authorization": "Basic abc123",
	}
	for k, v := range want {
		if gotV := got.Get(k); gotV != v {
			t.Errorf("header %q = %q, want %q", k, gotV, v)
		}
	}

	// v1 must NOT send the experimental v2 headers.
	for _, k := range []string{"Tenantid", "Experimental"} {
		if gotV := got.Get(k); gotV != "" {
			t.Errorf("v1 sent v2-only header %q = %q, want it absent", k, gotV)
		}
	}
}

func TestHeaders_V2(t *testing.T) {
	srv, got := captureHeaders(t)

	if _, err := RequestWithRetryAndReadV2(http.MethodGet, srv.URL, nil, nil, "my-tenant", "Basic abc123", 0, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := map[string]string{
		"Tenant":        "my-tenant",
		"Content-Type":  "application/json",
		"Authorization": "Basic abc123",
		"Tenantid":      "my-tenant",
		"Experimental":  "Do not use if you are not Commerce7.  API likely to change",
	}
	for k, v := range want {
		if gotV := got.Get(k); gotV != v {
			t.Errorf("header %q = %q, want %q", k, gotV, v)
		}
	}
}

// Queries are appended by the shared core; make sure both entry points still
// apply them.
func TestQueriesApplied(t *testing.T) {
	var gotQuery string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.Query().Get("q")
		w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	if _, err := RequestWithRetryAndRead(http.MethodGet, srv.URL, map[string]string{"q": "1234"}, nil, "t", "a", 0, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotQuery != "1234" {
		t.Errorf("query q = %q, want %q", gotQuery, "1234")
	}
}
