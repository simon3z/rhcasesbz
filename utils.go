package rhcasesbz

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

/* cspell:ignore rhcasesbz */

var ErrAPIRequestFailure = errors.New("api request failure")

type BasicAuthTransport struct {
	Parent   http.RoundTripper
	Username string
	Password string
}

func NewBasicAuthTransport(p http.RoundTripper, username, password string) *BasicAuthTransport {
	if p == nil {
		p = http.DefaultTransport
	}
	return &BasicAuthTransport{p, username, password}
}

func (t *BasicAuthTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r.SetBasicAuth(t.Username, t.Password)
	return t.Parent.RoundTrip(r)
}

type BearerAuthTransport struct {
	Parent http.RoundTripper
	Token  string
}

func NewBearerAuthTransport(p http.RoundTripper, token string) *BearerAuthTransport {
	if p == nil {
		p = http.DefaultTransport
	}
	return &BearerAuthTransport{p, token}
}

func (t *BearerAuthTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t.Token))
	return t.Parent.RoundTrip(r)
}

type JSONTransport struct {
	Parent http.RoundTripper
}

func NewJSONTransport(p http.RoundTripper) *JSONTransport {
	if p == nil {
		p = http.DefaultTransport
	}
	return &JSONTransport{p}
}

func (t *JSONTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Set("Accept", "application/json")
	r.Header.Set("Content-Type", "application/json")

	response, err := t.Parent.RoundTrip(r)

	if err != nil {
		return response, err
	}

	if response.StatusCode != http.StatusOK {
		return response, fmt.Errorf("%w: %s", ErrAPIRequestFailure, http.StatusText(response.StatusCode))
	}

	return response, err
}

type JSONClient struct {
	BaseURL   *url.URL
	Transport http.RoundTripper
}

func (c *JSONClient) JSONGetRequest(path string, query *url.Values, v interface{}) error {
	u := url.URL{
		Scheme: c.BaseURL.Scheme,
		Host:   c.BaseURL.Host,
		Path:   fmt.Sprintf("%s/%s", c.BaseURL.Path, path),
	}

	if query != nil {
		u.RawQuery = query.Encode()
	}

	request, err := http.NewRequest("GET", u.String(), nil)

	if err != nil {
		return err
	}

	client := &http.Client{Transport: NewJSONTransport(c.Transport)}
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
