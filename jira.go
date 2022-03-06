package rhcasesbz

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type JiraClient struct {
	BaseURL   *url.URL
	Transport http.RoundTripper
}

func NewJiraClient(baseURL string, transport http.RoundTripper) (*JiraClient, error) {
	u, err := url.Parse(baseURL)

	if err != nil {
		return nil, err
	}

	return &JiraClient{u, transport}, nil
}

type JiraIssue struct {
	Id   string
	Key  string
	Link string
}

func (j *JiraClient) FindIssues(jql string) ([]JiraIssue, error) {
	u := url.URL{
		Scheme: j.BaseURL.Scheme,
		Host:   j.BaseURL.Host,
		Path:   fmt.Sprintf("%s/rest/api/2/search", j.BaseURL.Path),
	}

	q := u.Query()
	q.Add("jql", jql)

	u.RawQuery = q.Encode()

	request, err := http.NewRequest("GET", u.String(), nil)

	if err != nil {
		return nil, err
	}

	client := &http.Client{Transport: NewJSONTransport(j.Transport)}
	r, err := client.Do(request)

	if err != nil {
		return nil, err
	}

	defer r.Body.Close()

	jiraResponse := new(struct {
		Issues []JiraIssue
	})

	err = json.NewDecoder(r.Body).Decode(jiraResponse)

	if err != nil {
		return nil, err
	}

	for i := range jiraResponse.Issues {
		jiraResponse.Issues[i].Link = fmt.Sprintf("%s/browse/%s", j.BaseURL.String(), jiraResponse.Issues[i].Key)
	}

	return jiraResponse.Issues, nil
}

func JQLEscapeString(s string) string {
	return fmt.Sprintf("\"%s\"", strings.ReplaceAll(s, "\"", "\\\""))
}
