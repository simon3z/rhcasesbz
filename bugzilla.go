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
	Host   string
	ApiKey string
}

type BugzillaBug struct {
	Summary       string
	Status        string
	TargetRelease []string `json:"target_release"`
}

func NewBugzillaClient(apikey string) *BugzillaClient {
	return &BugzillaClient{"bugzilla.redhat.com", apikey}
}

func (b *BugzillaClient) FetchBug(id string) (*BugzillaBug, error) {
	u := url.URL{Scheme: "https", Host: b.Host, Path: fmt.Sprintf("/rest/bug/%s", id)}

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
