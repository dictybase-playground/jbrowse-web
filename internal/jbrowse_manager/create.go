package jbrowse_manager

import (
	"context"
	"fmt"
	gh "github.com/google/go-github/v84/github"
	cli "github.com/urfave/cli/v3"
	"net/http"
	"strings"
)

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

func Create(ctx context.Context, cmd *cli.Command) error {
	client := gh.NewClient(&http.Client{Timeout: 10_000})

	latest, _, err := client.Repositories.GetLatestRelease(ctx, client)

	if err != nil {
		return fmt.Errorf("could not get latest repository release: %s", err)
	}

	// TODO: Process filtered releases
	return nil
}
