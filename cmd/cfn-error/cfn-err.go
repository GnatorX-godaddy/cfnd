package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "cfn-err",
		Usage: "Find Cloudtrail events for your CloudFormation errors",
		Action: func(c *cli.Context) error {
			fmt.Println("Surprise")
			return nil
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
