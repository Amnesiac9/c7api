package c7api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type titled struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

func TestListPageReadsAnyListKey(t *testing.T) {
	var page ListPage[titled]
	body := `{"promotions": [{"id": "p1", "title": "Ten off"}, {"id": "p2", "title": "Free ship"}], "total": 2, "cursor": "x"}`
	if err := json.Unmarshal([]byte(body), &page); err != nil {
		t.Fatal(err)
	}
	if page.GetTotal() != 2 || len(page.GetItems()) != 2 || page.GetItems()[1].Title != "Free ship" {
		t.Fatalf("got %+v", page)
	}
}

func TestListPageWithGetAll(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"posProfiles": [{"id": "a", "title": "Tasting Room"}], "total": 1}`)
	}))
	defer srv.Close()

	items, err := GetAll[titled, ListPage[titled]](srv.URL, nil, nil, "tenant", "auth", 0, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(*items) != 1 || (*items)[0].Title != "Tasting Room" {
		t.Fatalf("got %+v", *items)
	}
}
