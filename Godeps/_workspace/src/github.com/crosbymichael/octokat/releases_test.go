package octokat

import (
	"github.com/bmizerany/assert"
	"net/http"
	"testing"
)

func TestAddPreviewMediaType(t *testing.T) {
	options := addPreviewMediaType(nil)
	assert.Equal(t, PreviewMediaType, options.Headers["Accept"])
}

func TestReleases(t *testing.T) {
	setup()
	defer tearDown()

	mux.HandleFunc("/repos/jingweno/gh/releases", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", PreviewMediaType)
		respondWith(w, loadFixture("releases.json"))
	})

	repo := Repo{UserName: "jingweno", Name: "gh"}
	releases, _ := client.Releases(repo, nil)
	assert.Equal(t, 1, len(releases))

	firstRelease := releases[0]
	assert.Equal(t, 50013, firstRelease.ID)
	assert.Equal(t, "v0.23.0", firstRelease.TagName)
	assert.Equal(t, "master", firstRelease.TargetCommitsh)
	assert.Equal(t, "v0.23.0", firstRelease.Name)
	assert.T(t, !firstRelease.Draft)
	assert.T(t, !firstRelease.Prerelease)
	assert.Equal(t, "* Windows works!: https://github.com/jingweno/gh/commit/6cb80cb09fd9f624a64d85438157955751a9ac70", firstRelease.Body)
	assert.Equal(t, "https://api.github.com/repos/jingweno/gh/releases/50013", firstRelease.URL)
	assert.Equal(t, "https://api.github.com/repos/jingweno/gh/releases/50013/assets", firstRelease.AssetsURL)
	assert.Equal(t, "https://uploads.github.com/repos/jingweno/gh/releases/50013/assets{?name}", firstRelease.UploadURL)
	assert.Equal(t, "https://github.com/jingweno/gh/releases/v0.23.0", firstRelease.HTMLURL)
	assert.Equal(t, "2013-09-23 00:59:10 +0000 UTC", firstRelease.CreatedAt.String())
	assert.Equal(t, "2013-09-23 01:07:56 +0000 UTC", firstRelease.PublishedAt.String())

	firstReleaseAssets := firstRelease.Assets
	assert.Equal(t, 8, len(firstReleaseAssets))

	firstAsset := firstReleaseAssets[0]
	assert.Equal(t, 20428, firstAsset.ID)
	assert.Equal(t, "gh_0.23.0-snapshot_amd64.deb", firstAsset.Name)
	assert.Equal(t, "gh_0.23.0-snapshot_amd64.deb", firstAsset.Label)
	assert.Equal(t, "application/x-deb", firstAsset.ContentType)
	assert.Equal(t, "uploaded", firstAsset.State)
	assert.Equal(t, 1562984, firstAsset.Size)
	assert.Equal(t, 0, firstAsset.DownloadCount)
	assert.Equal(t, "https://api.github.com/repos/jingweno/gh/releases/assets/20428", firstAsset.URL)
	assert.Equal(t, "2013-09-23 01:05:20 +0000 UTC", firstAsset.CreatedAt.String())
	assert.Equal(t, "2013-09-23 01:07:56 +0000 UTC", firstAsset.UpdatedAt.String())
}

func TestCreateRelease(t *testing.T) {
	setup()
	defer tearDown()

	mux.HandleFunc("/repos/octokat/Hello-World/releases", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testHeader(t, r, "Accept", PreviewMediaType)
		testBody(t, r, `{"tag_name":"v1.0.0","target_commitish":"master"}`)
		respondWith(w, loadFixture("create_release.json"))
	})

	repo := Repo{UserName: "octokat", Name: "Hello-World"}
	params := ReleaseParams{
		TagName:         "v1.0.0",
		TargetCommitish: "master",
	}
	options := Options{Params: params}
	release, _ := client.CreateRelease(repo, &options)

	assert.Equal(t, "v1.0.0", release.TagName)
}
