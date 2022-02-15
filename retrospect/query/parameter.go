package query

import (
	"time"

	graphql "github.com/cli/shurcooL-graphql"
)

type Parameter map[string]interface{}

func NewParameter(opts ...func(Parameter)) Parameter {
	param := map[string]interface{}{}
	for _, opt := range opts {
		opt(param)
	}
	return param
}

func WithFrom(from time.Time) func(Parameter) {
	return func(p Parameter) {
		p["from"] = graphql.String(from.Format(time.RFC3339))
	}
}

func WithUserName(userName string) func(Parameter) {
	return func(p Parameter) {
		p["userName"] = graphql.String(userName)
	}
}
