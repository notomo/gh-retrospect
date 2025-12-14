package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/henvic/httpretty"
	"github.com/notomo/gh-retrospect/retrospect"
	"github.com/notomo/gh-retrospect/retrospect/outputter"
	"github.com/notomo/gh-retrospect/retrospect/query"
	"github.com/notomo/httpwriter"

	"github.com/urfave/cli/v2"
)

const (
	paramUser   = "user"
	paramLimit  = "limit"
	paramFrom   = "from"
	paramTo     = "to"
	paramOutput = "output"
	paramLog    = "log"
)

func main() {
	app := &cli.App{
		Name: "gh-retrospect",
		Action: func(c *cli.Context) error {
			opts := api.ClientOptions{}
			logDirPath := c.String(paramLog)
			if logDirPath != "" {
				opts.Transport = &httpwriter.Transport{
					TransportFactory: func(writer io.Writer) http.RoundTripper {
						logger := &httpretty.Logger{
							Time:           true,
							TLS:            false,
							RequestHeader:  true,
							RequestBody:    true,
							ResponseHeader: true,
							ResponseBody:   true,
							Formatters:     []httpretty.Formatter{&httpretty.JSONFormatter{}},
						}
						logger.SetOutput(writer)
						return logger.RoundTripper(nil)
					},
					GetWriter: httpwriter.MustDirectoryWriter(
						&httpwriter.Directory{Path: logDirPath},
					),
				}
			}
			gql, err := api.NewGraphQLClient(opts)
			if err != nil {
				return fmt.Errorf("create gql client: %w", err)
			}
			return Run(
				gql,
				c.String(paramUser),
				c.Int(paramLimit),
				c.String(paramFrom),
				c.String(paramTo),
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
				Usage: "YYYY-mm-dd format date or duration, default: last week date",
			},
			&cli.StringFlag{
				Name:  paramTo,
				Value: "",
				Usage: "YYYY-mm-dd format date or duration, default: no upper bound",
			},
			&cli.StringFlag{
				Name:  paramOutput,
				Value: "json",
				Usage: "outputter type",
			},
			&cli.StringFlag{
				Name:  paramLog,
				Value: "",
				Usage: "log directory path",
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func Run(
	gql *api.GraphQLClient,
	userName string,
	limit int,
	fromDateOrDuration string,
	toDateOrDuration string,
	outputterType string,
	writer io.Writer,
) error {
	from, err := retrospect.ParseFrom(fromDateOrDuration)
	if err != nil {
		return fmt.Errorf("parse from date: %w", err)
	}

	to, err := retrospect.ParseTo(toDateOrDuration)
	if err != nil {
		return fmt.Errorf("parse to date: %w", err)
	}

	outputter, err := outputter.Get(outputterType)
	if err != nil {
		return fmt.Errorf("get outputter: %w", err)
	}

	client := query.NewClient(gql)
	collected, err := retrospect.Collect(client, userName, from, to, limit)
	if err != nil {
		return fmt.Errorf("collect: %w", err)
	}

	if err := outputter.Output(writer, collected); err != nil {
		return fmt.Errorf("output: %w", err)
	}
	return nil
}
