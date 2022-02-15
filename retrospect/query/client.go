package query

import (
	"github.com/cli/go-gh/pkg/api"
)

type Client struct {
	GQL api.GQLClient
}

func NewClient(
	gql api.GQLClient,
) *Client {
	return &Client{
		GQL: gql,
	}
}
