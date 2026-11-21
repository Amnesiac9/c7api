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

type recordedRequest struct {
	method string
	path   string
	query  string
	body   []byte
}

// fakeC7 answers every request with the handler's status and body, recording each request.
func fakeC7(t *testing.T, handler func(r *http.Request) (int, string)) *[]recordedRequest {
	t.Helper()
	var reqs []recordedRequest
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		reqs = append(reqs, recordedRequest{r.Method, r.URL.Path, r.URL.RawQuery, body})
		status, resp := handler(r)
		w.WriteHeader(status)
		w.Write([]byte(resp))
	}))
	real := Endpoints
	Endpoints = GetEndpoints(srv.URL)
	t.Cleanup(func() {
		srv.Close()
		Endpoints = real
	})
	return &reqs
}

func TestListWebHooks(t *testing.T) {
	reqs := fakeC7(t, func(r *http.Request) (int, string) {
		return http.StatusOK, `{"webHooks":[
			{"id":"h1","object":"Order","action":"Create","url":"https://a.example/hook","username":"u","appId":"my-app","status":"Enabled"},
			{"id":"h2","object":"Customer","action":"Update","url":"https://b.example/hook","appId":"other-app","status":"Disabled"}
		],"total":2}`
	})

	hooks, err := ListWebHooksContext(context.Background(), "my-tenant", "Basic abc", 0, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(*reqs) != 1 || (*reqs)[0].path != "/web-hook" {
		t.Fatalf("want one GET /web-hook, got %+v", *reqs)
	}
	if len(hooks) != 2 || hooks[0].AppId != "my-app" || hooks[1].Object != "Customer" {
		t.Fatalf("unexpected hooks: %+v", hooks)
	}
	if !hooks[0].IsEnabled() || hooks[1].IsEnabled() {
		t.Errorf("enabled = %v, %v; want true, false", hooks[0].IsEnabled(), hooks[1].IsEnabled())
	}
}

func TestCreateWebHook(t *testing.T) {
	reqs := fakeC7(t, func(r *http.Request) (int, string) {
		return http.StatusOK, `{"id":"new-id","object":"Order","action":"Create","url":"https://a.example/hook","username":"u","status":"Enabled"}`
	})

	hook, err := CreateWebHookContext(context.Background(), WebHook{
		Id: "ignored", Object: WebhookObjectOrder, Action: WebhookActionCreate,
		URL: "https://a.example/hook", Username: "u", Password: "p",
	}, "my-tenant", "Basic abc", 0, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if hook.Id != "new-id" {
		t.Errorf("id = %q, want new-id", hook.Id)
	}

	req := (*reqs)[0]
	if req.method != http.MethodPost || req.path != "/web-hook" {
		t.Fatalf("got %s %s, want POST /web-hook", req.method, req.path)
	}
	var sent map[string]any
	json.Unmarshal(req.body, &sent)
	if _, hasId := sent["id"]; hasId {
		t.Errorf("create should not send an id: %s", req.body)
	}
	if sent["password"] != "p" || sent["object"] != "Order" || sent["action"] != "Create" {
		t.Errorf("unexpected body: %s", req.body)
	}
}

func TestCreateWebHook_NoIdReturned(t *testing.T) {
	fakeC7(t, func(r *http.Request) (int, string) { return http.StatusOK, `{}` })

	_, err := CreateWebHookContext(context.Background(), WebHook{Object: "Order", Action: "Create", URL: "u"}, "t", "Basic abc", 0, nil)
	if err == nil || !strings.Contains(err.Error(), "without an id") {
		t.Fatalf("want a missing id error, got %v", err)
	}
}

func TestUpdateWebHook(t *testing.T) {
	reqs := fakeC7(t, func(r *http.Request) (int, string) {
		return http.StatusOK, `{"id":"h1","object":"Order","action":"Update","url":"https://a.example/hook"}`
	})

	_, err := UpdateWebHookContext(context.Background(), WebHook{Id: "h1", Object: "Order", Action: "Update", URL: "https://a.example/hook"}, "t", "Basic abc", 0, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req := (*reqs)[0]; req.method != http.MethodPut || req.path != "/web-hook/h1" {
		t.Fatalf("got %s %s, want PUT /web-hook/h1", req.method, req.path)
	}
}

func TestDeleteWebHook(t *testing.T) {
	reqs := fakeC7(t, func(r *http.Request) (int, string) { return http.StatusOK, `{}` })

	if err := DeleteWebHookContext(context.Background(), "h1", "t", "Basic abc", 0, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req := (*reqs)[0]; req.method != http.MethodDelete || req.path != "/web-hook/h1" {
		t.Fatalf("got %s %s, want DELETE /web-hook/h1", req.method, req.path)
	}
}

func TestDeleteWebHook_NotFoundIsSuccess(t *testing.T) {
	fakeC7(t, func(r *http.Request) (int, string) {
		return http.StatusNotFound, `{"statusCode":404,"type":"notFound","message":"Not found"}`
	})

	if err := DeleteWebHookContext(context.Background(), "gone", "t", "Basic abc", 0, nil); err != nil {
		t.Fatalf("a 404 should count as deleted, got %v", err)
	}
}

func TestDeleteWebHook_MissingId(t *testing.T) {
	if err := DeleteWebHookContext(context.Background(), "", "t", "Basic abc", 0, nil); err == nil {
		t.Fatal("want an error for a missing id")
	}
}

func TestCreateNote_Task(t *testing.T) {
	reqs := fakeC7(t, func(r *http.Request) (int, string) {
		return http.StatusOK, `{"id":"n1","type":"Task","content":"Check this order","orderId":"o1"}`
	})

	incomplete := false
	note, err := CreateNoteContext(context.Background(), Note{
		Type: NoteTypeTask, Content: "Check this order", OrderId: "o1", IsComplete: &incomplete,
	}, "t", "Basic abc", 0, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if note.Id != "n1" {
		t.Errorf("id = %q, want n1", note.Id)
	}

	req := (*reqs)[0]
	if req.method != http.MethodPost || req.path != "/note" {
		t.Fatalf("got %s %s, want POST /note", req.method, req.path)
	}
	var sent map[string]any
	json.Unmarshal(req.body, &sent)
	if sent["type"] != "Task" || sent["orderId"] != "o1" || sent["isComplete"] != false {
		t.Errorf("unexpected body: %s", req.body)
	}
	if sent["noteDate"] == "0001-01-01T00:00:00Z" {
		t.Error("a zero note date should be set to now")
	}
	if _, has := sent["customerId"]; has {
		t.Errorf("empty customerId should be left out: %s", req.body)
	}
}

func TestCreateNote_NeedsTarget(t *testing.T) {
	if _, err := CreateNoteContext(context.Background(), Note{Content: "x"}, "t", "Basic abc", 0, nil); err == nil {
		t.Fatal("want an error when neither order nor customer is set")
	}
}
