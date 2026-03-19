package jbrowsemanager

import (
	"context"
	"fmt"

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	P "github.com/IBM/fp-go/v2/pair"
	cli "github.com/urfave/cli/v3"
)

var (
	createError = F.Bind1st(
		P.MakePair[DownloadResult, error],
		DownloadResult{},
	)
	createSuccess = F.Bind2nd(
		P.MakePair[DownloadResult, error],
		(error)(nil),
	)
)

func CreateAction(ctx context.Context, cmd *cli.Command) error {
	output := F.Pipe2(
		CreateParams{
			Cfg: NewConfig(),
			Ctx: ctx,
		},
		RunCreate,
		E.Fold(createError, createSuccess),
	)
	if err := P.Second(output); err != nil {
		return err
	}
	dr := P.First(output)
	defer dr.Body.Close()
	fmt.Printf("Downloaded jbrowse-web %s\n", dr.Version)
	return nil
}
