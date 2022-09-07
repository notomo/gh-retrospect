package query

import (
	"fmt"
	"time"

	graphql "github.com/cli/shurcooL-graphql"
)

type MergedPullRequests struct {
	Search struct {
		Nodes []struct {
			PullRequest `graphql:"... on PullRequest"`
		}
		PageInfo PageInfo
	} `graphql:"search(query: $searchQuery, type: ISSUE, first: $limit, after: $after)"`
}

func (c *Client) MergedPullRequests(
	userName string,
	from time.Time,
	limit int,
) ([]PullRequest, error) {
	searchQuery := fmt.Sprintf("author:%s is:pr is:merged sort:created-asc created:>=%s", userName, from.Format(TimeFormat))

	pullRequests := []PullRequest{}
	var query MergedPullRequests
	if err := c.Paginate(
		"MergedPullRequests",
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
