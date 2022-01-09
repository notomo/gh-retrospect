package query

import "time"

type Issue struct {
	Title     string     `json:"title"`
	URL       string     `json:"url"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	ClosedAt  *time.Time `json:"closed_at,omitempty"`
}
