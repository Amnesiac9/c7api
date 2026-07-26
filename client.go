package c7api

import (
	"context"
	"net/http"
	"time"
)

// DefaultTimeout bounds a single HTTP attempt: connect, send, and read of the
// response body.
//
// It is a per-attempt budget, not a budget for the whole retry loop. A call
// with retryCount N can still take up to (N+1) * DefaultTimeout plus backoff,
// so pass a context if you need to bound the total.
const DefaultTimeout = 30 * time.Second

// httpClient is shared by every request this package makes. A nil Transport
// means http.DefaultTransport, so connections stay pooled process-wide.
var httpClient = &http.Client{Timeout: DefaultTimeout}

// SetHTTPClient replaces the client used for all requests, for callers who need
// their own transport, proxy, or timeout. Passing nil restores the default.
//
// This is not safe to call concurrently with in-flight requests; set it once
// during startup.
func SetHTTPClient(client *http.Client) {
	if client == nil {
		httpClient = &http.Client{Timeout: DefaultTimeout}
		return
	}
	httpClient = client
}

// MaxBackoff caps the wait between retries. Without a cap, a call with the
// maximum retryCount of 10 would end with a single sleep of over eight minutes.
const MaxBackoff = 8 * time.Second

// backoffDuration returns how long to wait before the next attempt, where i is
// the zero-based index of the attempt that just failed: 500ms, 1s, 2s, 4s, 8s,
// then held at MaxBackoff.
//
// The previous formula was SLEEP_TIME * i, which is linear rather than
// exponential and — worse against a rate limiter — returns zero for i == 0, so
// the first retry of a 429 went straight back out with no delay at all.
func backoffDuration(i int) time.Duration {
	if i < 0 {
		i = 0
	}

	d := SLEEP_TIME << i
	if d > MaxBackoff || d <= 0 { // d <= 0 guards against the shift overflowing
		return MaxBackoff
	}
	return d
}

// sleepCtx waits for d, returning early with ctx.Err() if the context is
// cancelled first. A retry loop that sleeps with time.Sleep can't be
// cancelled, which is most of the value of accepting a context at all.
func sleepCtx(ctx context.Context, d time.Duration) error {
	if d <= 0 {
		return nil
	}
	timer := time.NewTimer(d)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
