package c7api

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// A cancelled context must abort an in-flight request rather than run it to
// completion.
func TestContext_CancelInFlight(t *testing.T) {
	released := make(chan struct{})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-released // hold the request open until the test is done
	}))
	defer srv.Close()
	defer close(released)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	start := time.Now()
	_, err := RequestWithRetryAndReadContext(ctx, http.MethodGet, srv.URL, nil, nil, "t", "a", 5, nil)
	elapsed := time.Since(start)

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("err = %v, want context.Canceled", err)
	}
	if elapsed > 2*time.Second {
		t.Errorf("took %v, expected to abort promptly on cancel", elapsed)
	}
}

// The retry loop sleeps between attempts. Those sleeps must be cancellable, or
// a caller with a deadline still waits out the full backoff.
func TestContext_CancelDuringBackoff(t *testing.T) {
	var hits int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits++
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message":"boom"}`))
	}))
	defer srv.Close()

	// retryCount 10 at 500ms per sleep would be ~5s without cancellation.
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	start := time.Now()
	_, err := RequestWithRetryAndReadContext(ctx, http.MethodGet, srv.URL, nil, nil, "t", "a", 10, nil)
	elapsed := time.Since(start)

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("err = %v, want context.DeadlineExceeded", err)
	}
	if elapsed > 2*time.Second {
		t.Errorf("took %v, expected the backoff sleep to be cut short by the deadline", elapsed)
	}
	if hits > 2 {
		t.Errorf("made %d attempts after the deadline passed, expected to stop", hits)
	}
}

// The package client must carry a timeout, so a server that accepts a
// connection and never replies cannot hang the caller forever.
func TestClientTimeout(t *testing.T) {
	if DefaultTimeout <= 0 {
		t.Fatalf("DefaultTimeout = %v, want a positive per-attempt budget", DefaultTimeout)
	}

	released := make(chan struct{})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-released
	}))
	defer srv.Close()
	defer close(released)

	SetHTTPClient(&http.Client{Timeout: 100 * time.Millisecond})
	t.Cleanup(func() { SetHTTPClient(nil) })

	start := time.Now()
	_, err := RequestWithRetryAndReadContext(context.Background(), http.MethodGet, srv.URL, nil, nil, "t", "a", 0, nil)
	elapsed := time.Since(start)

	if err == nil {
		t.Fatal("expected a timeout error, got nil")
	}
	if elapsed > 2*time.Second {
		t.Errorf("took %v, expected the client timeout to fire at ~100ms", elapsed)
	}
}

func TestSetHTTPClient_NilRestoresDefault(t *testing.T) {
	SetHTTPClient(&http.Client{Timeout: time.Second})
	SetHTTPClient(nil)

	if httpClient.Timeout != DefaultTimeout {
		t.Errorf("timeout = %v, want %v", httpClient.Timeout, DefaultTimeout)
	}
}

// The non-retry entry point takes a context too.
func TestRequestContext_Cancel(t *testing.T) {
	released := make(chan struct{})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-released
	}))
	defer srv.Close()
	defer close(released)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	start := time.Now()
	_, err := RequestContext(ctx, http.MethodGet, srv.URL, nil, "t", "a", false)
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if elapsed := time.Since(start); elapsed > 2*time.Second {
		t.Errorf("took %v, expected to abort at the deadline", elapsed)
	}
}
