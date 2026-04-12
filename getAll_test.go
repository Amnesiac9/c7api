package c7api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

// ── test item + wrapper types ──────────────────────────────────────

type testProduct struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Paginator wrapper
type testProductsPage struct {
	Products []testProduct `json:"products"`
	Total    int           `json:"total"`
}

func (p testProductsPage) GetItems() []testProduct { return p.Products }
func (p testProductsPage) GetTotal() int           { return p.Total }

// Cursornator wrapper
type testProductsCursor struct {
	Products []testProduct `json:"products"`
	Cursor   string        `json:"cursor"`
}

func (p testProductsCursor) GetItems() []testProduct { return p.Products }
func (p testProductsCursor) GetCursor() string       { return p.Cursor }

// ── helpers ────────────────────────────────────────────────────────

// makeProducts builds n testProducts with sequential IDs.
func makeProducts(n int) []testProduct {
	out := make([]testProduct, n)
	for i := range out {
		out[i] = testProduct{
			ID:   fmt.Sprintf("id-%d", i+1),
			Name: fmt.Sprintf("Product %d", i+1),
		}
	}
	return out
}

// ── GetAll (page-based) ────────────────────────────────────────────

func TestGetAll_MultiplePages(t *testing.T) {
	// 120 products → 3 pages (PageSize is 50)
	allProducts := makeProducts(120)
	total := len(allProducts)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		if page < 1 {
			page = 1
		}

		start := (page - 1) * PageSize
		end := start + PageSize
		if end > total {
			end = total
		}

		resp := testProductsPage{
			Products: allProducts[start:end],
			Total:    total,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	got, err := GetAll[testProduct, testProductsPage](
		srv.URL+"/product",
		nil,                   // no extra queries
		nil,                   // no body
		"test-tenant",         // tenant
		"Basic dGVzdDp0ZXN0", // auth
		0,                     // retryCount
		&rateLimiterMock{},    // rate limiter
	)
	if err != nil {
		t.Fatalf("GetAll returned error: %v", err)
	}

	if len(*got) != total {
		t.Fatalf("expected %d products, got %d", total, len(*got))
	}

	// Spot-check first and last item
	if (*got)[0].ID != "id-1" {
		t.Errorf("first item ID = %q, want %q", (*got)[0].ID, "id-1")
	}
	if (*got)[total-1].ID != fmt.Sprintf("id-%d", total) {
		t.Errorf("last item ID = %q, want %q", (*got)[total-1].ID, fmt.Sprintf("id-%d", total))
	}
}

func TestGetAll_SinglePage(t *testing.T) {
	// Fewer items than PageSize → single request
	allProducts := makeProducts(10)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := testProductsPage{
			Products: allProducts,
			Total:    len(allProducts),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	got, err := GetAll[testProduct, testProductsPage](
		srv.URL+"/product",
		nil, nil, "test-tenant", "Basic dGVzdDp0ZXN0", 0, &rateLimiterMock{},
	)
	if err != nil {
		t.Fatalf("GetAll returned error: %v", err)
	}
	if len(*got) != 10 {
		t.Fatalf("expected 10 products, got %d", len(*got))
	}
}

func TestGetAll_EmptyResult(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := testProductsPage{Products: []testProduct{}, Total: 0}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	got, err := GetAll[testProduct, testProductsPage](
		srv.URL+"/product",
		nil, nil, "test-tenant", "Basic dGVzdDp0ZXN0", 0, &rateLimiterMock{},
	)
	if err != nil {
		t.Fatalf("GetAll returned error: %v", err)
	}
	if len(*got) != 0 {
		t.Fatalf("expected 0 products, got %d", len(*got))
	}
}

// ── GetAllWithCursor ───────────────────────────────────────────────

func TestGetAllWithCursor_MultiplePages(t *testing.T) {
	allProducts := makeProducts(120)
	chunkSize := PageSize

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cursor := r.URL.Query().Get("cursor")

		// Determine offset from cursor
		offset := 0
		if cursor != "" && cursor != "start" {
			offset, _ = strconv.Atoi(cursor)
		}

		end := offset + chunkSize
		nextCursor := ""
		if end < len(allProducts) {
			nextCursor = strconv.Itoa(end)
		} else {
			end = len(allProducts)
		}

		resp := testProductsCursor{
			Products: allProducts[offset:end],
			Cursor:   nextCursor,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	got, err := GetAllWithCursor[testProduct, testProductsCursor](
		srv.URL+"/product",
		nil,                   // no extra queries
		nil,                   // no body
		"test-tenant",         // tenant
		"Basic dGVzdDp0ZXN0", // auth
		0,                     // retryCount
	)
	if err != nil {
		t.Fatalf("GetAllWithCursor returned error: %v", err)
	}

	if len(*got) != len(allProducts) {
		t.Fatalf("expected %d products, got %d", len(allProducts), len(*got))
	}

	if (*got)[0].ID != "id-1" {
		t.Errorf("first item ID = %q, want %q", (*got)[0].ID, "id-1")
	}
	if (*got)[119].ID != "id-120" {
		t.Errorf("last item ID = %q, want %q", (*got)[119].ID, "id-120")
	}
}
