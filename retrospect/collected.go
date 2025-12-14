package retrospect

import (
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/notomo/gh-retrospect/retrospect/expose"
	"github.com/notomo/gh-retrospect/retrospect/expose/model"
	"github.com/notomo/gh-retrospect/retrospect/query"
)

type Collected struct {
	From                 time.Time           `json:"from"`
	To                   time.Time           `json:"to"`
	ClosedIssues         []model.Issue       `json:"closed_issues"`
	ReportedIssues       []model.Issue       `json:"reported_issues"`
	MergedPullRequests   []model.PullRequest `json:"merged_pull_requests"`
	ReviewedPullRequests []model.PullRequest `json:"reviewed_pull_requests"`
}

func Collect(
	client *query.Client,
	userName string,
	from time.Time,
	to time.Time,
	limit int,
) (*Collected, error) {
	name, err := client.ViewerName(userName)
	if err != nil {
		return nil, fmt.Errorf("get viewer name: %w", err)
	}

	var closedIssues []model.Issue
	var reportedIssues []model.Issue
	var mergedPullRequests []model.PullRequest
	var reviewedPullRequests []model.PullRequest

	g := new(errgroup.Group)

	g.Go(func() error {
		res, err := client.ClosedIssues(name, from, to, limit)
		if err != nil {
			return fmt.Errorf("collect closed issues: %w", err)
		}
		closedIssues = expose.Issues(res)
		return nil
	})

	g.Go(func() error {
		res, err := client.ReportedIssues(name, from, to, limit)
		if err != nil {
			return fmt.Errorf("collect reported issues: %w", err)
		}
		reportedIssues = expose.Issues(res)
		return nil
	})

	g.Go(func() error {
		res, err := client.MergedPullRequests(name, from, to, limit)
		if err != nil {
			return fmt.Errorf("collect merged pull requests: %w", err)
		}
		mergedPullRequests = expose.PullRequests(res)
		return nil
	})

	g.Go(func() error {
		res, err := client.ReviewedPullRequests(name, from, to, limit)
		if err != nil {
			return fmt.Errorf("collect reviewed pull requests: %w", err)
		}
		reviewedPullRequests = expose.PullRequests(res)
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return &Collected{
		From:                 from,
		To:                   to,
		ClosedIssues:         closedIssues,
		ReportedIssues:       reportedIssues,
		MergedPullRequests:   mergedPullRequests,
		ReviewedPullRequests: reviewedPullRequests,
	}, nil
}
