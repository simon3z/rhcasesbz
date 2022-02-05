package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/simon3z/rhcasesbz"
)

/* cspell:ignore rhuser rhpass rhcasesbz bugzillas rhbzkey */

func main() {
	rhuser := os.Getenv("RHUSER")
	rhpass := os.Getenv("RHPASS")
	rhbzkey := os.Getenv("RHBZKEY")

	h := rhcasesbz.NewHydraClient(rhuser, rhpass)
	b := rhcasesbz.NewBugzillaClient(rhbzkey)

	r := csv.NewReader(os.Stdin)
	r.Comma = '\t'

	w := csv.NewWriter(os.Stdout)
	w.Comma = '\t'

	accountMap := map[string]string{}

	for {
		record, err := r.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			panic(err)
		}

		c, err := h.FetchCase(record[0])

		if err != nil {
			panic(err)
		}

		if _, ok := accountMap[c.Account]; !ok {
			a, err := h.FetchAccount(c.Account)

			if err != nil {
				panic(err)
			}

			accountMap[c.Account] = a.Name
		}

		e := []string{Hyperlink(record[0], c.Link()), accountMap[c.Account]}

		if len(record) > 1 {
			e = append(e, record[1:]...)
		}

		e = append(e, PreviewString(c.Summary, 40), fmt.Sprintf("%s %s", c.Product, c.Version), c.Status, c.Severity, FormatDate(c.LastModified))

		if len(c.Bugzillas.Items) == 0 {
			w.Write(e)
		} else {
			for _, i := range c.Bugzillas.Items {
				u, err := b.FetchBug(i.ID)

				if err != nil {
					panic(err)
				}

				z := append(e, Hyperlink(fmt.Sprintf("BZ#%s", i.ID), i.Link), PreviewString(u.Summary, 40), u.Status, strings.Join(u.TargetRelease, ","))

				w.Write(z)
			}
		}

		w.Flush()
	}
}