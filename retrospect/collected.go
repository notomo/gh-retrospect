package retrospect

import (
	"fmt"
	"time"

	"github.com/notomo/gh-retrospect/retrospect/query"
)

type Collected struct {
	From           time.Time     `json:"from"`
	ClosedIssues   []query.Issue `json:"closedIssues"`
	ReportedIssues []query.Issue `json:"reportedIssues"`
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
		collected.ClosedIssues = res
	}
	{
		res, err := client.ReportedIssues(name, from, limit)
		if err != nil {
			return nil, fmt.Errorf("collect reported issues: %w", err)
		}
		collected.ReportedIssues = res
	}

	return &collected, nil
}
