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
	var cursor *graphql.String
	limit := c.Limit
	for {
		var query ReportedIssues
		variables := map[string]interface{}{
			"limit":    graphql.Int(limit),
			"userName": graphql.String(userName),
			"from":     graphql.String(c.From.Format(time.RFC3339)),
			"after":    cursor,
		}
		if cursor != nil {
			variables["after"] = graphql.NewString(*cursor)
		}
		if err := c.GQL.Query("ReportedIssues", &query, variables); err != nil {
			return nil, err
		}
		for _, issue := range query.User.Issues.Nodes {
			if issue.CreatedAt.Before(c.From) {
				continue
			}
			issues = append(issues, issue)
		}
		limit -= len(issues)
		pageInfo := query.User.Issues.PageInfo
		if !pageInfo.HasNextPage || limit <= 0 {
			break
		}
		endCursor := graphql.NewString(graphql.String(pageInfo.EndCursor))
		cursor = endCursor
	}
	return issues, nil
}
