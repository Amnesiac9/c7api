package c7api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

// countAttempts serves status for every request and reports how many requests
// arrived, plus the error the call returned.
func countAttempts(t *testing.T, status int, retryCount int) (int32, error) {
	t.Helper()

	var hits int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&hits, 1)
		w.WriteHeader(status)
		// No statusCode field, so this also covers the fallback to the HTTP
		// status added alongside the original 404 fix.
		w.Write([]byte(`{"message":"nope"}`))
	}))
	defer srv.Close()

	_, err := RequestWithRetryAndRead(http.MethodGet, srv.URL, nil, nil, "t", "a", retryCount, nil)
	return atomic.LoadInt32(&hits), err
}

func TestRetry_StatusClassification(t *testing.T) {
	const retryCount = 1 // keep the suite quick; 1 retry is enough to tell

	tests := []struct {
		name     string
		status   int
		wantHits int32
	}{
		// Nothing a retry can change.
		{"400 bad request", http.StatusBadRequest, 1},
		{"401 unauthorized", http.StatusUnauthorized, 1},
		{"403 forbidden", http.StatusForbidden, 1},
		{"404 not found", http.StatusNotFound, 1},
		{"422 unprocessable", http.StatusUnprocessableEntity, 1},

		// Transient.
		{"408 request timeout", http.StatusRequestTimeout, 2},
		{"429 too many requests", http.StatusTooManyRequests, 2},
		{"500 internal", http.StatusInternalServerError, 2},
		{"502 bad gateway", http.StatusBadGateway, 2},
		{"503 unavailable", http.StatusServiceUnavailable, 2},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			hits, err := countAttempts(t, tc.status, retryCount)

			if hits != tc.wantHits {
				t.Errorf("made %d requests, want %d", hits, tc.wantHits)
			}

			var c7err *C7Error
			if !errors.As(err, &c7err) {
				t.Fatalf("err = %T (%v), want *C7Error", err, err)
			}
			if c7err.StatusCode != tc.status {
				t.Errorf("C7Error.StatusCode = %d, want %d", c7err.StatusCode, tc.status)
			}
		})
	}
}

// Regression guard for b791434: a 404 must cost exactly one request even with
// the retry budget wide open.
func TestRetry_404FailsFastAtMaxRetries(t *testing.T) {
	hits, err := countAttempts(t, http.StatusNotFound, 10)

	if hits != 1 {
		t.Errorf("made %d requests for a 404, want 1", hits)
	}

	var c7err *C7Error
	if !errors.As(err, &c7err) {
		t.Fatalf("err = %T (%v), want *C7Error", err, err)
	}
	if c7err.StatusCode != http.StatusNotFound {
		t.Errorf("C7Error.StatusCode = %d, want 404", c7err.StatusCode)
	}
}

// A retryable status that clears must return the eventual success, not the
// earlier failure.
func TestRetry_RecoversAfter429(t *testing.T) {
	var hits int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if atomic.AddInt32(&hits, 1) < 3 {
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`{"message":"slow down"}`))
			return
		}
		w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()

	body, err := RequestWithRetryAndRead(http.MethodGet, srv.URL, nil, nil, "t", "a", 3, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := string(*body); got != `{"ok":true}` {
		t.Errorf("body = %q, want the successful response", got)
	}
	if hits != 3 {
		t.Errorf("made %d requests, want 3", hits)
	}
}

// hijackAndDrop answers every request by killing the connection, which
// surfaces to the client as a transport error rather than a status code.
func hijackAndDrop(hits *int32) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(hits, 1)
		conn, _, err := w.(http.Hijacker).Hijack()
		if err != nil {
			return
		}
		conn.Close()
	}
}

// Transport errors used to abort on the first failure with zero retries — the
// inverse of the right policy, and newly reachable now that a client timeout
// can produce them.
func TestRetry_TransportErrorIsRetried(t *testing.T) {
	var hits int32
	srv := httptest.NewServer(hijackAndDrop(&hits))
	defer srv.Close()

	_, err := RequestWithRetryAndRead(http.MethodGet, srv.URL, nil, nil, "t", "a", 2, nil)

	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if hits != 3 {
		t.Errorf("made %d requests, want 3 (1 + 2 retries)", hits)
	}

	// No response body ever arrived, so there is no C7 error message to report.
	var c7err *C7Error
	if errors.As(err, &c7err) {
		t.Errorf("err = %v, want a plain transport error rather than a C7Error", err)
	}
}

// A transport failure that clears on a later attempt must succeed.
func TestRetry_RecoversAfterTransportError(t *testing.T) {
	var hits int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if atomic.AddInt32(&hits, 1) < 2 {
			conn, _, err := w.(http.Hijacker).Hijack()
			if err != nil {
				return
			}
			conn.Close()
			return
		}
		w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()

	body, err := RequestWithRetryAndRead(http.MethodGet, srv.URL, nil, nil, "t", "a", 2, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := string(*body); got != `{"ok":true}` {
		t.Errorf("body = %q, want the successful response", got)
	}
}

// A transport error on a later attempt must not resurrect an earlier response,
// and must not panic on the nil response Do returns with its error.
func TestRetry_TransportErrorAfterStatusFailure(t *testing.T) {
	var hits int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if atomic.AddInt32(&hits, 1) == 1 {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message":"boom"}`))
			return
		}
		conn, _, err := w.(http.Hijacker).Hijack()
		if err != nil {
			return
		}
		conn.Close()
	}))
	defer srv.Close()

	_, err := RequestWithRetryAndRead(http.MethodGet, srv.URL, nil, nil, "t", "a", 1, nil)
	if err == nil {
		t.Fatal("expected an error, got nil")
	}

	var c7err *C7Error
	if errors.As(err, &c7err) {
		t.Errorf("err = %v, want the last failure (transport), not the earlier 500", err)
	}
}

// The classification applies to v2 as well, since both share one core.
func TestRetry_V2FailsFastOn404(t *testing.T) {
	var hits int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&hits, 1)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message":"missing"}`))
	}))
	defer srv.Close()

	if _, err := RequestWithRetryAndReadV2(http.MethodGet, srv.URL, nil, nil, "t", "a", 5, nil); err == nil {
		t.Fatal("expected an error, got nil")
	}
	if hits != 1 {
		t.Errorf("v2 made %d requests for a 404, want 1", hits)
	}
}

// A 401 with the retry budget wide open should return essentially immediately.
func TestRetry_NonRetryableReturnsPromptly(t *testing.T) {
	start := time.Now()
	if _, err := countAttempts(t, http.StatusUnauthorized, 10); err == nil {
		t.Fatal("expected an error, got nil")
	}
	if elapsed := time.Since(start); elapsed > time.Second {
		t.Errorf("took %v, want a near-immediate return for a 401", elapsed)
	}
}
