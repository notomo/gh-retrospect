package main

import (
	"bytes"
	"testing"

	"github.com/notomo/gh-retrospect/retrospect/gqltest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	t.Setenv("TZ", "UTC")

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
    "search": {
      "nodes": [
        {
          "title": "title1",
          "url": "https://github.com/notomo/example/issues/1",
          "createdAt": "1971-01-01T00:00:00Z",
          "closedAt": "1971-01-01T01:00:00Z",
          "updatedAt": "1971-01-01T01:00:00Z",
          "labels": {
            "nodes": [
              {
                "name": "label1"
              }
            ]
          }
        },
        {
          "title": "ignored",
          "url": "https://github.com/notomo/example/issues/8888",
          "createdAt": "1971-08-08T00:00:00Z",
          "closedAt": "1971-08-08T08:00:00Z",
          "updatedAt": "1971-08-08T08:00:00Z",
          "labels": {
            "nodes": [
              {
                "name": "label1"
              }
            ]
          }
        }
      ],
      "pageInfo": {
        "endCursor": "88888888888888888888888888888888888888888888888888888888",
        "hasNextPage": false
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
            "createdAt": "1972-01-02T00:00:00Z",
            "closedAt": null,
            "labels": {
              "nodes": []
            }
          },
          {
            "title": "ignored",
            "url": "https://github.com/notomo/example/issues/8888",
            "createdAt": "1971-08-08T00:00:00Z",
            "closedAt": null,
            "labels": {
              "nodes": []
            }
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

		gqltest.WithOK("MergedPullRequests", `{
  "data": {
    "search": {
      "nodes": [
        {
          "title": "title3",
          "url": "https://github.com/notomo/example/pull/3",
          "createdAt": "1971-01-03T00:00:00Z",
          "closedAt": "1972-01-03T03:00:00Z",
          "labels": {
            "nodes": [
              {
                "name": "label1"
              }
            ]
          }
        },
        {
          "title": "ignored",
          "url": "https://github.com/notomo/example/pull/8888",
          "createdAt": "1971-08-08T00:00:00Z",
          "closedAt": "1971-08-08T08:00:00Z",
          "labels": {
            "nodes": []
          }
        }
      ],
      "pageInfo": {
        "endCursor": "88888888888888888888888888888888888888888888888888888888",
        "hasNextPage": false
      }
    }
  }
}`),

		gqltest.WithOK("ReviewedPullRequests", `{
  "data": {
    "search": {
      "nodes": [
        {
          "title": "title4",
          "url": "https://github.com/notomo/example/pull/4",
          "createdAt": "1971-01-04T00:00:00Z",
          "closedAt": "1972-01-04T04:00:00Z",
          "labels": {
            "nodes": [
              {
                "name": "label1"
              }
            ]
          }
        },
        {
          "title": "ignored",
          "url": "https://github.com/notomo/example/pull/8888",
          "createdAt": "1971-08-08T00:00:00Z",
          "closedAt": "1971-08-08T08:00:00Z",
          "labels": {
            "nodes": []
          }
        }
      ],
      "pageInfo": {
        "endCursor": "88888888888888888888888888888888888888888888888888888888",
        "hasNextPage": false
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
		"",
		"json",
		output,
	))

	want := `{
  "from": "1970-01-01T00:00:00Z",
  "to": "0001-01-01T00:00:00Z",
  "closed_issues": [
    {
      "title": "title1",
      "url": "https://github.com/notomo/example/issues/1",
      "created_at": "1971-01-01T00:00:00Z",
      "closed_at": "1971-01-01T01:00:00Z",
      "updated_at": "1971-01-01T01:00:00Z",
      "label_names": ["label1"]
    }
  ],
  "reported_issues": [
    {
      "title": "title2",
      "url": "https://github.com/notomo/example/issues/2",
      "created_at": "1972-01-02T00:00:00Z",
      "label_names": []
    }
  ],
  "merged_pull_requests": [
    {
      "title": "title3",
      "url": "https://github.com/notomo/example/pull/3",
      "created_at": "1971-01-03T00:00:00Z",
      "closed_at": "1972-01-03T03:00:00Z",
      "label_names": ["label1"]
    }
  ],
  "reviewed_pull_requests": [
    {
      "title": "title4",
      "url": "https://github.com/notomo/example/pull/4",
      "created_at": "1971-01-04T00:00:00Z",
      "closed_at": "1972-01-04T04:00:00Z",
      "label_names": ["label1"]
    }
  ]
}`
	got := output.String()
	assert.JSONEq(t, want, got)
}
