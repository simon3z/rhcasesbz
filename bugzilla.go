package rhcasesbz

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

/* cspell:ignore rhcasesbz */

var ErrBugNotFound = errors.New("bugzilla: bug not found")

type BugzillaClient struct {
	BaseURL *url.URL
	ApiKey  string
}

type BugzillaBug struct {
	Summary       string
	Status        string
	TargetRelease []string `json:"target_release"`
}

func NewBugzillaClient(baseURL string, apikey string) (*BugzillaClient, error) {
	u, err := url.Parse(baseURL)

	if err != nil {
		return nil, err
	}

	return &BugzillaClient{u, apikey}, nil
}

func (b *BugzillaClient) FetchBug(id string) (*BugzillaBug, error) {
	u := url.URL{
		Scheme: b.BaseURL.Scheme,
		Host:   b.BaseURL.Host,
		Path:   fmt.Sprintf("%s/rest/bug/%s", b.BaseURL.Path, id),
	}

	request, err := http.NewRequest("GET", u.String(), nil)

	if err != nil {
		return nil, err
	}

	client := &http.Client{Transport: NewBearerAuthTransport(NewJSONTransport(http.DefaultTransport), b.ApiKey)}
	r, err := client.Do(request)

	if err != nil {
		return nil, err
	}

	defer r.Body.Close()

	bzResponse := new(struct {
		Bugs []BugzillaBug
	})

	err = json.NewDecoder(r.Body).Decode(bzResponse)

	if err != nil {
		return nil, err
	}

	if len(bzResponse.Bugs) != 1 {
		return nil, ErrBugNotFound
	}

	return &bzResponse.Bugs[0], nil
}
