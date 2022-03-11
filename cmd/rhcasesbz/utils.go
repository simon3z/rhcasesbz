package main

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/simon3z/rhcasesbz"
)

/* cspell:ignore rhcasesbz kubernetes rhacm rhocs rhodf rhelplan ocpbugsm */

func Hyperlink(text, url string) string {
	return fmt.Sprintf("=HYPERLINK(\"%s\",\"%s\")", url, text)
}

func PreviewString(text string, length int) string {
	runes := bytes.Runes([]byte(text))

	if len(runes) > length {
		return string(runes[:length]) + "…"
	}

	return string(runes)
}

func FormatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

var productAcronymsMap = map[string]string{
	"OpenShift Container Platform":                       "OCP",
	"Red Hat Enterprise Linux":                           "RHEL",
	"Red Hat Advanced Cluster Management for Kubernetes": "RHACM",
	"Red Hat OpenShift Container Storage":                "RHOCS",
	"Red Hat OpenShift Data Foundation":                  "RHODF",
	"Red Hat Enterprise Linux Fast Datapath":             "RHEL-FD",
}

var productVersionRE = regexp.MustCompile(` [0-9]*$`)

func ShortenProductRelease(product, release string, l bool) string {
	p := productVersionRE.ReplaceAllString(product, "")

	p, ok := productAcronymsMap[p]

	if !ok {
		p = product
	}

	s := " "

	if l {
		p = strings.ToLower(p)
		s = "-"
	}

	if release == "" {
		return p
	}

	return fmt.Sprintf("%s%s%s", p, s, release)
}

func GetBugTargetRelease(b *rhcasesbz.BugzillaBug) string {
	if b.ZStreamTarget != "" && b.ZStreamTarget != "---" {
		return b.ZStreamTarget
	}

	if b.InternalTargetRelease != "" && b.InternalTargetRelease != "---" {
		return b.InternalTargetRelease
	}

	r := []string{}

	for _, i := range b.TargetRelease {
		if i != "" && i != "---" {
			r = append(r, i)
		}
	}

	if len(r) == 0 {
		return ""
	}

	return strings.Join(r, ",")
}

var ErrJiraBugIssueNotFound = errors.New("jira: bug issue not found")
var ErrJiraMultipleBugIssuesFound = errors.New("jira: multiple issues found for single bugs")

func FindJiraIssueByBugzillaID(j *rhcasesbz.JiraClient, id string) (*rhcasesbz.JiraIssue, error) {
	l, err := j.FindIssues(fmt.Sprintf("cf[12316840]=%s AND project IN ('RHELPLAN', 'OCPBUGSM')", rhcasesbz.JQLEscapeString(id)))

	if err != nil {
		return nil, err
	}

	switch {
	case len(l) < 1:
		return nil, ErrJiraBugIssueNotFound
	case len(l) > 1:
		return nil, ErrJiraMultipleBugIssuesFound
	}

	return &l[0], nil
}
