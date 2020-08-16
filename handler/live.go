package handler

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/zjyl1994/livetv/global"
	"github.com/zjyl1994/livetv/service"
	"github.com/zjyl1994/livetv/util"
)

func M3UHandler(c *gin.Context) {
	content, err := service.M3UGenerate()
	if err != nil {
		log.Println(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Data(http.StatusOK, "application/vnd.apple.mpegurl", []byte(content))
}

func LiveHandler(c *gin.Context) {
	channelNumber := c.Query("channel")
	if channelNumber == "" {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	channelInfo, err := service.GetChannel(util.String2Uint(channelNumber))
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			c.AbortWithStatus(http.StatusNotFound)
		} else {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}
	baseUrl, err := service.GetConfig("base_url")
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	liveM3U8, err := service.GetYoutubeLiveM3U8(channelInfo.URL)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if channelInfo.Proxy {
		liveProxyM3U8 := baseUrl + "/p/live.m3u8?url=" + url.QueryEscape(liveM3U8)
		c.Redirect(http.StatusTemporaryRedirect, liveProxyM3U8)
	} else {
		c.Redirect(http.StatusTemporaryRedirect, liveM3U8)
	}
}

func TsProxyHandler(c *gin.Context) {
	remoteURL := c.Query("url")
	client := http.Client{Timeout: global.HttpClientTimeout}
	resp, err := client.Get(remoteURL)
	if err != nil {
		log.Println(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer resp.Body.Close()
	c.DataFromReader(http.StatusOK, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)
}

func M3U8ProxyHandler(c *gin.Context) {
	m3u8URL := c.Query("url")
	if m3u8URL == "" {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	baseURL, err := service.GetConfig("base_url")
	if err != nil {
		log.Println(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	client := http.Client{Timeout: global.HttpClientTimeout}
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
	processedBody := service.M3U8Process(bodyString, baseURL+"/p/live.ts?url=")
	c.Data(http.StatusOK, resp.Header.Get("Content-Type"), []byte(processedBody))
}
