package global

import (
	"sync"
	"time"
)

var DefaultConfigValue = map[string]string{
	"ytdl_cmd":  "youtube-dl",
	"ytdl_args": "-f best -g {url}",
	"base_url":  "http://127.0.0.1:9000",
}

var (
	HttpClientTimeout = 10 * time.Second
	VersionString     = "LiveTV! v0.0.2"
	ConfigCache       sync.Map
	URLCache          sync.Map
)
