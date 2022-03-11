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

type HydraBug struct {
	ID      string `json:"bugzillaNumber"`
	Link    string `json:"bugzillaLink"`
	Summary string
	Status  string
}

type HydraCase struct {
	ID           string `json:"caseNumber"`
	Link         string `json:"-"`
	Account      string `json:"accountNumberRef"`
	Summary      string
	Status       string
	Product      string
	Version      string
	Severity     string
	Escalated    bool `json:"customerEscalation"`
	Bugzillas    []HydraBug
	LastModified time.Time `json:"-"`
}

func (c *HydraCase) UnmarshalJSON(data []byte) error {
	type HydraCaseJSON HydraCase

	d := new(struct {
		HydraCaseJSON
		LastModifiedString string `json:"lastModifiedDate"`
	})

	err := json.Unmarshal(data, d)

	if err != nil {
		return err
	}

	*c = (HydraCase)(d.HydraCaseJSON)
	c.LastModified, err = time.Parse(time.RFC3339, d.LastModifiedString)

	if err != nil {
		return err
	}

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

	err := h.JSONGetRequest(fmt.Sprintf("/hydra/rest/v1/cases/%s", id), nil, c)

	if err != nil {
		return nil, err
	}

	c.Link = fmt.Sprintf("%s/support/cases/#/case/%s", h.BaseURL.String(), c.ID)

	return c, nil
}

func (h *HydraClient) FetchAccount(id string) (*HydraAccount, error) {
	a := new(HydraAccount)

	err := h.JSONGetRequest(fmt.Sprintf("/hydra/rest/v1/accounts/%s", id), nil, a)

	if err != nil {
		return nil, err
	}

	return a, nil
}
