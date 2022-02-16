package main

import (
	"bytes"
	"testing"

	"github.com/notomo/gh-retrospect/retrospect/gqltest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	gql, err := gqltest.New(
		t,
		gqltest.WithOK("ViewerName", `{
  "data": {
    "viewer": {
      "login": "notomo"
    }
  }
}`),
		gqltest.WithOK("ClosedIssues", `{
  "data": {
    "user": {
      "issues": {
        "nodes": [
          {
            "title": "title1",
            "url": "https://github.com/notomo/example/issues/1",
            "createdAt": "1971-01-01T00:00:00Z",
            "closedAt": "1971-01-02T00:00:00Z"
          },
          {
            "title": "ignored",
            "url": "https://github.com/notomo/example/issues/8888",
            "createdAt": "1972-02-01T00:00:00Z",
            "closedAt": "1971-01-02T00:00:00Z"
          }
        ],
        "pageInfo": {
          "endCursor": "88888888888888888888888888888888888888888888888888888888",
          "hasNextPage": false
        }
      }
    }
  }
}`),
		gqltest.WithOK("ReportedIssues", `{
  "data": {
    "user": {
      "issues": {
        "nodes": [
          {
            "title": "title2",
            "url": "https://github.com/notomo/example/issues/2",
            "createdAt": "1972-02-01T00:00:00Z",
            "closedAt": null
          },
          {
            "title": "ignored",
            "url": "https://github.com/notomo/example/issues/8888",
            "createdAt": "1972-02-01T00:00:00Z",
            "closedAt": null
          }
        ],
        "pageInfo": {
          "endCursor": "88888888888888888888888888888888888888888888888888888888",
          "hasNextPage": false
        }
      }
    }
  }
}`),
	)
	require.NoError(t, err)

	noUserName := ""
	output := &bytes.Buffer{}
	require.NoError(t, Run(
		gql,
		noUserName,
		1,
		"1970-01-01",
		"json",
		output,
	))

	want := `{
  "from": "1970-01-01T00:00:00Z",
  "closedIssues": [
    {
      "title": "title1",
      "url": "https://github.com/notomo/example/issues/1",
      "created_at": "1971-01-01T00:00:00Z",
      "closed_at": "1971-01-02T00:00:00Z"
    }
  ],
  "reportedIssues": [
    {
      "title": "title2",
      "url": "https://github.com/notomo/example/issues/2",
      "created_at": "1972-02-01T00:00:00Z"
    }
  ]
}`
	got := output.String()
	assert.JSONEq(t, want, got)
}
