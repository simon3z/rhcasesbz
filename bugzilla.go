package rhcasesbz

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

/* cspell:ignore rhcasesbz bzapikey */

var ErrBugNotFound = errors.New("bugzilla: bug not found")

type BugzillaClient struct {
	Host     string
	BZApiKey string
}

type BugzillaBug struct {
	Summary       string
	Status        string
	TargetRelease []string `json:"target_release"`
}

func NewBugzillaClient(bzapikey string) *BugzillaClient {
	return &BugzillaClient{"bugzilla.redhat.com", bzapikey}
}

func (b *BugzillaClient) FetchBug(id string) (*BugzillaBug, error) {
	u := url.URL{Scheme: "https", Host: b.Host, Path: fmt.Sprintf("/rest/bug/%s", id)}

	request, err := http.NewRequest("GET", u.String(), nil)

	if err != nil {
		return nil, err
	}

	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", b.BZApiKey))

	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
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
