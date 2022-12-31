package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/cli/go-gh"
	"github.com/cli/go-gh/pkg/api"
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
	paramLog    = "log"
)

func main() {
	app := &cli.App{
		Name: "gh-retrospect",
		Action: func(c *cli.Context) error {
			opts := &api.ClientOptions{}
			logFilePath := c.String(paramLog)
			if logFilePath != "" {
				f, err := os.Create(logFilePath)
				if err != nil {
					return fmt.Errorf("create log file: %w", err)
				}
				defer f.Close()
				opts.Log = f
				opts.LogVerboseHTTP = true
			}
			gql, err := gh.GQLClient(opts)
			if err != nil {
				return fmt.Errorf("create gql client: %w", err)
			}
			return Run(
				gql,
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
			&cli.StringFlag{
				Name:  paramLog,
				Value: "",
				Usage: "log file path",
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func Run(
	gql api.GQLClient,
	userName string,
	limit int,
	fromDate string,
	outputterType string,
	writer io.Writer,
) error {
	from, err := retrospect.ParseDate(fromDate)
	if err != nil {
		return fmt.Errorf("parse date: %w", err)
	}

	outputter, err := outputter.Get(outputterType)
	if err != nil {
		return fmt.Errorf("get outputter: %w", err)
	}

	client := query.NewClient(gql)
	collected, err := retrospect.Collect(client, userName, from, limit)
	if err != nil {
		return fmt.Errorf("collect: %w", err)
	}

	if err := outputter.Output(writer, collected); err != nil {
		return fmt.Errorf("output: %w", err)
	}
	return nil
}
