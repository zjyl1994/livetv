package service

import (
	"bufio"
	"net/url"
	"strings"
)

func M3U8Process(data string, prefixURL string) string {
	var sb strings.Builder
	scanner := bufio.NewScanner(strings.NewReader(data))
	for scanner.Scan() {
		l := scanner.Text()
		if strings.HasPrefix(l, "#") {
			sb.WriteString(l)
		} else {
			sb.WriteString(prefixURL)
			sb.WriteString(url.QueryEscape(l))
		}
		sb.WriteString("\n")
	}
	return sb.String()
}
