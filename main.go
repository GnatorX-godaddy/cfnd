package main

import (
	"log"
	"os"

	"github.com/GnatorX/cfnd/pkg/ctl"

	"github.com/urfave/cli/v2"
)

func main() {
	var region string
	var stackName string
	var outputFile string
	var readOnly bool
	var all bool

	app := &cli.App{
		Name:  "cfnd",
		Usage: "Find Cloudtrail events for your CloudFormation errors",
		Action: func(c *cli.Context) error {
			ctl.Find(c.Context, stackName, region, outputFile, readOnly, all)
			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "stackname",
				Usage:       "Name of the stack",
				Destination: &stackName,
				Aliases:     []string{"s"},
			},
			&cli.StringFlag{
				Name:        "output",
				Usage:       "Output file name",
				Destination: &outputFile,
				Aliases:     []string{"o"},
				DefaultText: "cf_error.json",
				Value:       "cf_error.json",
			},
			&cli.BoolFlag{
				Name:        "readonly",
				Usage:       "Return readonly events from CloudTrail. Add the flag if you want readonly to be true",
				Aliases:     []string{"ro"},
				Destination: &readOnly,
				Value:       false,
			},
			&cli.BoolFlag{
				Name:        "all",
				Usage:       "Return all events from CloudTrail. By default, only Events with error is returned. Add the flag if you want all events",
				Aliases:     []string{"a"},
				Destination: &all,
				Value:       false,
			},
			&cli.StringFlag{
				Name:        "region",
				Usage:       "AWS region for the search",
				Destination: &region,
				Aliases:     []string{"r"},
				DefaultText: "us-west-2",
				Value:       "us-west-2",
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
