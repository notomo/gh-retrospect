package expose

import (
	"github.com/notomo/gh-retrospect/retrospect/expose/model"
	"github.com/notomo/gh-retrospect/retrospect/query"
)

func PullRequests(src []query.PullRequest) []model.PullRequest {
	dest := make([]model.PullRequest, len(src))
	for i, s := range src {
		dest[i] = PullRequest(s)
	}
	return dest
}

func PullRequest(src query.PullRequest) model.PullRequest {
	return model.PullRequest{
		PullRequestPrimitive: src.PullRequestPrimitive,
		LabelNames:           LabelNames(src.Labels),
	}
}
