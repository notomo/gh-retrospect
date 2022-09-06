package model

import "time"

type PullRequestPrimitive struct {
	Title     string     `json:"title"`
	URL       string     `json:"url"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	ClosedAt  *time.Time `json:"closed_at,omitempty"`
}

type PullRequest struct {
	PullRequestPrimitive
	LabelNames []string `json:"label_names"`
}
