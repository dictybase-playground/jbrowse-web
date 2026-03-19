package jbrowsemanager

import (
	A "github.com/IBM/fp-go/v2/array"
	F "github.com/IBM/fp-go/v2/function"
	P "github.com/IBM/fp-go/v2/predicate"
	S "github.com/IBM/fp-go/v2/string"
	gh "github.com/google/go-github/v84/github"
)

var (
	isBuildAsset = func(asset *gh.ReleaseAsset) bool {
		return F.Pipe1(
			asset.GetName(),
			S.Includes("jbrowse-web"),
		)
	}
	hasBuildAsset = func(release *gh.RepositoryRelease) bool {
		return F.Pipe1(release.Assets, A.Any(isBuildAsset))
	}
	isPrerelease = func(release *gh.RepositoryRelease) bool {
		return release.GetPrerelease()
	}
	isNotPrerelease    = P.Not(isPrerelease)
	isVersionedRelease = func(release *gh.RepositoryRelease) bool {
		return F.Pipe1(release.GetTagName(), S.Includes("v"))
	}
	isDownloadableRelease = F.Pipe2(
		isVersionedRelease,
		P.And(isNotPrerelease),
		P.And(hasBuildAsset),
	)
)

func getAssetID(
	asset *gh.ReleaseAsset,
) int64 {
	return asset.GetID()
}
