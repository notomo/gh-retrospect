package model

import "time"

type IssuePrimitive struct {
	Title     string     `json:"title"`
	URL       string     `json:"url"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	ClosedAt  *time.Time `json:"closed_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

type Issue struct {
	IssuePrimitive
	LabelNames []string `json:"label_names"`
}
