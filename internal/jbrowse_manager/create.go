package jbrowse_manager

import (
	"context"
	"fmt"
	A "github.com/IBM/fp-go/v2/array"
	F "github.com/IBM/fp-go/v2/function"
	O "github.com/IBM/fp-go/v2/option"
	S "github.com/IBM/fp-go/v2/string"
	gh "github.com/google/go-github/v84/github"
	"io"
	"net/http"
	"strings"
	"time"
)

type GithubManager struct {
	client *gh.Client
	owner  string
	repo   string
}

func listReleases(ctx context.Context, client *gh.Client) ([]*gh.RepositoryRelease, error) {
	var filteredReleases []*gh.RepositoryRelease
	page := 1

	for {
		opts := &gh.ListOptions{
			Page:    page,
			PerPage: 30,
		}

		releases, resp, err := client.Repositories.ListReleases(ctx, "GMOD", "jbrowse-components", opts)
		if err != nil {
			return nil, fmt.Errorf("could not fetch JBrowse repository releases: %w", err)
		}

		// Filter releases whose tag_name starts with 'v'
		for _, release := range releases {
			if release.TagName != nil && strings.HasPrefix(*release.TagName, "v") {
				filteredReleases = append(filteredReleases, release)
			}
		}

		// Check if there are more pages
		if resp.NextPage == 0 || len(releases) == 0 {
			break
		}
		page = resp.NextPage
	}

	return filteredReleases, nil
}

func (ghm *GithubManager) getLatestRelease(ctx context.Context) (*gh.RepositoryRelease, error) {

	latest, _, err := ghm.client.Repositories.GetLatestRelease(ctx, ghm.owner, ghm.repo)

	if err != nil {
		return nil, fmt.Errorf("could not get latest repository release: %s", err)
	}

	return latest, nil
}

func isNotPrerelease(release *gh.RepositoryRelease) bool {
	return !(*release.Prerelease)
}

func getAssetName(asset *gh.ReleaseAsset) string { return asset.GetName() }
func getAssetID(asset *gh.ReleaseAsset) int64    { return asset.GetID() }

func downloadAsset(ghm *GithubManager, ctx context.Context, id int64) (io.ReadCloser, string, error) {
	return ghm.client.Repositories.DownloadReleaseAsset(ctx, ghm.owner, ghm.repo, id, http.DefaultClient)
}

func isBuildAsset(asset *gh.ReleaseAsset) bool {
	return F.Pipe2(asset, getAssetName, S.Includes("jbrowse-web"))
}

func hasBuildAsset(release *gh.RepositoryRelease) bool {
	return F.Pipe1(release.Assets, A.Any(isBuildAsset))
}

func (ghm *GithubManager) FetchReleases(ctx context.Context) ([]*gh.RepositoryRelease, error) {
	versions, _, err := ghm.client.Repositories.ListReleases(ctx, ghm.owner, ghm.repo, &gh.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("could not get repository releases: %s", err)
	}
	filteredReleases := F.Pipe2(versions, A.Filter(isNotPrerelease), A.Filter(hasBuildAsset))

	return filteredReleases, nil
}

// Get the latest release with that has build assets labeled `jbrowse-web`
func Create(ctx context.Context) error {
	ghm := &GithubManager{
		client: gh.NewClient(&http.Client{Timeout: time.Second * 10}),
		owner:  "GMOD",
		repo:   "jbrowse-components",
	}

	latest, err := ghm.getLatestRelease(ctx)

	if err != nil {
		return fmt.Errorf("could not get latest repository release: %s", err)
	}

	F.Pipe3(latest.Assets, A.FindFirst(isBuildAsset), O.Map(getAssetID), O.Fold(
		func() error { return nil },
		func(id int64) error { return nil },
	))

	return nil
}
