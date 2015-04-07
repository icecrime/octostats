package octokit

import (
	"io"
	"net/url"
	"time"

	"github.com/jingweno/go-sawyer/hypermedia"
)

// CommitsURL is a template for accessing commits in a specific owner's
// repository with a particular sha hash that can be expanded to a full address.
var CommitsURL = Hyperlink("repos/{owner}/{repo}/commits{/sha}")

// Commits creates a CommitsService with a base url.
func (c *Client) Commits(url *url.URL) (commits *CommitsService) {
	commits = &CommitsService{client: c, URL: url}
	return
}

// CommitsService is a service providing access to commits from a particular url
type CommitsService struct {
	client *Client
	URL    *url.URL
}

// All gets a list of all commits associated with the URL of the service
func (c *CommitsService) All() (commits []Commit, result *Result) {
	result = c.client.get(c.URL, &commits)
	return
}

// One gets a specific commit based on the url of the service
func (c *CommitsService) One() (commit *Commit, result *Result) {
	result = c.client.get(c.URL, &commit)
	return
}

// Patch gets a specific commit patch based on the url of the service
func (c *CommitsService) Patch() (patch io.ReadCloser, result *Result) {
	patch, result = c.client.getBody(c.URL, patchMediaType)
	return
}

// CommitFile is a representation of a file within a commit
type CommitFile struct {
	Additions   int    `json:"additions,omitempty"`
	BlobURL     string `json:"blob_url,omitempty"`
	Changes     int    `json:"changes,omitempty"`
	ContentsURL string `json:"contents_url,omitempty"`
	Deletions   int    `json:"deletions,omitempty"`
	Filename    string `json:"filename,omitempty"`
	Patch       string `json:"patch,omitempty"`
	RawURL      string `json:"raw_url,omitempty"`
	Sha         string `json:"sha,omitempty"`
	Status      string `json:"status,omitempty"`
}

// CommitStats represents the statistics on the changes made in a commit
type CommitStats struct {
	Additions int `json:"additions,omitempty"`
	Deletions int `json:"deletions,omitempty"`
	Total     int `json:"total,omitempty"`
}

// CommitCommit is the representation of the metadata regarding the commit as a subset
// of the full commit structure
type CommitCommit struct {
	Author struct {
		Date  *time.Time `json:"date,omitempty"`
		Email string     `json:"email,omitempty"`
		Name  string     `json:"name,omitempty"`
	} `json:"author,omitempty"`
	CommentCount int `json:"comment_count,omitempty"`
	Committer    struct {
		Date  *time.Time `json:"date,omitempty"`
		Email string     `json:"email,omitempty"`
		Name  string     `json:"name,omitempty"`
	} `json:"committer,omitempty"`
	Message string `json:"message,omitempty"`
	Tree    struct {
		Sha string `json:"sha,omitempty"`
		URL string `json:"url,omitempty"`
	} `json:"tree,omitempty"`
	URL string `json:"url,omitempty"`
}

// Commit is a representation of a full commit in git
type Commit struct {
	*hypermedia.HALResource

	Author      *User         `json:"author,omitempty"`
	CommentsURL string        `json:"comments_url,omitempty"`
	Commit      *CommitCommit `json:"commit,omitempty"`
	Committer   *User         `json:"committer,omitempty"`
	Files       []CommitFile  `json:"files,omitempty"`
	HtmlURL     string        `json:"html_url,omitempty"`
	Parents     []Commit      `json:"parents,omitempty"`
	Sha         string        `json:"sha,omitempty"`
	Stats       CommitStats   `json:"stats,omitempty"`
	URL         string        `json:"url,omitempty"`
}
