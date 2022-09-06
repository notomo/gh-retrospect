package query

import (
	"github.com/notomo/gh-retrospect/retrospect/expose/model"
)

type PullRequest struct {
	model.PullRequestPrimitive
	Labels Labels `graphql:"labels(first: 10)"`
}
