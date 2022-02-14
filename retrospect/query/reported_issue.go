package query

import (
	"time"

	graphql "github.com/cli/shurcooL-graphql"
)

type ReportedIssues struct {
	User struct {
		Issues struct {
			Nodes    []Issue
			PageInfo PageInfo
		} `graphql:"issues(first: $limit, after: $after, filterBy: {createdBy: $userName, since: $from}, orderBy: {field:CREATED_AT, direction:ASC})"`
	} `graphql:"user(login: $userName)"`
}

func (c *Client) ReportedIssues(userName string) ([]Issue, error) {
	issues := []Issue{}

	var query ReportedIssues
	if err := c.Paginate(
		"ReportedIssues",
		&query,
		map[string]interface{}{
			"userName": graphql.String(userName),
			"from":     graphql.String(c.From.Format(time.RFC3339)),
		},
		func() (PageInfo, int) {
			for _, issue := range query.User.Issues.Nodes {
				if issue.CreatedAt.Before(c.From) {
					continue
				}
				issues = append(issues, issue)
			}
			return query.User.Issues.PageInfo, len(issues)
		},
	); err != nil {
		return nil, err
	}

	if len(issues) > c.Limit {
		return issues[:c.Limit], nil
	}
	return issues, nil
}
