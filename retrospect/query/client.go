package query

import (
	"time"

	"github.com/cli/go-gh/pkg/api"
)

type Client struct {
	GQL   api.GQLClient
	From  time.Time
	Limit int
}

func NewClient(
	gql api.GQLClient,
	from time.Time,
	limit int,
) *Client {
	return &Client{
		GQL:   gql,
		From:  from,
		Limit: limit,
	}
}
