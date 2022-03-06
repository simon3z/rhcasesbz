package rhcasesbz

import (
	"errors"
	"fmt"
	"net/http"
)

/* cspell:ignore rhcasesbz */

var ErrAPIRequestFailure = errors.New("api request failure")

type BasicAuthTransport struct {
	Parent   http.RoundTripper
	Username string
	Password string
}

func NewBasicAuthTransport(p http.RoundTripper, username, password string) *BasicAuthTransport {
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
