package rhcasesbz

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

/* cspell:ignore rhcasesbz bugzilla bugzillas  bzapikey */

type HydraClient struct {
	JSONClient
}

type HydraAccount struct {
	Name string
}

type HydraBugsList struct {
	Items []HydraBug `json:"bugzilla"`
}

type HydraBug struct {
	ID      string `json:"bugzilla_number"`
	Link    string `json:"resource_view_uri"`
	Summary string
	Release string
	Status  string
}

type HydraCase struct {
	ID           string `json:"case_number"`
	Account      string `json:"account_number"`
	Summary      string
	Status       string
	Product      string
	Version      string
	Severity     string
	Escalated    bool
	Bugzillas    HydraBugsList
	LastModified time.Time `json:"-"`
}

func (c *HydraCase) UnmarshalJSON(data []byte) error {
	type HydraCaseJSON HydraCase

	d := new(struct {
		HydraCaseJSON
		LastModifiedInt int64 `json:"last_modified_date"`
	})

	err := json.Unmarshal(data, d)

	if err != nil {
		return err
	}

	*c = (HydraCase)(d.HydraCaseJSON)
	c.LastModified = time.Unix(d.LastModifiedInt/1000, d.LastModifiedInt%1000)

	return nil
}

func NewHydraClient(baseURL string, transport http.RoundTripper) (*HydraClient, error) {
	u, err := url.Parse(baseURL)

	if err != nil {
		return nil, err
	}

	return &HydraClient{JSONClient{u, transport}}, nil
}

func (h *HydraClient) FetchCase(id string) (*HydraCase, error) {
	c := new(HydraCase)

	err := h.JSONGetRequest(fmt.Sprintf("/rs/cases/%s", id), nil, c)

	if err != nil {
		return nil, err
	}

	return c, nil
}

func (h *HydraClient) FetchAccount(id string) (*HydraAccount, error) {
	a := new(HydraAccount)

	err := h.JSONGetRequest(fmt.Sprintf("/rs/accounts/%s", id), nil, a)

	if err != nil {
		return nil, err
	}

	return a, nil
}
