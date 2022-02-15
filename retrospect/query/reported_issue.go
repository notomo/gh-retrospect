package query

import (
	"time"
)

type ReportedIssues struct {
	User struct {
		Issues struct {
			Nodes    []Issue
			PageInfo PageInfo
		} `graphql:"issues(first: $limit, after: $after, filterBy: {createdBy: $userName, since: $from}, orderBy: {field:CREATED_AT, direction:ASC})"`
	} `graphql:"user(login: $userName)"`
}

func (c *Client) ReportedIssues(
	userName string,
	from time.Time,
	limit int,
) ([]Issue, error) {
	issues := []Issue{}

	var query ReportedIssues
	if err := c.Paginate(
		"ReportedIssues",
		&query,
		NewParameter(
			WithUserName(userName),
			WithFrom(from),
		),
		func() (PageInfo, int) {
			for _, issue := range query.User.Issues.Nodes {
				if issue.CreatedAt.Before(from) {
					continue
				}
				issues = append(issues, issue)
			}
			return query.User.Issues.PageInfo, len(issues)
		},
		limit,
	); err != nil {
		return nil, err
	}

	if len(issues) > limit {
		return issues[:limit], nil
	}
	return issues, nil
}
