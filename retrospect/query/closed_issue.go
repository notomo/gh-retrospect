package query

import (
	"time"
)

type ClosedIssues struct {
	User struct {
		Issues struct {
			Nodes    []Issue
			PageInfo PageInfo
		} `graphql:"issues(first: $limit, after: $after, filterBy: {states: CLOSED, since: $from, assignee: $userName}, orderBy: {field:UPDATED_AT, direction:ASC})"`
	} `graphql:"user(login: $userName)"`
}

func (c *Client) ClosedIssues(
	userName string,
	from time.Time,
	to time.Time,
	limit int,
) ([]Issue, error) {
	issues := []Issue{}

	var query ClosedIssues
	if err := c.Paginate(
		"ClosedIssues",
		&query,
		NewParameter(
			WithUserName(userName),
			WithFrom(from),
			WithTo(to),
		),
		func() (PageInfo, int) {
			for _, issue := range query.User.Issues.Nodes {
				if issue.ClosedAt.Before(from) {
					continue
				}
				if !to.IsZero() && issue.ClosedAt.After(to) {
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
