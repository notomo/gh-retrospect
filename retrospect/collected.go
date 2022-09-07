package retrospect

import (
	"fmt"
	"time"

	"github.com/notomo/gh-retrospect/retrospect/expose"
	"github.com/notomo/gh-retrospect/retrospect/expose/model"
	"github.com/notomo/gh-retrospect/retrospect/query"
)

type Collected struct {
	From                 time.Time           `json:"from"`
	ClosedIssues         []model.Issue       `json:"closed_issues"`
	ReportedIssues       []model.Issue       `json:"reported_issues"`
	MergedPullRequests   []model.PullRequest `json:"merged_pull_requests"`
	ReviewedPullRequests []model.PullRequest `json:"reviewed_pull_requests"`
}

func Collect(
	client *query.Client,
	userName string,
	from time.Time,
	limit int,
) (*Collected, error) {
	name, err := client.ViewerName(userName)
	if err != nil {
		return nil, fmt.Errorf("get viewer name: %w", err)
	}

	var closedIssues []model.Issue
	{
		res, err := client.ClosedIssues(name, from, limit)
		if err != nil {
			return nil, fmt.Errorf("collect closed issues: %w", err)
		}
		closedIssues = expose.Issues(res)
	}

	var reportedIssues []model.Issue
	{
		res, err := client.ReportedIssues(name, from, limit)
		if err != nil {
			return nil, fmt.Errorf("collect reported issues: %w", err)
		}
		reportedIssues = expose.Issues(res)
	}

	var mergedPullRequests []model.PullRequest
	{
		res, err := client.MergedPullRequests(name, from, limit)
		if err != nil {
			return nil, fmt.Errorf("collect merged pull requests: %w", err)
		}
		mergedPullRequests = expose.PullRequests(res)
	}

	var reviewedPullRequests []model.PullRequest
	{
		res, err := client.ReviewedPullRequests(name, from, limit)
		if err != nil {
			return nil, fmt.Errorf("collect reviewed pull requests: %w", err)
		}
		reviewedPullRequests = expose.PullRequests(res)
	}

	return &Collected{
		From:                 from,
		ClosedIssues:         closedIssues,
		ReportedIssues:       reportedIssues,
		MergedPullRequests:   mergedPullRequests,
		ReviewedPullRequests: reviewedPullRequests,
	}, nil
}
