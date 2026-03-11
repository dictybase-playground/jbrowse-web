package jbrowse_manager

import (
	"context"
	"net/http"
	"testing"
	"time"

	gh "github.com/google/go-github/v84/github"
)

func TestListReleases(t *testing.T) {
	ctx := context.Background()
	client := gh.NewClient(&http.Client{Timeout: 10 * time.Second})

	releases, err := listReleases(ctx, client)
	if err != nil {
		t.Fatalf("listReleases failed: %v", err)
	}

	t.Logf("Found %d JBrowse releases", len(releases))

	// Print first 5 releases
	for i, release := range releases {
		if i >= 5 {
			break
		}
		if release.TagName != nil {
			t.Logf("  - %s", *release.TagName)
		}
	}
}
