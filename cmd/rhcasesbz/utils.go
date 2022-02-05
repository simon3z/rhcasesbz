package main

import (
	"bytes"
	"fmt"
	"time"
)

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
