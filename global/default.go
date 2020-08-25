package global

import (
	"sync"
	"time"
)

var defaultConfigValue = map[string]string{
	"ytdl_cmd":  "youtube-dl",
	"ytdl_args": "-f best -g {url}",
	"base_url":  "http://127.0.0.1:9000",
	"password":  "password",
}

var (
	HttpClientTimeout = 30 * time.Second
	ConfigCache       sync.Map
	URLCache          sync.Map
)
