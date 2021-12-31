package query

import (
	"time"

	graphql "github.com/cli/shurcooL-graphql"
)

type ReportedIssues struct {
	User struct {
		Issues struct {
			Nodes []Issue
		} `graphql:"issues(first: $limit, filterBy: {createdBy: $userName, since: $from}, orderBy: {field:CREATED_AT, direction:ASC})"`
	} `graphql:"user(login: $userName)"`
}

func (c *Client) ReportedIssues(userName string) ([]Issue, error) {
	variables := map[string]interface{}{
		"limit":    graphql.Int(c.Limit),
		"userName": graphql.String(userName),
		"from":     graphql.String(c.From.Format(time.RFC3339)),
	}
	var query ReportedIssues
	if err := c.GQL.Query("ReportedIssues", &query, variables); err != nil {
		return nil, err
	}
	return query.User.Issues.Nodes, nil
}
