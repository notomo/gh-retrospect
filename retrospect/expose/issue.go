package expose

import (
	"github.com/notomo/gh-retrospect/retrospect/expose/model"
	"github.com/notomo/gh-retrospect/retrospect/query"
)

func Issues(src []query.Issue) []model.Issue {
	dest := make([]model.Issue, len(src))
	for i, s := range src {
		dest[i] = Issue(s)
	}
	return dest
}

func Issue(src query.Issue) model.Issue {
	return model.Issue{
		IssuePrimitive: src.IssuePrimitive,
		LabelNames:     LabelNames(src.Labels),
	}
}

func LabelNames(src query.Labels) []string {
	dest := make([]string, len(src.Nodes))
	for i, s := range src.Nodes {
		dest[i] = s.Name
	}
	return dest
}
