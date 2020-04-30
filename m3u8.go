package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

func liveHandler(c *gin.Context) {
	liveURL := c.Query("url")
	if liveURL == "" {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	liveM3U8, err := getYoutubeLiveM3U8(liveURL)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if cfg.ProxyStream {
		liveProxyM3U8 := cfg.BaseURL + "/p/live.m3u8?url=" + url.QueryEscape(liveM3U8)
		c.Redirect(http.StatusTemporaryRedirect, liveProxyM3U8)
	} else {
		c.Redirect(http.StatusTemporaryRedirect, liveM3U8)
	}
}

func tsProxyHandler(c *gin.Context) {
	remoteURL := c.Query("url")
	client := http.Client{Timeout: httpClientTimeout}
	resp, err := client.Get(remoteURL)
	if err != nil {
		log.Println(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer resp.Body.Close()
	c.DataFromReader(http.StatusOK, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)
}

func m3u8ProxyHandler(c *gin.Context) {
	m3u8URL := c.Query("url")
	if m3u8URL == "" {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	client := http.Client{Timeout: httpClientTimeout}
	resp, err := client.Get(m3u8URL)
	if err != nil {
		log.Println(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	bodyString := string(bodyBytes)
	processedBody := m3u8Proc(bodyString, cfg.BaseURL+"/p/live.ts?url=")
	c.Data(200, resp.Header.Get("Content-Type"), []byte(processedBody))
}

func m3u8Proc(data string, prefixURL string) string {
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
