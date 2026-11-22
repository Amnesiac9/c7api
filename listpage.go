package c7api

import "encoding/json"

// ListPage reads any Commerce7 list response, e.g. {"promotions": [...], "total": 3}, so a
// GetAll call doesn't need a wrapper type per object:
//
//	GetAllContext[Promotion, ListPage[Promotion]](ctx, Endpoints.Promotion, ...)
//
// The records are the response's one array; any other fields besides total are ignored.
type ListPage[T any] struct {
	Items []T
	Total int
}

func (p ListPage[T]) GetItems() []T { return p.Items }
func (p ListPage[T]) GetTotal() int { return p.Total }

func (p *ListPage[T]) UnmarshalJSON(b []byte) error {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(b, &fields); err != nil {
		return err
	}
	for key, raw := range fields {
		switch {
		case key == "total":
			if err := json.Unmarshal(raw, &p.Total); err != nil {
				return err
			}
		case len(raw) > 0 && raw[0] == '[':
			if err := json.Unmarshal(raw, &p.Items); err != nil {
				return err
			}
		}
	}
	return nil
}
