package octokit

// Organization is a representation of an organization on GitHub, containing
// all identifying information related to the specific organization.
type Organization struct {
	AvatarURL        string `json:"avatar_url,omitempty"`
	PublicMembersURL string `json:"public_member_url,omitempty"`
	MembersURL       string `json:"members_url,omitempty"`
	EventsURL        string `json:"events_url,omitempty"`
	ReposURL         string `json:"repos_url,omitempty"`
	URL              string `json:"url,omitempty"`
	ID               int    `json:"id,omitempty"`
	Login            string `json:"login,omitempty"`
}
