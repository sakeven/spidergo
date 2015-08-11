package page

import (
	"net/url"
	"path/filepath"
)

// FixUri fixs uri like example.com/../../data.html, example.com/data//data.html
func FixUri(uri string) string {
	u, err := url.Parse(uri)
	if err != nil {
		return ""
	}

	u.Path = filepath.Clean(u.Path)

	return u.String()
}
