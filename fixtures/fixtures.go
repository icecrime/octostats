package fixtures

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path"
	"strconv"
	"testing"

	"github.com/octokit/go-octokit/octokit"
)

var (
	server *httptest.Server
	mux    *http.ServeMux
	Client *octokit.Client
)

func Setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	serverURL, _ := url.Parse(server.URL)

	Client = octokit.NewClientWith(
		serverURL.String(),
		"test user agent",
		octokit.TokenAuth{AccessToken: "token"},
		nil,
	)
}

func TearDown() {
	server.Close()
}

func SetupMux(t *testing.T, resourceType string) {
	rPath := fmt.Sprintf("/repos/docker/docker/%s", resourceType)

	mux.HandleFunc(rPath, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Fatalf("Expected GET but it was %s\n", r.Method)
		}

		page := r.URL.Query().Get("page")
		header := w.Header()

		if page == "" {
			link := fmt.Sprintf(`<%s>; rel="next", <%s>; rel="last"`, testURLOf(rPath+"?page=2"), testURLOf(rPath+"?page=4"))

			header.Set("Link", link)
			respondWithJSON(w, loadFixture(resourceType+"/page1.json"))
		} else {
			p, _ := strconv.Atoi(page)
			next := fmt.Sprintf(rPath+"?page=%d", p+1)
			link := fmt.Sprintf(`<%s>; rel="next", <%s>; rel="last"`, testURLOf(next), testURLOf(rPath+"?page=4"))

			header.Set("Link", link)
			respondWithJSON(w, loadFixture(fmt.Sprintf(resourceType+"/page%s.json", page)))
		}
	})
}

func testURLOf(path string) *url.URL {
	u, _ := url.ParseRequestURI(testURLStringOf(path))
	return u
}

func testURLStringOf(path string) string {
	return fmt.Sprintf("%s/%s", server.URL, path)
}

func loadFixture(f string) string {
	pwd, _ := os.Getwd()
	p := path.Join(pwd, "..", "fixtures", f)
	c, _ := ioutil.ReadFile(p)
	return string(c)
}

func respondWithJSON(w http.ResponseWriter, s string) {
	header := w.Header()
	header.Set("Content-Type", "application/json")
	respondWith(w, s)
}

func respondWith(w http.ResponseWriter, s string) {
	fmt.Fprint(w, s)
}
