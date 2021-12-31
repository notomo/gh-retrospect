package retrospect

import (
	"fmt"

	"github.com/notomo/gh-retrospect/retrospect/query"
)

type Collected struct {
	ClosedIssues   []query.Issue `json:"closedIssues"`
	ReportedIssues []query.Issue `json:"reportedIssues"`
}

func Collect(client *query.Client, userName string) (*Collected, error) {
	collected := Collected{}

	name, err := client.ViewerName(userName)
	if err != nil {
		return nil, fmt.Errorf("get viewr name: %w", err)
	}

	{
		res, err := client.ClosedIssues(name)
		if err != nil {
			return nil, fmt.Errorf("collect closed issues: %w", err)
		}
		collected.ClosedIssues = res
	}
	{
		res, err := client.ReportedIssues(name)
		if err != nil {
			return nil, fmt.Errorf("collect reported issues: %w", err)
		}
		collected.ReportedIssues = res
	}

	return &collected, nil
}
