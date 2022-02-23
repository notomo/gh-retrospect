package retrospect

import (
	"fmt"
	"time"

	"github.com/notomo/gh-retrospect/retrospect/expose"
	"github.com/notomo/gh-retrospect/retrospect/expose/model"
	"github.com/notomo/gh-retrospect/retrospect/query"
)

type Collected struct {
	From           time.Time     `json:"from"`
	ClosedIssues   []model.Issue `json:"closed_issues"`
	ReportedIssues []model.Issue `json:"reported_issues"`
}

func Collect(
	client *query.Client,
	userName string,
	from time.Time,
	limit int,
) (*Collected, error) {
	collected := Collected{From: from}

	name, err := client.ViewerName(userName)
	if err != nil {
		return nil, fmt.Errorf("get viewr name: %w", err)
	}

	{
		res, err := client.ClosedIssues(name, from, limit)
		if err != nil {
			return nil, fmt.Errorf("collect closed issues: %w", err)
		}
		collected.ClosedIssues = expose.Issues(res)
	}
	{
		res, err := client.ReportedIssues(name, from, limit)
		if err != nil {
			return nil, fmt.Errorf("collect reported issues: %w", err)
		}
		collected.ReportedIssues = expose.Issues(res)
	}

	return &collected, nil
}
