package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/cli/go-gh"
	"github.com/notomo/gh-retrospect/retrospect"
	"github.com/notomo/gh-retrospect/retrospect/outputter"
	"github.com/notomo/gh-retrospect/retrospect/query"

	"github.com/urfave/cli/v2"
)

const (
	paramUser   = "user"
	paramLimit  = "limit"
	paramFrom   = "from"
	paramOutput = "output"
)

func main() {
	app := &cli.App{
		Name: "gh-retrospect",
		Action: func(c *cli.Context) error {
			return Run(
				c.String(paramUser),
				c.Int(paramLimit),
				c.String(paramFrom),
				c.String(paramOutput),
				os.Stdout,
			)
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  paramUser,
				Value: "",
				Usage: "user name (default: authenticated user)",
			},
			&cli.IntFlag{
				Name:  paramLimit,
				Value: 100,
				Usage: "limit to retrieve by each section",
			},
			&cli.StringFlag{
				Name:  paramFrom,
				Value: "",
				Usage: "YYYY-mm-dd format date, default: last week date",
			},
			&cli.StringFlag{
				Name:  paramOutput,
				Value: "json",
				Usage: "outputter type",
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func Run(
	userName string,
	limit int,
	fromDate string,
	outputterType string,
	writer io.Writer,
) error {
	gql, err := gh.GQLClient(nil)
	if err != nil {
		return fmt.Errorf("create gql client: %w", err)
	}

	from, err := retrospect.ParseDate(fromDate)
	if err != nil {
		return fmt.Errorf("parse date: %w", err)
	}

	client := query.NewClient(gql, from, limit)

	outputter, err := outputter.Get(outputterType)
	if err != nil {
		return fmt.Errorf("get outputter: %w", err)
	}

	collected, err := retrospect.Collect(client, userName)
	if err != nil {
		return fmt.Errorf("collect: %w", err)
	}

	if err := outputter.Output(writer, collected); err != nil {
		return fmt.Errorf("output: %w", err)
	}
	return nil
}
