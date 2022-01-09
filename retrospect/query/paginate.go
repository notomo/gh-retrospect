package query

type PageInfo struct {
	EndCursor   string `json:"end_cursor"`
	HasNextPage bool   `json:"has_next_page"`
}
