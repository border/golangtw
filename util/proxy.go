package util

import (
    "strings"
    "fmt"
)

func CheckUrlPrefix(url string) string {
    if len(url) > 5 {
		prefix := strings.ToLower(url[0:5])
		if !strings.HasPrefix(prefix, "http") && !strings.HasPrefix(prefix, "http") {
			url = fmt.Sprintf("http://%s", url)
		}

	} else {
		url = fmt.Sprintf("http://%s", url)
	}
    return url
}
