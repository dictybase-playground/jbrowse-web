package main

import (
	"context"
	cli "github.com/urfave/cli/v3"
	"os"
)

func jbrowseCommands() *cli.Command {
	return &cli.Command{
		Name:  "create",
		Usage: "Downloads Jbrowse",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "version",
				Usage: "specifies the version of jbrowse to fetch",
			},
		},
		Action: func(context.Context, *cli.Command) error {
			return nil
		},
	}
}

func setupApp() *cli.Command {
	return &cli.Command{
		Name:  "jbrowse-launch",
		Usage: "Jbrowse management tools",
		Commands: []*cli.Command{
			jbrowseCommands(),
		},
	}
}

func main() {

	app := setupApp()

	if err := app.Run(context.Background(), os.Args); err != nil {
		os.Exit(1)
	}
}
