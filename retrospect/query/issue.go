package query

import (
	"github.com/notomo/gh-retrospect/retrospect/expose/model"
)

type Issue struct {
	model.IssuePrimitive
	Labels Labels `graphql:"labels(first: 10)"`
}

type Labels struct {
	Nodes []struct {
		Name string `json:"name"`
	}
}
