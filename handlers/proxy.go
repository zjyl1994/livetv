package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zjyl1994/livetv/services/proxy"
)

func TsProxyHandler(c *gin.Context) {
	strClipId := c.Query("id")
	if strClipId == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	contentType, data, err := proxy.TSProxy(strClipId)
	if err != nil {
		if err == proxy.ErrNotFound {
			c.AbortWithStatus(http.StatusNotFound)
		} else {
			c.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}
	c.Data(http.StatusOK, contentType, data)
}

func M3U8ProxyHandler(c *gin.Context) {
	strChannelNumber := c.Query("id")
	if strChannelNumber == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	contentType, data, err := proxy.M3U8Proxy(strChannelNumber)
	if err != nil {
		if err == proxy.ErrNotFound {
			c.AbortWithStatus(http.StatusNotFound)
		} else {
			c.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}
	c.Data(http.StatusOK, contentType, data)
}
