package octokit

import (
	"net/url"
	"time"

	"github.com/jingweno/go-sawyer/hypermedia"
)

// ReleasesURL is a template for accessing releases in a particular repository
// for a particular owner that can be expanded to a full address.
var ReleasesURL = Hyperlink("repos/{owner}/{repo}/releases{/id}")

// Release is a representation of a release on GitHub. Published releases are
// available to everyone.
type Release struct {
	*hypermedia.HALResource

	ID              int        `json:"id,omitempty"`
	URL             string     `json:"url,omitempty"`
	HTMLURL         string     `json:"html_url,omitempty"`
	AssetsURL       string     `json:"assets_url,omitempty"`
	UploadURL       Hyperlink  `json:"upload_url,omitempty"`
	TagName         string     `json:"tag_name,omitempty"`
	TargetCommitish string     `json:"target_commitish,omitempty"`
	Name            string     `json:"name,omitempty"`
	Body            string     `json:"body,omitempty"`
	Draft           bool       `json:"draft,omitempty"`
	Prerelease      bool       `json:"prerelease,omitempty"`
	CreatedAt       *time.Time `json:"created_at,omitempty"`
	PublishedAt     *time.Time `json:"published_at,omitempty"`
	Assets          []Asset    `json:"assets,omitempty"`
}

// Asset represents a piece of content produced and associated with a given
// released that may be downloaded
type Asset struct {
	ID            int        `json:"id,omitempty"`
	Name          string     `json:"name,omitempty"`
	Label         string     `json:"label,omitempty"`
	ContentType   string     `json:"content_type,omitempty"`
	State         string     `json:"state,omitempty"`
	Size          int        `json:"size,omitempty"`
	DownloadCount int        `json:"download_count,omitempty"`
	URL           string     `json:"url,omitempty"`
	CreatedAt     *time.Time `json:"created_at,omitempty"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty"`
}

// Releases creates a ReleasesService with a base url
func (c *Client) Releases(url *url.URL) (releases *ReleasesService) {
	releases = &ReleasesService{client: c, URL: url}
	return
}

// ReleasesService is a service providing access to releases from a particular url
type ReleasesService struct {
	client *Client
	URL    *url.URL
}

// One gets a specific release based on the url of the service
func (r *ReleasesService) All() (releases []Release, result *Result) {
	result = r.client.get(r.URL, &releases)
	return
}

// Create posts a new release based on the relase parameters to the releases service url
func (r *ReleasesService) Create(params interface{}) (release *Release, result *Result) {
	result = r.client.post(r.URL, params, &release)
	return
}

// Update modifies a release based on the release parameters on the service url
func (r *ReleasesService) Update(params interface{}) (release *Release, result *Result) {
	result = r.client.patch(r.URL, params, &release)
	return
}

// ReleaseParams represent the parameters used to create or update a release
type ReleaseParams struct {
	TagName         string `json:"tag_name,omitempty"`
	TargetCommitish string `json:"target_commitish,omitempty"`
	Name            string `json:"name,omitempty"`
	Body            string `json:"body,omitempty"`
	Draft           bool   `json:"draft,omitempty"`
	Prerelease      bool   `json:"prerelease,omitempty"`
}
