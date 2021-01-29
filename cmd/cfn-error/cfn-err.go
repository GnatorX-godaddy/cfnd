package main

import (
	"fmt"
	"log"
	"os"

	"cloudformation-error/pkg/ctl"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "cfnerr",
		Usage: "Find Cloudtrail events for your CloudFormation errors",
		Action: func(c *cli.Context) error {
			fmt.Println(ctl.Find(c.Context, c.Args().Get(0), c.Args().Get(1)))

			return nil
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
