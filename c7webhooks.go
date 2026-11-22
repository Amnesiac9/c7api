package c7api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// WebHook is a Commerce7 web-hook registration.
//
// Username and Password are the Basic auth Commerce7 sends with every delivery. Listing
// returns the username but never the password.
//
// Listing returns every app's webhooks on the tenant, so filter on AppId to find your own.
type WebHook struct {
	Id                string     `json:"id,omitempty"`
	Object            string     `json:"object"`
	Action            string     `json:"action"`
	URL               string     `json:"url"`
	Username          string     `json:"username,omitempty"`
	Password          string     `json:"password,omitempty"`
	ApiVersion        string     `json:"apiVersion,omitempty"`
	AppId             string     `json:"appId,omitempty"`
	AppTitle          string     `json:"appTitle,omitempty"`
	Status            string     `json:"status,omitempty"`
	NotificationEmail string     `json:"notificationEmail,omitempty"`
	CreatedAt         *time.Time `json:"createdAt,omitempty"`
	UpdatedAt         *time.Time `json:"updatedAt,omitempty"`
}

// IsEnabled reports whether the webhook still delivers. A disabled webhook must be recreated.
func (h WebHook) IsEnabled() bool {
	return h.Status == "" || h.Status == WebhookStatusEnabled
}

type WebHookList struct {
	WebHooks []WebHook `json:"webHooks"`
	Total    int       `json:"total"`
}

func (l WebHookList) GetItems() []WebHook { return l.WebHooks }
func (l WebHookList) GetTotal() int       { return l.Total }

// ListWebHooksContext returns every webhook registered on the tenant, by any app.
func ListWebHooksContext(ctx context.Context, tenant, auth string, retryCount int, rl genericRateLimiter) ([]WebHook, error) {
	hooks, err := GetAllContext[WebHook, WebHookList](ctx, Endpoints.WebHook, nil, nil, tenant, auth, retryCount, rl)
	if err != nil {
		return nil, fmt.Errorf("listing web-hooks: %w", err)
	}
	return *hooks, nil
}

func GetWebHookContext(ctx context.Context, id, tenant, auth string, retryCount int, rl genericRateLimiter) (*WebHook, error) {
	if id == "" {
		return nil, errors.New("missing web-hook id")
	}
	return GetContext[WebHook](ctx, Endpoints.WebHook+"/"+id, nil, nil, tenant, auth, retryCount, rl)
}

// CreateWebHookContext registers a webhook and returns it with the id Commerce7 assigned.
// Set Username and Password, or deliveries arrive without auth.
func CreateWebHookContext(ctx context.Context, hook WebHook, tenant, auth string, retryCount int, rl genericRateLimiter) (*WebHook, error) {
	hook.Id = ""
	body, err := PostContext(ctx, &hook, Endpoints.WebHook, tenant, auth, retryCount, rl)
	if err != nil {
		return nil, fmt.Errorf("creating web-hook: %w", err)
	}
	return decodeWebHook(body)
}

// UpdateWebHookContext replaces the webhook with hook.Id.
func UpdateWebHookContext(ctx context.Context, hook WebHook, tenant, auth string, retryCount int, rl genericRateLimiter) (*WebHook, error) {
	if hook.Id == "" {
		return nil, errors.New("missing web-hook id")
	}
	url := Endpoints.WebHook + "/" + hook.Id
	hook.Id = ""
	body, err := PutContext(ctx, &hook, url, tenant, auth, retryCount, rl)
	if err != nil {
		return nil, fmt.Errorf("updating web-hook: %w", err)
	}
	return decodeWebHook(body)
}

// DeleteWebHookContext removes a webhook. A 404 counts as success: it's already gone.
func DeleteWebHookContext(ctx context.Context, id, tenant, auth string, retryCount int, rl genericRateLimiter) error {
	if id == "" {
		return errors.New("missing web-hook id")
	}
	_, err := RequestWithRetryAndReadContext(ctx, http.MethodDelete, Endpoints.WebHook+"/"+id, nil, nil, tenant, auth, retryCount, rl)
	if err != nil && !IsNotFound(err) {
		return fmt.Errorf("deleting web-hook: %w", err)
	}
	return nil
}

func decodeWebHook(body *[]byte) (*WebHook, error) {
	var hook WebHook
	if err := json.Unmarshal(*body, &hook); err != nil {
		return nil, fmt.Errorf("decoding web-hook: %w", err)
	}
	if hook.Id == "" {
		return nil, fmt.Errorf("commerce7 returned a web-hook without an id: %s", string(*body))
	}
	return &hook, nil
}
