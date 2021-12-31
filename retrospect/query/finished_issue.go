package query

import (
	"time"

	graphql "github.com/cli/shurcooL-graphql"
)

type ClosedIssues struct {
	User struct {
		Issues struct {
			Nodes []Issue
		} `graphql:"issues(first: $limit, filterBy: {states: CLOSED, since: $from, assignee: $userName}, orderBy: {field:UPDATED_AT, direction:ASC})"`
	} `graphql:"user(login: $userName)"`
}

func (c *Client) ClosedIssues(userName string) ([]Issue, error) {
	variables := map[string]interface{}{
		"limit":    graphql.Int(c.Limit),
		"userName": graphql.String(userName),
		"from":     graphql.String(c.From.Format(time.RFC3339)),
	}
	var query ClosedIssues
	if err := c.GQL.Query("ClosedIssues", &query, variables); err != nil {
		return nil, err
	}
	return query.User.Issues.Nodes, nil
}
