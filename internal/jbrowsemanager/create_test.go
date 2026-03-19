package jbrowsemanager

import (
	"context"
	"testing"

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	gh "github.com/google/go-github/v84/github"
)

func TestFetchReleases(t *testing.T) {
	cfg := NewConfig()
	result := toEither(FetchReleases(cfg, context.Background()))

	if E.IsLeft(result) {
		msg := E.Fold(
			func(err error) string { return err.Error() },
			func(_ []*gh.RepositoryRelease) string { return "" },
		)(result)
		t.Fatalf("FetchReleases failed: %s", msg)
	}

	releases := E.GetOrElse(
		func(_ error) []*gh.RepositoryRelease { return nil },
	)(
		result,
	)
	if len(releases) == 0 {
		t.Fatal("expected at least one release")
	}

	t.Logf("Found %d releases", len(releases))
	for i, release := range releases {
		if i >= 5 {
			break
		}
		t.Logf("  - %s", release.GetTagName())
	}
}

func TestCreate(t *testing.T) {
	result := F.Pipe2(
		CreateParams{
			Cfg: NewConfig(),
			Ctx: context.Background(),
		},
		Create,
		toEither,
	)

	if E.IsLeft(result) {
		msg := E.Fold(
			func(err error) string { return err.Error() },
			func(_ DownloadResult) string { return "" },
		)(result)
		t.Fatalf("Create failed: %s", msg)
	}

	dr := E.GetOrElse(
		func(_ error) DownloadResult { return DownloadResult{} },
	)(
		result,
	)
	if dr.Body == nil {
		t.Fatal("expected non-nil Body in DownloadResult")
	}
	defer dr.Body.Close()

	t.Logf("Create succeeded: version=%s", dr.Version)
}
