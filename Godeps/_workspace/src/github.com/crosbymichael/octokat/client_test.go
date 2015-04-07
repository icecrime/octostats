package octokat

import (
	"bytes"
	"github.com/bmizerany/assert"
	"net/http"
	"testing"
)

func TestGet(t *testing.T) {
	setup()
	defer tearDown()

	mux.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", MediaType)
		testHeader(t, r, "User-Agent", UserAgent)
		testHeader(t, r, "Content-Type", DefaultContentType)
		respondWith(w, "ok")
	})

	client.get("foo", nil)
}

func TestPost(t *testing.T) {
	setup()
	defer tearDown()

	content := "content"
	mux.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testHeader(t, r, "Accept", MediaType)
		testHeader(t, r, "User-Agent", UserAgent)
		testHeader(t, r, "Content-Type", DefaultContentType)
		testBody(t, r, content)
		respondWith(w, "ok")
	})

	client.post("foo", nil, bytes.NewBufferString(content))
}

func TestJSONPost(t *testing.T) {
	setup()
	defer tearDown()

	mux.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testHeader(t, r, "Accept", MediaType)
		testHeader(t, r, "User-Agent", UserAgent)
		testHeader(t, r, "Content-Type", DefaultContentType)
		testBody(t, r, "")
		respondWith(w, `{"ok": "foo"}`)
	})

	m := make(map[string]interface{})
	client.jsonPost("foo", nil, &m)

	assert.Equal(t, "foo", m["ok"])
}

func TestBuildURL(t *testing.T) {
	url, _ := client.buildURL("https://api.github.com")
	assert.Equal(t, "https://api.github.com", url.String())

	url, _ = client.buildURL("repos")
	assert.Equal(t, testURLOf("repos"), url.String())
}
