package c7api

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// fakeTagServer serves one customer tag and records the last request.
func fakeTagServer(t *testing.T, status int) (*httptest.Server, *http.Request, *[]byte) {
	t.Helper()
	var lastReq http.Request
	var lastBody []byte
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lastReq = *r
		lastBody, _ = io.ReadAll(r.Body)
		w.WriteHeader(status)
		if status == http.StatusNotFound {
			w.Write([]byte(`{"statusCode":404,"type":"notFound","message":"Tag not found"}`))
			return
		}
		w.Write([]byte(`{"id":"tag-1","title":"Renamed","objectType":"Customer"}`))
	}))
	real := Endpoints
	Endpoints = GetEndpoints(srv.URL)
	t.Cleanup(func() {
		srv.Close()
		Endpoints = real
	})
	return srv, &lastReq, &lastBody
}

func TestGetTagById(t *testing.T) {
	_, req, _ := fakeTagServer(t, http.StatusOK)

	tag, err := GetTagById("my-tenant", "Basic abc", "customer", "tag-1", 0, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.Method != http.MethodGet || req.URL.Path != "/tag/customer/tag-1" {
		t.Errorf("got %s %s, want GET /tag/customer/tag-1", req.Method, req.URL.Path)
	}
	if tag.ID != "tag-1" {
		t.Errorf("tag id = %q, want tag-1", tag.ID)
	}
}

func TestGetTagById_NotFound(t *testing.T) {
	fakeTagServer(t, http.StatusNotFound)

	_, err := GetTagById("my-tenant", "Basic abc", "customer", "tag-1", 0, nil)
	if !IsNotFound(err) {
		t.Fatalf("want a not found error, got %v", err)
	}
}

func TestUpdateTag(t *testing.T) {
	_, req, body := fakeTagServer(t, http.StatusOK)

	tag, err := UpdateTag("my-tenant", "Basic abc", "Customer", "tag-1", "Renamed", 0, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.Method != http.MethodPut || req.URL.Path != "/tag/customer/tag-1" {
		t.Errorf("got %s %s, want PUT /tag/customer/tag-1", req.Method, req.URL.Path)
	}
	var sent TagPayload_Create
	if err := json.Unmarshal(*body, &sent); err != nil {
		t.Fatalf("bad payload: %v", err)
	}
	if sent.Title != "Renamed" || sent.Type != "Manual" {
		t.Errorf("payload = %+v, want title Renamed and type Manual", sent)
	}
	if tag.Title != "Renamed" {
		t.Errorf("tag title = %q, want Renamed", tag.Title)
	}
}

func TestUpdateTag_RejectsBadInput(t *testing.T) {
	cases := map[string][3]string{
		"bad object type": {"product", "tag-1", "Title"},
		"no tag id":       {"customer", "", "Title"},
		"blank title":     {"customer", "tag-1", "  "},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			if _, err := UpdateTag("my-tenant", "Basic abc", c[0], c[1], c[2], 0, nil); err == nil {
				t.Error("want an error, got nil")
			}
		})
	}
}

func TestDeleteTag(t *testing.T) {
	_, req, _ := fakeTagServer(t, http.StatusOK)

	if err := DeleteTag("my-tenant", "Basic abc", "customer", "tag-1", 0, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.Method != http.MethodDelete || req.URL.Path != "/tag/customer/tag-1" {
		t.Errorf("got %s %s, want DELETE /tag/customer/tag-1", req.Method, req.URL.Path)
	}
}

func TestDeleteTag_AlreadyGone(t *testing.T) {
	fakeTagServer(t, http.StatusNotFound)

	if err := DeleteTag("my-tenant", "Basic abc", "customer", "tag-1", 0, nil); err != nil {
		t.Fatalf("a missing tag should count as deleted, got %v", err)
	}
}
