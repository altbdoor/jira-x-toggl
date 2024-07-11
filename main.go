package main

import (
	"log"
	"os"

	"jira-x-toggl/actions"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:     "jira-x-toggl",
		HelpName: "jiggl",
		Usage:    "An application that generates the data based on toggl and Jira API",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "The path to the JSON configuration file",
				Value:   "./config.json",
			},
			&cli.BoolFlag{
				Name:  "debug",
				Usage: "Prints debug information",
				Value: false,
			},
		},
		// Action: actions.RunAction,
		Commands: []*cli.Command{
			{
				Name:    "run",
				Aliases: []string{"r"},
				Usage:   "Fetches toggl and Jira data, and compiles into a CSV file",
				Action:  actions.RunAction,
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "skip-fetch-toggl",
						Aliases: []string{"s"},
						Usage:   "Skip fetching toggl data",
						Value:   false,
					},
					&cli.IntFlag{
						Name:  "start",
						Usage: "The start_date query in toggl in days. E.g., 30 or 90.",
						Value: 90,
					},
					&cli.IntFlag{
						Name:  "end",
						Usage: "The end_date query in toggl in days. E.g., 30 or 90.",
						Value: 0,
					},
				},
			},
			{
				Name:   "config-init",
				Usage:  "Creates a blank JSON configuration file",
				Action: actions.ConfigInitAction,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
