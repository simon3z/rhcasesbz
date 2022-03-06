package rhcasesbz

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

/* cspell:ignore rhcasesbz zstream */

var ErrBugNotFound = errors.New("bugzilla: bug not found")

type BugzillaClient struct {
	BaseURL   *url.URL
	Transport http.RoundTripper
}

type BugzillaBug struct {
	Summary               string
	Status                string
	Product               string
	TargetRelease         []string  `json:"target_release"`
	ZStreamTarget         string    `json:"cf_zstream_target_release"`
	InternalTargetRelease string    `json:"cf_internal_target_release"`
	LastChangeTime        time.Time `json:"-"`
}

func NewBugzillaClient(baseURL string, transport http.RoundTripper) (*BugzillaClient, error) {
	u, err := url.Parse(baseURL)

	if err != nil {
		return nil, err
	}

	return &BugzillaClient{u, transport}, nil
}

func (c *BugzillaBug) UnmarshalJSON(data []byte) error {
	type BugzillaBugJSON BugzillaBug

	d := new(struct {
		BugzillaBugJSON
		LastChangeTimeString string `json:"last_change_time"`
	})

	err := json.Unmarshal(data, d)

	if err != nil {
		return err
	}

	*c = (BugzillaBug)(d.BugzillaBugJSON)
	c.LastChangeTime, err = time.Parse(time.RFC3339, d.LastChangeTimeString)

	if err != nil {
		return err
	}

	return nil
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

	client := &http.Client{Transport: NewJSONTransport(b.Transport)}
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
