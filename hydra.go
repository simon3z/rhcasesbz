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
	Host     string
	Username string
	Password string
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

func NewHydraClient(username, password string) *HydraClient {
	return &HydraClient{"api.access.redhat.com", username, password}
}

func (h *HydraClient) getRequest(path string, v interface{}) error {
	u := url.URL{Scheme: "https", Host: h.Host, Path: path}

	request, err := http.NewRequest("GET", u.String(), nil)

	if err != nil {
		return err
	}

	client := &http.Client{Transport: NewBasicAuthTransport(NewJSONTransport(http.DefaultTransport), h.Username, h.Password)}
	r, err := client.Do(request)

	if err != nil {
		return err
	}

	defer r.Body.Close()

	err = json.NewDecoder(r.Body).Decode(v)

	if err != nil {
		return err
	}

	return nil
}

func (h *HydraClient) FetchCase(id string) (*HydraCase, error) {
	c := new(HydraCase)

	err := h.getRequest(fmt.Sprintf("/rs/cases/%s", id), c)

	return c, err
}

func (c *HydraCase) Link() string {
	return fmt.Sprintf("https://access.redhat.com/support/cases/#/case/%s", c.ID)
}

func (h *HydraClient) FetchAccount(id string) (*HydraAccount, error) {
	a := new(HydraAccount)

	err := h.getRequest(fmt.Sprintf("/rs/accounts/%s", id), a)

	return a, err
}
