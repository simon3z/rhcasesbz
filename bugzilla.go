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
	JSONClient
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

	return &BugzillaClient{JSONClient{u, transport}}, nil
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
	bzResponse := new(struct {
		Bugs []BugzillaBug
	})

	err := b.JSONGetRequest(fmt.Sprintf("/rest/bug/%s", id), nil, &bzResponse)

	if err != nil {
		return nil, err
	}

	if len(bzResponse.Bugs) != 1 {
		return nil, ErrBugNotFound
	}

	return &bzResponse.Bugs[0], nil
}
