package service

import (
	"log"
	"regexp"
	"time"

	"github.com/zjyl1994/livetv/global"
	"github.com/zjyl1994/livetv/util"
)

func LoadChannelCache() {
	channels, err := GetAllChannel()
	if err != nil {
		log.Println(err)
		return
	}
	for _, v := range channels {
		log.Println("caching", v.URL)
		liveURL, err := RealGetYoutubeLiveM3U8(v.URL)
		if err != nil {
			log.Println(err)
			return
		}
		global.URLCache.Store(v.URL, liveURL)
		log.Println(v.URL, "cached")
	}
}

func UpdateURLCache() {
	channels, err := GetAllChannel()
	if err != nil {
		log.Println(err)
		return
	}
	for _, v := range channels {
		log.Println("caching", v.URL)
		liveURL, err := RealGetYoutubeLiveM3U8(v.URL)
		if err != nil {
			log.Println(err)
		} else {
			global.URLCache.Store(v.URL, liveURL)
			log.Println(v.URL, "cached")
		}
	}
	global.URLCache.Range(func(k, v interface{}) bool {
		value := v.(string)
		regex := regexp.MustCompile(`/expire/(\d+)/`)
		matched := regex.FindStringSubmatch(value)
		if len(matched) < 2 {
			global.URLCache.Delete(k)
		}
		expireTime := time.Unix(util.String2Int64(matched[1]), 0)
		if time.Now().After(expireTime) {
			global.URLCache.Delete(k)
		}
		return true
	})
}
