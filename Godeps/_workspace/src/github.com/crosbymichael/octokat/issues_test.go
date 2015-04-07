package octokat

import (
	"github.com/bmizerany/assert"
	"net/http"
	"testing"
	"time"
)

func TestIssues(t *testing.T) {
	setup()
	defer tearDown()

	mux.HandleFunc("/repos/octocat/Hello-World/issues", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		respondWith(w, loadFixture("issues.json"))
	})

	repo := Repo{UserName: "octocat", Name: "Hello-World"}
	issues, _ := client.Issues(repo, nil)

	assert.Equal(t, 1, len(issues))

	issue := issues[0]

	validateIssue(t, issue)

}

func TestIssue(t *testing.T) {
	setup()
	defer tearDown()

	mux.HandleFunc("/repos/octocat/Hello-World/issues/1347", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		respondWith(w, loadFixture("issue.json"))
	})

	repo := Repo{UserName: "octocat", Name: "Hello-World"}
	issue, _ := client.Issue(repo, 1347, nil)

	validateIssue(t, issue)
}

func validateIssue(t *testing.T, issue *Issue) {

	assert.Equal(t, "https://api.github.com/repos/octocat/Hello-World/issues/1347", issue.URL)
	assert.Equal(t, "https://github.com/octocat/Hello-World/issues/1347", issue.HTMLURL)
	assert.Equal(t, 1347, issue.Number)
	assert.Equal(t, "open", issue.State)
	assert.Equal(t, "Found a bug", issue.Title)
	assert.Equal(t, "I'm having a problem with this.", issue.Body)

	assert.Equal(t, "octocat", issue.User.Login)
	assert.Equal(t, 1, issue.User.ID)
	assert.Equal(t, "https://github.com/images/error/octocat_happy.gif", issue.User.AvatarURL)
	assert.Equal(t, "somehexcode", issue.User.GravatarID)
	assert.Equal(t, "https://api.github.com/users/octocat", issue.User.URL)

	assert.Equal(t, 1, len(issue.Labels))
	assert.Equal(t, "https://api.github.com/repos/octocat/Hello-World/labels/bug", issue.Labels[0].URL)
	assert.Equal(t, "bug", issue.Labels[0].Name)

	assert.Equal(t, "octocat", issue.Assignee.Login)
	assert.Equal(t, 1, issue.Assignee.ID)
	assert.Equal(t, "https://github.com/images/error/octocat_happy.gif", issue.Assignee.AvatarURL)
	assert.Equal(t, "somehexcode", issue.Assignee.GravatarID)
	assert.Equal(t, "https://api.github.com/users/octocat", issue.Assignee.URL)

	assert.Equal(t, "https://api.github.com/repos/octocat/Hello-World/milestones/1", issue.Milestone.URL)
	assert.Equal(t, 1, issue.Milestone.Number)
	assert.Equal(t, "open", issue.Milestone.State)
	assert.Equal(t, "v1.0", issue.Milestone.Title)
	assert.Equal(t, "", issue.Milestone.Description)

	assert.Equal(t, "octocat", issue.Milestone.Creator.Login)
	assert.Equal(t, 1, issue.Milestone.Creator.ID)
	assert.Equal(t, "https://github.com/images/error/octocat_happy.gif", issue.Milestone.Creator.AvatarURL)
	assert.Equal(t, "somehexcode", issue.Milestone.Creator.GravatarID)
	assert.Equal(t, "https://api.github.com/users/octocat", issue.Milestone.Creator.URL)

	assert.Equal(t, 4, issue.Milestone.OpenIssues)
	assert.Equal(t, 8, issue.Milestone.ClosedIssues)
	assert.Equal(t, "2011-04-10 20:09:31 +0000 UTC", issue.Milestone.CreatedAt.String())
	assert.Equal(t, (*time.Time)(nil), issue.Milestone.DueOn)

	assert.Equal(t, 0, issue.Comments)
	assert.Equal(t, "https://github.com/octocat/Hello-World/pull/1347", issue.PullRequest.HTMLURL)
	assert.Equal(t, "https://github.com/octocat/Hello-World/pull/1347.diff", issue.PullRequest.DiffURL)
	assert.Equal(t, "https://github.com/octocat/Hello-World/pull/1347.patch", issue.PullRequest.PatchURL)

	assert.Equal(t, (*time.Time)(nil), issue.ClosedAt)
	assert.Equal(t, "2011-04-22 13:33:48 +0000 UTC", issue.CreatedAt.String())
	assert.Equal(t, "2011-04-22 13:33:48 +0000 UTC", issue.UpdatedAt.String())

	// phew!
}
