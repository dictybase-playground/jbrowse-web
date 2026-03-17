package jbrowse_manager

import (
	"context"
	"fmt"
	A "github.com/IBM/fp-go/v2/array"
	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	O "github.com/IBM/fp-go/v2/option"
	S "github.com/IBM/fp-go/v2/string"
	gh "github.com/google/go-github/v84/github"
	"io"
	"net/http"
	"time"
)

var (
	downloadAsset = F.Bind12of3(E.Eitherize3(uncurriedDownloadAsset))
)

type githubManager struct {
	client *gh.Client
	owner  string
	repo   string
}

func (ghm *githubManager) getLatestRelease(ctx context.Context) (*gh.RepositoryRelease, error) {
	latest, _, err := ghm.client.Repositories.GetLatestRelease(ctx, ghm.owner, ghm.repo)

	if err != nil {
		return nil, fmt.Errorf("could not get latest repository release: %s", err)
	}

	return latest, nil
}

func isNotPrerelease(release *gh.RepositoryRelease) bool {
	return !(*release.Prerelease)
}

func isVersionedRelease(release *gh.RepositoryRelease) bool {
	return F.Pipe1(release.GetTagName(), S.Includes("v"))
}

func hasBuildAsset(release *gh.RepositoryRelease) bool {
	return F.Pipe1(release.Assets, A.Any(isBuildAsset))
}

func isBuildAsset(asset *gh.ReleaseAsset) bool {
	return F.Pipe1(asset.GetName(), S.Includes("jbrowse-web"))
}

func getAssetID(asset *gh.ReleaseAsset) int64 { return asset.GetID() }
func (ghm *githubManager) FetchReleases(ctx context.Context) ([]*gh.RepositoryRelease, error) {
	releases, _, err := ghm.client.Repositories.ListReleases(ctx, ghm.owner, ghm.repo, &gh.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("could not get repository releases: %s", err)
	}
	filteredReleases := F.Pipe3(releases, A.Filter(isVersionedRelease), A.Filter(isNotPrerelease), A.Filter(hasBuildAsset))

	return filteredReleases, nil
}

func uncurriedDownloadAsset(ghm *githubManager, ctx context.Context, id int64) (io.ReadCloser, error) {
	rc, _, err := ghm.client.Repositories.DownloadReleaseAsset(ctx, ghm.owner, ghm.repo, id, http.DefaultClient)
	return rc, err
}

// Get the latest release with that has build assets labeled `jbrowse-web`
func Create(ctx context.Context) error {
	ghm := &githubManager{
		client: gh.NewClient(&http.Client{Timeout: time.Second * 10}),
		owner:  "GMOD",
		repo:   "jbrowse-components",
	}

	latest, err := ghm.getLatestRelease(ctx)

	if err != nil {
		return fmt.Errorf("could not get latest repository release: %s", err)
	}

	F.Pipe5(
		latest.Assets,
		A.FindFirst(isBuildAsset),
		O.Map(getAssetID),
		E.FromOption[int64](func() error { return fmt.Errorf("could not find build asset in latest release") }),
		E.FilterOrElse(
			func(id int64) bool { return id != 0 },
			func(id int64) error { return fmt.Errorf("build asset has invalid id: %d", id) }),
		E.Map[error](downloadAsset(ghm, ctx)),
	)

	return nil
}
