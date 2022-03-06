package main

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"time"
)

/* cspell:ignore kubernetes rhacm rhocs rhodf */

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
