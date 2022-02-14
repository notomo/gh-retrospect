package query

import graphql "github.com/cli/shurcooL-graphql"

type PageInfo struct {
	EndCursor   string `json:"end_cursor"`
	HasNextPage bool   `json:"has_next_page"`
}

const limitPerRequest = 100

func (c *Client) Paginate(
	name string,
	query interface{},
	variables map[string]interface{},
	each func() (PageInfo, int),
) error {
	var cursor *graphql.String
	currentCount := 0
	for {
		limit := c.Limit - currentCount
		if limit <= 0 {
			break
		}

		vars := map[string]interface{}{
			"limit": graphql.Int(limitPerRequest),
			"after": cursor,
		}
		for k, v := range variables {
			vars[k] = v
		}
		if err := c.GQL.Query(name, query, vars); err != nil {
			return err
		}

		pageInfo, count := each()
		if !pageInfo.HasNextPage {
			break
		}
		endCursor := graphql.NewString(graphql.String(pageInfo.EndCursor))
		cursor = endCursor
		currentCount = count
	}
	return nil
}
