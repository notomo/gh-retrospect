package main

import (
	"fmt"
	"log"
	"os"

	"github.com/cli/go-gh"
	"github.com/notomo/gh-retrospect/retrospect"
	"github.com/notomo/gh-retrospect/retrospect/outputter"
	"github.com/notomo/gh-retrospect/retrospect/query"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:   "gh-retrospect",
		Action: run,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "user",
				Value: "",
				Usage: "user name (default: authenticated user)",
			},
			&cli.IntFlag{
				Name:  "limit",
				Value: 100,
				Usage: "limit to retrieve by each section",
			},
			&cli.StringFlag{
				Name:  "from",
				Value: "",
				Usage: "YYYY-mm-dd format date, default: last week date",
			},
			&cli.StringFlag{
				Name:  "output",
				Value: "json",
				Usage: "outputter type",
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	gql, err := gh.GQLClient(nil)
	if err != nil {
		return fmt.Errorf("create gql client: %w", err)
	}

	from, err := retrospect.ParseDate(c.String("from"))
	if err != nil {
		return fmt.Errorf("parse 'from' date: %w", err)
	}

	client := query.NewClient(gql, from, c.Int("limit"))

	outputter, err := outputter.Get(c.String("output"))
	if err != nil {
		return fmt.Errorf("get outputter: %w", err)
	}

	collected, err := retrospect.Collect(client, c.String("user"))
	if err != nil {
		return fmt.Errorf("collect: %w", err)
	}

	writer := os.Stdout
	if err := outputter.Output(writer, collected); err != nil {
		return fmt.Errorf("output: %w", err)
	}
	return nil
}
