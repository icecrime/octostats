package octokit

import (
	"net/url"
)

// EmojisURL is an address for accessing the emojis available for use on GitHub.
var EmojisURL = Hyperlink("/emojis")

// Emojis creates an EmojisService with a base url
func (c *Client) Emojis(url *url.URL) (emojis *EmojisService) {
	emojis = &EmojisService{client: c, URL: url}
	return
}

// EmojisService is a service providing access to all the emojis available from a
// particular url
type EmojisService struct {
	client *Client
	URL    *url.URL
}

// All gets a list all the available emoji paths from the service
func (s *EmojisService) All() (emojis map[string]string, result *Result) {
	result = s.client.get(s.URL, &emojis)
	return
}
