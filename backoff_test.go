package c7api

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func TestBackoffDuration(t *testing.T) {
	tests := []struct {
		attempt int
		want    time.Duration
	}{
		{-1, 500 * time.Millisecond}, // defensive, never happens in the loop
		{0, 500 * time.Millisecond},  // was 0 for 429s, hammering the limiter
		{1, 1 * time.Second},
		{2, 2 * time.Second},
		{3, 4 * time.Second},
		{4, MaxBackoff},
		{10, MaxBackoff}, // the highest attempt index the loop allows
		{62, MaxBackoff}, // shift overflow must still land on the cap
	}

	for _, tc := range tests {
		if got := backoffDuration(tc.attempt); got != tc.want {
			t.Errorf("backoffDuration(%d) = %v, want %v", tc.attempt, got, tc.want)
		}
	}
}

// retryTimer counts requests to a server that always fails with status, and
// reports how long the whole call took.
func retryTimer(t *testing.T, status int, retryCount int) (int32, time.Duration) {
	t.Helper()

	var hits int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&hits, 1)
		w.WriteHeader(status)
		w.Write([]byte(`{"message":"nope"}`))
	}))
	defer srv.Close()

	start := time.Now()
	_, err := RequestWithRetryAndRead(http.MethodGet, srv.URL, nil, nil, "t", "a", retryCount, nil)
	elapsed := time.Since(start)

	if err == nil {
		t.Fatalf("expected an error for status %d", status)
	}
	return atomic.LoadInt32(&hits), elapsed
}

// The loop used to sleep after the final attempt, delaying the error return by
// a full backoff interval for no benefit.
func TestBackoff_NoSleepAfterFinalAttempt(t *testing.T) {
	hits, elapsed := retryTimer(t, http.StatusInternalServerError, 1)

	if hits != 2 {
		t.Fatalf("made %d requests, want 2 (1 + 1 retry)", hits)
	}
	// One backoff of 500ms between the two attempts, and nothing after the
	// second. Anything near 1s means the trailing sleep is back.
	if elapsed > 900*time.Millisecond {
		t.Errorf("took %v, want ~500ms — looks like it slept after the last attempt", elapsed)
	}
	if elapsed < 500*time.Millisecond {
		t.Errorf("took %v, want at least one 500ms backoff between attempts", elapsed)
	}
}

// 429 is the case the backoff exists for. The old formula slept
// SLEEP_TIME * i, so the first retry left with no delay at all.
func TestBackoff_FirstRetryOf429Waits(t *testing.T) {
	hits, elapsed := retryTimer(t, http.StatusTooManyRequests, 1)

	if hits != 2 {
		t.Fatalf("made %d requests, want 2", hits)
	}
	if elapsed < 500*time.Millisecond {
		t.Errorf("took %v, want the first 429 retry to wait 500ms", elapsed)
	}
}

// Successive retries must grow, not stay flat.
func TestBackoff_IsExponential(t *testing.T) {
	hits, elapsed := retryTimer(t, http.StatusInternalServerError, 3)

	if hits != 4 {
		t.Fatalf("made %d requests, want 4 (1 + 3 retries)", hits)
	}
	// 500ms + 1s + 2s = 3.5s of backoff. The old flat 500ms gave 1.5s.
	if elapsed < 3*time.Second {
		t.Errorf("took %v, want ~3.5s of exponential backoff", elapsed)
	}
	if elapsed > 6*time.Second {
		t.Errorf("took %v, want ~3.5s — backoff is growing faster than expected", elapsed)
	}
}
