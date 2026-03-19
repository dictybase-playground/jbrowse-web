package jbrowsemanager

import (
	"context"
	"fmt"
	"io"
	"net/http"

	A "github.com/IBM/fp-go/v2/array"
	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	O "github.com/IBM/fp-go/v2/option"
	gh "github.com/google/go-github/v84/github"
)

type DownloadResult struct {
	Body    io.ReadCloser
	Version string
}

type releaseAsset struct {
	ID      int64
	Version string
}

type CreateParams struct {
	Cfg Config
	Ctx context.Context
}

func getLatestRelease(
	params CreateParams,
) IOE.IOEither[error, *gh.RepositoryRelease] {
	return IOE.TryCatchError(
		func() (*gh.RepositoryRelease, error) {
			release, _, err := params.Cfg.Client.Repositories.GetLatestRelease(
				params.Ctx,
				params.Cfg.Owner,
				params.Cfg.Repo,
			)
			return release, err
		},
	)
}

func FetchReleases(
	cfg Config,
	ctx context.Context,
) IOE.IOEither[error, []*gh.RepositoryRelease] {
	return F.Pipe1(
		IOE.TryCatchError(
			func() ([]*gh.RepositoryRelease, error) {
				releases, _, err := cfg.Client.Repositories.ListReleases(
					ctx,
					cfg.Owner,
					cfg.Repo,
					&gh.ListOptions{},
				)
				return releases, err
			},
		),
		IOE.Map[error](A.Filter(isDownloadableRelease)),
	)
}

func downloadAsset(
	cfg Config,
	ctx context.Context,
	id int64,
) IOE.IOEither[error, io.ReadCloser] {
	return IOE.TryCatchError(func() (io.ReadCloser, error) {
		rc, _, err := cfg.Client.Repositories.DownloadReleaseAsset(
			ctx,
			cfg.Owner,
			cfg.Repo,
			id,
			http.DefaultClient,
		)
		return rc, err
	})
}

func extractReleaseAsset(
	release *gh.RepositoryRelease,
) E.Either[error, releaseAsset] {
	return F.Pipe5(
		release.Assets,
		A.FindFirst(isBuildAsset),
		O.Map(getAssetID),
		E.FromOption[int64](func() error {
			return fmt.Errorf(
				"no jbrowse-web asset in release %s",
				release.GetTagName(),
			)
		}),
		E.FilterOrElse(
			func(id int64) bool { return id != 0 },
			func(id int64) error { return fmt.Errorf("build asset has invalid id: %d", id) },
		),
		E.Map[error](func(id int64) releaseAsset {
			return releaseAsset{
				ID:      id,
				Version: release.GetTagName(),
			}
		}),
	)
}

func Create(
	params CreateParams,
) IOE.IOEither[error, DownloadResult] {
	return F.Pipe3(
		params,
		getLatestRelease,
		IOE.ChainEitherK(extractReleaseAsset),
		IOE.Chain(
			func(ra releaseAsset) IOE.IOEither[error, DownloadResult] {
				return F.Pipe1(
					downloadAsset(params.Cfg, params.Ctx, ra.ID),
					IOE.Map[error](
						func(body io.ReadCloser) DownloadResult {
							return DownloadResult{
								Body:    body,
								Version: ra.Version,
							}
						},
					),
				)
			},
		),
	)
}

func RunCreate(
	params CreateParams,
) E.Either[error, DownloadResult] {
	return toEither(Create(params))
}

func toEither[ERR, A any](
	ioe IOE.IOEither[ERR, A],
) E.Either[ERR, A] {
	return ioe()
}
