package c7api

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
)

// requestWithRetryAndRead is the shared implementation behind
// RequestWithRetryAndRead (v1) and RequestWithRetryAndReadV2 (v2).
//
// extraHeaders is applied after the standard tenant/content-type/auth headers,
// so a caller can add to them or override them. The v2 API uses this for the
// two headers it requires on top of the v1 set.
func requestWithRetryAndRead(ctx context.Context, method string, url string, queries map[string]string, reqBody *[]byte, tenant string, c7AppAuthEncoded string, retryCount int, rl genericRateLimiter, extraHeaders map[string]string) (*[]byte, error) {
	//
	if url == "" || tenant == "" || c7AppAuthEncoded == "" {
		return nil, fmt.Errorf("error getting JSON from C7: nil or blank value in arguments")
	}

	if ctx == nil {
		ctx = context.Background()
	}

	if reqBody == nil {
		reqBody = &[]byte{}
	}

	minRetryCount := 0
	maxRetryCount := 10

	if retryCount < minRetryCount {
		retryCount = minRetryCount
	} else if retryCount > maxRetryCount {
		retryCount = maxRetryCount
	}

	response := &http.Response{StatusCode: 0}
	body := []byte{}

	for i := 0; i <= retryCount; i++ {
		// The rate limiter and the sleeps below can hold us for a while, so
		// check for cancellation before spending another attempt.
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		if rl != nil && !reflect.ValueOf(rl).IsNil() {
			rl.Wait()
		}

		req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(*reqBody))
		if err != nil {
			return nil, fmt.Errorf("error creating GET request for C7: %v", err)
		}

		if queries != nil {
			query := req.URL.Query()
			for k, v := range queries {
				query.Add(k, v)
			}
			req.URL.RawQuery = query.Encode()
		}

		req.Header.Set("tenant", tenant)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("Authorization", c7AppAuthEncoded)

		for k, v := range extraHeaders {
			req.Header.Set(k, v)
		}

		response, err = httpClient.Do(req)
		if err != nil {
			// A cancelled context surfaces here as an opaque *url.Error, so
			// report ctx.Err() to keep errors.Is usable by the caller.
			if ctxErr := ctx.Err(); ctxErr != nil {
				return nil, ctxErr
			}
			return nil, fmt.Errorf("error making GET request to C7: %v", err)
		}

		body, err = io.ReadAll(response.Body)
		response.Body.Close()
		if err != nil {
			if ctxErr := ctx.Err(); ctxErr != nil {
				return nil, ctxErr
			}
			return nil, fmt.Errorf("error reading response body from C7: %v", err)
		}

		// 200-299 is success, return body and nil error
		if ResponseIsOK(response.StatusCode) {
			return &body, nil
		}

		// A missing resource won't appear on a retry, so fail fast.
		if response.StatusCode == http.StatusNotFound {
			break
		}

		// Exponential backoff before the next attempt. Skipped on the final
		// pass, where there is no next attempt to wait for.
		if i < retryCount {
			if err := sleepCtx(ctx, backoffDuration(i)); err != nil {
				return nil, err
			}
		}
	}

	// Read the C7 Error if present
	// Always return as C7Error after this point, since this means C7 sent an error message.
	// If we have trouble reading it for some reason, handle that here.
	c7Error := C7Error{}
	err := c7Error.UnmarshalJSON(body)
	if err != nil {
		c7Error.StatusCode = response.StatusCode
		c7Error.Err = errors.New("error unmarshalling Commerce7 Error Message: " + err.Error() + "json: " + string(body))
		return &body, &c7Error
	}

	// Fall back to the HTTP status when the body didn't carry one,
	// so callers can still switch on things like 404.
	if c7Error.StatusCode == 0 {
		c7Error.StatusCode = response.StatusCode
	}

	// Add the raw json body to the err as well in case needed.
	// TODO: Handle this better to allow slog nested json...
	c7Error.Err = errors.New(string(body))
	return &body, &c7Error
}
