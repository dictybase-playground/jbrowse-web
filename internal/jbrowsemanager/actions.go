package jbrowsemanager

import (
	"context"
	"fmt"

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	cli "github.com/urfave/cli/v3"
)

func CreateAction(ctx context.Context, cmd *cli.Command) error {
	return F.Pipe2(
		CreateParams{
			Cfg: NewConfig(),
			Ctx: ctx,
		},
		RunCreate,
		E.Fold(
			F.Identity[error],
			func(dr DownloadResult) error {
				defer dr.Body.Close()
				fmt.Printf(
					"Downloaded jbrowse-web %s\n",
					dr.Version,
				)
				return nil
			},
		),
	)
}
