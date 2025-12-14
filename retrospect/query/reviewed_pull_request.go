package query

import (
	"fmt"
	"time"

	graphql "github.com/cli/shurcooL-graphql"
)

type ReviewedPullRequests struct {
	Search struct {
		Nodes []struct {
			PullRequest `graphql:"... on PullRequest"`
		}
		PageInfo PageInfo
	} `graphql:"search(query: $searchQuery, type: ISSUE, first: $limit, after: $after)"`
}

func (c *Client) ReviewedPullRequests(
	userName string,
	from time.Time,
	to time.Time,
	limit int,
) ([]PullRequest, error) {
	searchQuery := fmt.Sprintf("reviewed-by:%s is:pr sort:created-asc created:>=%s", userName, from.Format(TimeFormat))
	if !to.IsZero() {
		searchQuery += fmt.Sprintf(" created:<=%s", to.Format(TimeFormat))
	}

	pullRequests := []PullRequest{}
	var query MergedPullRequests
	if err := c.Paginate(
		"ReviewedPullRequests",
		&query,
		NewParameter(
			func(p Parameter) {
				p["searchQuery"] = graphql.String(searchQuery)
			},
		),
		func() (PageInfo, int) {
			for _, node := range query.Search.Nodes {
				pullRequest := node.PullRequest
				if pullRequest.CreatedAt.Before(from) {
					continue
				}
				if !to.IsZero() && pullRequest.CreatedAt.After(to) {
					// Since results are sorted by created-asc, stop pagination
					return query.Search.PageInfo, limit
				}
				pullRequests = append(pullRequests, pullRequest)
			}
			return query.Search.PageInfo, len(pullRequests)
		},
		limit,
	); err != nil {
		return nil, err
	}

	if len(pullRequests) > limit {
		return pullRequests[:limit], nil
	}
	return pullRequests, nil
}
