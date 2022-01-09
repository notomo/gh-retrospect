package query

import graphql "github.com/cli/shurcooL-graphql"

type PageInfo struct {
	EndCursor   string `json:"end_cursor"`
	HasNextPage bool   `json:"has_next_page"`
}

func (c *Client) Paginate(
	name string,
	query interface{},
	variables map[string]interface{},
	each func() (PageInfo, int),
) error {
	var cursor *graphql.String
	limit := c.Limit
	for {
		vars := map[string]interface{}{
			"limit": graphql.Int(limit),
			"after": cursor,
		}
		if cursor != nil {
			vars["after"] = graphql.NewString(*cursor)
		}
		for k, v := range variables {
			vars[k] = v
		}

		if err := c.GQL.Query(name, query, vars); err != nil {
			return err
		}

		pageInfo, count := each()
		limit = c.Limit - count
		if !pageInfo.HasNextPage || limit <= 0 {
			break
		}
		endCursor := graphql.NewString(graphql.String(pageInfo.EndCursor))
		cursor = endCursor
	}
	return nil
}
