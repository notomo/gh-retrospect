package query

import (
	"time"

	graphql "github.com/cli/shurcooL-graphql"
)

type ClosedIssues struct {
	User struct {
		Issues struct {
			Nodes    []Issue
			PageInfo PageInfo
		} `graphql:"issues(first: $limit, after: $after, filterBy: {states: CLOSED, since: $from, assignee: $userName}, orderBy: {field:UPDATED_AT, direction:ASC})"`
	} `graphql:"user(login: $userName)"`
}

func (c *Client) ClosedIssues(userName string) ([]Issue, error) {
	var query ClosedIssues
	issues := []Issue{}
	if err := c.Paginate(
		"ClosedIssues",
		&query,
		map[string]interface{}{
			"userName": graphql.String(userName),
			"from":     graphql.String(c.From.Format(time.RFC3339)),
		},
		func() (PageInfo, int) {
			for _, issue := range query.User.Issues.Nodes {
				if issue.ClosedAt.Before(c.From) {
					continue
				}
				issues = append(issues, issue)
			}
			return query.User.Issues.PageInfo, len(issues)
		},
	); err != nil {
		return nil, err
	}
	return issues, nil
}
