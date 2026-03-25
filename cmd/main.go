package main

import (
	"context"
	"dictybase-playground/jbrowse-web/internal/server"
	cli "github.com/urfave/cli/v3"
	"os"
)

func serverCommands() *cli.Command {
	return &cli.Command{
		Name:  "dev",
		Usage: "Starts development server",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return server.StartServer(ctx)
		},
	}
}

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
			serverCommands(),
		},
	}
}

func main() {

	app := setupApp()

	if err := app.Run(context.Background(), os.Args); err != nil {
		os.Exit(1)
	}
}
