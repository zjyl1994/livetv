package global

import (
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
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
	M3U8Cache         = cache.New(3*time.Second, 10*time.Second)
)
