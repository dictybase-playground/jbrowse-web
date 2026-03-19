package main

import (
	"context"
	"dictybase-playground/jbrowse-cli/internal/jbrowsemanager"
	"dictybase-playground/jbrowse-cli/internal/server"
	"fmt"
	"os"

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	cli "github.com/urfave/cli/v3"
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
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return F.Pipe2(
				jbrowsemanager.CreateParams{
					Cfg: jbrowsemanager.NewConfig(),
					Ctx: ctx,
				},
				jbrowsemanager.RunCreate,
				E.Fold(
					F.Identity[error],
					func(dr jbrowsemanager.DownloadResult) error {
						defer dr.Body.Close()
						fmt.Printf(
							"Downloaded jbrowse-web %s\n",
							dr.Version,
						)
						return nil
					},
				),
			)
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

	if err := app.Run(
		context.Background(),
		os.Args,
	); err != nil {
		os.Exit(1)
	}
}
