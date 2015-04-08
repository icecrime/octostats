package octokit

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGistsService_One(t *testing.T) {
	setup()
	defer tearDown()

	stubGet(t, "/gists/a6bea192debdbec0d4ab", "gist", nil)

	url, _ := GistsURL.Expand(M{"gist_id": "a6bea192debdbec0d4ab"})
	gist, result := client.Gists(url).One()

	assert.False(t, result.HasError())
	assert.Equal(t, "a6bea192debdbec0d4ab", gist.ID)
	assert.Len(t, gist.Files, 1)

	file := gist.Files["grep_cellar"]
	assert.Equal(t, "grep_cellar", file.FileName)
	assert.Equal(t, "text/plain", file.Type)
	assert.Equal(t, "", file.Language)
	assert.Equal(t, "https://gist.githubusercontent.com/jingweno/a6bea192debdbec0d4ab/raw/80757419d2bd4cfddf7c6be24308eca11b3c330e/grep_cellar", file.RawURL)
	assert.Equal(t, 8107, file.Size)
	assert.Equal(t, false, file.Truncated)
}

func TestGistsService_Raw(t *testing.T) {
	setup()
	defer tearDown()

	stubGet(t, "/gists/a6bea192debdbec0d4ab", "gist", nil)

	mux.HandleFunc("/jingweno/a6bea192debdbec0d4ab/raw/80757419d2bd4cfddf7c6be24308eca11b3c330e/grep_cellar", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		assert.Equal(t, "gist.githubusercontent.com", r.Host)
		testHeader(t, r, "Accept", textMediaType)
		respondWith(w, "hello")
	})

	url, _ := GistsURL.Expand(M{"gist_id": "a6bea192debdbec0d4ab"})
	body, result := client.Gists(url).Raw()

	assert.False(t, result.HasError())
	content, err := ioutil.ReadAll(body)
	assert.NoError(t, err)
	assert.Equal(t, "hello", string(content))
}
