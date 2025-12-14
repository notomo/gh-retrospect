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

type (
	DateTime string
)

const TimeFormat = time.RFC3339

func WithFrom(from time.Time) func(Parameter) {
	return func(p Parameter) {
		p["from"] = DateTime(from.Format(TimeFormat))
	}
}

func WithTo(to time.Time) func(Parameter) {
	return func(p Parameter) {
		// Note: GitHub GraphQL API doesn't support 'until' parameter in filterBy
		// We rely on client-side filtering instead
		_ = to
	}
}

func WithUserName(userName string) func(Parameter) {
	return func(p Parameter) {
		p["userName"] = graphql.String(userName)
	}
}
