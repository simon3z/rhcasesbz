package rhcasesbz

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type JiraClient struct {
	JSONClient
}

func NewJiraClient(baseURL string, transport http.RoundTripper) (*JiraClient, error) {
	u, err := url.Parse(baseURL)

	if err != nil {
		return nil, err
	}

	return &JiraClient{JSONClient{u, transport}}, nil
}

type JiraIssue struct {
	Id   string
	Key  string
	Link string
}

func (j *JiraClient) FindIssues(jql string) ([]JiraIssue, error) {
	q := &url.Values{"jql": []string{jql}}

	jiraResponse := new(struct {
		Issues []JiraIssue
	})

	err := j.JSONGetRequest("/rest/api/2/search", q, jiraResponse)

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
