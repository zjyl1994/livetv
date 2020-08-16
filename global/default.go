package global

import "time"

var DefaultConfigValue = map[string]string{
	"ytdl_cmd":     "youtube-dl",
	"ytdl_args":    "-f best -g {url}",
	"preload_cron": "30 * * * *",
	"base_url":     "http://127.0.0.1:9000",
}

var (
	HttpClientTimeout = 10 * time.Second
)
