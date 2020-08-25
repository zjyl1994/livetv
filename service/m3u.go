package service

import (
	"log"
	"strconv"
	"strings"
)

func M3UGenerate() (string, error) {
	baseUrl, err := GetConfig("base_url")
	if err != nil {
		log.Println(err)
		return "", err
	}
	channels, err := GetAllChannel()
	if err != nil {
		log.Println(err)
		return "", err
	}
	var m3u strings.Builder
	m3u.WriteString("#EXTM3U\n")
	for _, v := range channels {
		m3u.WriteString("#EXTINF:-1,")
		m3u.WriteString(v.Name)
		m3u.WriteString("\n")
		m3u.WriteString(baseUrl)
		m3u.WriteString("/live.m3u8?c=")
		m3u.WriteString(strconv.Itoa(int(v.ID)))
		m3u.WriteString("\n")
	}
	return m3u.String(), nil
}
