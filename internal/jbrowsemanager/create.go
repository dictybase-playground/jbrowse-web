package jbrowsemanager

import (
	"context"
	"fmt"
	"io"
	"net/http"

	A "github.com/IBM/fp-go/v2/array"
	E "github.com/IBM/fp-go/v2/either"
	EQ "github.com/IBM/fp-go/v2/eq"
	F "github.com/IBM/fp-go/v2/function"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	P "github.com/IBM/fp-go/v2/pair"
	Pred "github.com/IBM/fp-go/v2/predicate"
	gh "github.com/google/go-github/v84/github"
)

var (
	int64Eq = EQ.FromStrictEquals[int64]()

	notEqualInt64 = F.Flow2(EQ.Equals(int64Eq), Pred.Not)
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
) IOE.IOEither[error, P.Pair[CreateParams, *gh.RepositoryRelease]] {
	return IOE.TryCatchError(
		func() (P.Pair[CreateParams, *gh.RepositoryRelease], error) {
			release, _, err := params.Cfg.Client.Repositories.GetLatestRelease(
				params.Ctx,
				params.Cfg.Owner,
				params.Cfg.Repo,
			)
			return P.MakePair(params, release), err
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
	pair P.Pair[CreateParams, releaseAsset],
) IOE.IOEither[error, DownloadResult] {
	cfg, asset := P.First(pair).Cfg, P.Second(pair)
	return F.Pipe1(
		IOE.TryCatchError(func() (io.ReadCloser, error) {
			rc, _, err := cfg.Client.Repositories.DownloadReleaseAsset(
				P.First(pair).Ctx,
				cfg.Owner,
				cfg.Repo,
				asset.ID,
				http.DefaultClient,
			)
			return rc, err
		}),
		IOE.Map[error](func(body io.ReadCloser) DownloadResult {
			return DownloadResult{
				Body:    body,
				Version: asset.Version,
			}
		}),
	)
}

func extractReleaseAsset(
	pair P.Pair[CreateParams, *gh.RepositoryRelease],
) E.Either[error, P.Pair[CreateParams, releaseAsset]] {
	release := P.Second(pair)
	return F.Pipe5(
		release.Assets,
		A.FindFirst(isBuildAsset),
		E.FromOption[*gh.ReleaseAsset](func() error {
			return fmt.Errorf(
				"no jbrowse-web asset in release %s",
				release.GetTagName(),
			)
		}),
		E.Map[error](getAssetID),
		E.Chain(
			E.FromPredicate(
				notEqualInt64(0),
				invalidAssetIDError,
			)),
		E.Map[error](
			func(id int64) P.Pair[CreateParams, releaseAsset] {
				return P.MakePair(P.First(pair), releaseAsset{
					ID:      id,
					Version: release.GetTagName(),
				})
			},
		),
	)
}

func toEither[ERR, A any](
	ioe IOE.IOEither[ERR, A],
) E.Either[ERR, A] {
	return ioe()
}

func invalidAssetIDError(id int64) error {
	return fmt.Errorf("invalid asset id: %d", id)
}
