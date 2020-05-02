package main

import (
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var urlCache sync.Map

func loadChannelCache() {
	channels, err := channelParser(cfg.ChannelFile)
	if err != nil {
		log.Println(err)
		return
	}
	for _, v := range channels {
		log.Println("caching", v.URL)
		liveURL, err := realGetYoutubeLiveM3U8(v.URL)
		if err != nil {
			log.Println(err)
			return
		}
		urlCache.Store(v.URL, liveURL)
		log.Println(v.URL, "cached")
	}
}

func updateURLCache() {
	channels, err := channelParser(cfg.ChannelFile)
	if err != nil {
		log.Println(err)
		return
	}
	for _, v := range channels {
		log.Println("caching", v.URL)
		liveURL, err := realGetYoutubeLiveM3U8(v.URL)
		if err != nil {
			log.Println(err)
		} else {
			urlCache.Store(v.URL, liveURL)
			log.Println(v.URL, "cached")
		}
	}
	urlCache.Range(func(k, v interface{}) bool {
		value := v.(string)
		regex := regexp.MustCompile(`/expire/(\d+)/`)
		matched := regex.FindStringSubmatch(value)
		if len(matched) < 2 {
			urlCache.Delete(k)
		}
		expireTime := time.Unix(string2Int64(matched[1]), 0)
		if time.Now().After(expireTime) {
			urlCache.Delete(k)
		}
		return true
	})
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func string2Int64(s string) int64 {
	i64, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	} else {
		return i64
	}
}

func cacheHandler(c *gin.Context) {
	var sb strings.Builder
	urlCache.Range(func(k, v interface{}) bool {
		sb.WriteString(k.(string))
		sb.WriteString(" => ")
		sb.WriteString(v.(string))
		sb.WriteString("\n")
		return true
	})
	c.Data(http.StatusOK, "text/plain", []byte(sb.String()))
}
