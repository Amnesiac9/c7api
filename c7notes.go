package c7api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// Note is a Commerce7 note, flag or task (see NoteType*) on an order or a customer.
// Set OrderId or CustomerId. DueDate, AssignedToAccountId and IsComplete are for tasks only.
type Note struct {
	Id         string    `json:"id,omitempty"`
	Type       string    `json:"type"`
	Content    string    `json:"content"`
	OrderId    string    `json:"orderId,omitempty"`
	CustomerId string    `json:"customerId,omitempty"`
	NoteDate   time.Time `json:"noteDate"`

	DueDate             *time.Time `json:"dueDate,omitempty"`
	AssignedToAccountId string     `json:"assignedToAccountId,omitempty"`
	IsComplete          *bool      `json:"isComplete,omitempty"`
}

// CreateNoteContext posts a note and returns it with the id Commerce7 assigned.
// A zero NoteDate is set to now.
func CreateNoteContext(ctx context.Context, note Note, tenant, auth string, retryCount int, rl genericRateLimiter) (*Note, error) {
	if note.OrderId == "" && note.CustomerId == "" {
		return nil, errors.New("a note needs an order id or a customer id")
	}
	if note.Type == "" {
		note.Type = NoteTypeNote
	}
	if note.NoteDate.IsZero() {
		note.NoteDate = time.Now().UTC()
	}
	note.Id = ""

	body, err := PostContext(ctx, &note, Endpoints.Note, tenant, auth, retryCount, rl)
	if err != nil {
		return nil, fmt.Errorf("creating note: %w", err)
	}
	var created Note
	if err := json.Unmarshal(*body, &created); err != nil {
		return nil, fmt.Errorf("decoding note: %w", err)
	}
	return &created, nil
}
