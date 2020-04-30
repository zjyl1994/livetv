package main

import (
	"bufio"
	"bytes"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

type Channel struct {
	Name string
	URL  string
}

func m3uHandler(c *gin.Context) {
	bTpl, err := ioutil.ReadFile(cfg.M3UTemplate)
	if err != nil {
		log.Println(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	t, err := template.New("m3u").Parse(string(bTpl))
	if err != nil {
		log.Println(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	m3uList, err := channelParser(cfg.ChannelFile)
	if err != nil {
		log.Println(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	for i, v := range m3uList {
		m3uList[i].URL = cfg.BaseURL + "/live.m3u8?url=" + url.QueryEscape(v.URL)
	}
	var m3u bytes.Buffer
	if err := t.Execute(&m3u, gin.H{"channels": m3uList}); err != nil {
		log.Println(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Data(http.StatusOK, "application/vnd.apple.mpegurl", m3u.Bytes())
}

func channelParser(channelFile string) ([]Channel, error) {
	bChannel, err := ioutil.ReadFile(channelFile)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(bytes.NewReader(bChannel))
	scanner.Split(bufio.ScanLines)
	var channels []Channel
	var currentChannel Channel
	for scanner.Scan() {
		l := scanner.Text()
		if strings.HasPrefix(l, "#") {
			currentChannel.Name = strings.TrimPrefix(l, "#")
		} else {
			currentChannel.URL = l
		}
		if currentChannel.Name != "" && currentChannel.URL != "" {
			channels = append(channels, currentChannel)
			currentChannel.Name = ""
			currentChannel.URL = ""
		}
	}
	return channels, nil
}
