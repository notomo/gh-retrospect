package query

import (
	"time"
)

// TODO: improve performance

type MergedPullRequest struct {
	User struct {
		PullRequests struct {
			Nodes    []PullRequest
			PageInfo PageInfo
		} `graphql:"pullRequests(first: $limit, after: $after, orderBy: {field:CREATED_AT, direction:ASC}, states: [MERGED])"`
	} `graphql:"user(login: $userName)"`
}

func (c *Client) MergedPullRequests(
	userName string,
	from time.Time,
	limit int,
) ([]PullRequest, error) {
	pullRequests := []PullRequest{}

	var query MergedPullRequest
	if err := c.Paginate(
		"MergedPullRequests",
		&query,
		NewParameter(
			WithUserName(userName),
		),
		func() (PageInfo, int) {
			for _, pullRequest := range query.User.PullRequests.Nodes {
				if pullRequest.CreatedAt.Before(from) {
					continue
				}
				pullRequests = append(pullRequests, pullRequest)
			}
			return query.User.PullRequests.PageInfo, len(pullRequests)
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
