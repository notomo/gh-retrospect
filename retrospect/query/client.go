package query

import (
	"github.com/cli/go-gh/v2/pkg/api"
)

type Client struct {
	GQL *api.GraphQLClient
}

func NewClient(
	gql *api.GraphQLClient,
) *Client {
	return &Client{
		GQL: gql,
	}
}
