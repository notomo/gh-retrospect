package query

import (
	"fmt"
	"time"

	graphql "github.com/cli/shurcooL-graphql"
)

type ClosedIssues struct {
	Search struct {
		Nodes []struct {
			Issue `graphql:"... on Issue"`
		}
		PageInfo PageInfo
	} `graphql:"search(query: $searchQuery, type: ISSUE, first: $limit, after: $after)"`
}

func (c *Client) ClosedIssues(
	userName string,
	from time.Time,
	to time.Time,
	limit int,
) ([]Issue, error) {
	searchQuery := fmt.Sprintf("assignee:%s is:issue is:closed sort:updated-asc closed:>=%s", userName, from.Format(TimeFormat))
	if !to.IsZero() {
		searchQuery += fmt.Sprintf(" closed:<=%s", to.Format(TimeFormat))
	}

	issues := []Issue{}

	var query ClosedIssues
	if err := c.Paginate(
		"ClosedIssues",
		&query,
		NewParameter(
			func(p Parameter) {
				p["searchQuery"] = graphql.String(searchQuery)
			},
		),
		func() (PageInfo, int) {
			for _, node := range query.Search.Nodes {
				issue := node.Issue

				if issue.ClosedAt.Before(from) {
					continue
				}
				if !to.IsZero() && issue.ClosedAt.After(to) {
					continue
				}

				if !to.IsZero() && issue.UpdatedAt != nil && issue.UpdatedAt.After(to) {
					return query.Search.PageInfo, limit
				}

				issues = append(issues, issue)
			}
			return query.Search.PageInfo, len(issues)
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
