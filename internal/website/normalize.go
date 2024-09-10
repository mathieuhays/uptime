package website

import "strings"

func NormalizeURL(url string) string {
	url = strings.ToLower(url)
	url = strings.Split(url, "?")[0]
	url = strings.Split(url, "#")[0]
	url = strings.TrimSuffix(url, "/")

	return url
}
