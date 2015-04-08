package main

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
	mux    *http.ServeMux
	client *octokit.Client
	server *httptest.Server
)

func setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	serverURL, _ := url.Parse(server.URL)

	client = octokit.NewClientWith(
		serverURL.String(),
		"test user agent",
		octokit.TokenAuth{AccessToken: "token"},
		nil,
	)
}

func tearDown() {
	server.Close()
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
	p := path.Join(pwd, "fixtures", f)
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

func setupMux(t *testing.T, resourceType string) {
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

func TestAllPullRequests(t *testing.T) {
	setup()
	setupMux(t, "pulls")
	defer tearDown()

	r := gitHubRepository{"docker", "docker", client}
	prs, err := r.PullRequests("open", "asc")
	if err != nil {
		t.Fatal(err)
	}

	if len(prs) != 4 {
		t.Fatalf("Expected 4 prs but it was %d\n", len(prs))
	}
}

func TestAllIssues(t *testing.T) {
	setup()
	setupMux(t, "issues")
	defer tearDown()

	r := gitHubRepository{"docker", "docker", client}
	issues, err := r.Issues("open", "asc")
	if err != nil {
		t.Fatal(err)
	}

	if len(issues) != 4 {
		t.Fatalf("Expected 4 issues but it was %d\n", len(issues))
	}
}
