package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zjyl1994/livetv/service"
)

func IndexHandler(c *gin.Context) {
	baseUrl, err := service.GetConfig("base_url")
	if err != nil {
		log.Println(err.Error())
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"ErrMsg": err.Error(),
		})
		return
	}
	channelModels, err := service.GetAllChannel()
	if err != nil {
		log.Println(err.Error())
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"ErrMsg": err.Error(),
		})
		return
	}
	channels := make([]Channel, len(channelModels))
	for i, v := range channelModels {
		channels[i] = Channel{
			ID:   v.ID,
			Name: v.Name,
			URL:  v.URL,
			M3U8: baseUrl + "/live.m3u8?channel=" + strconv.Itoa(int(v.ID)),
		}
	}
	c.HTML(http.StatusOK, "index.html", gin.H{
		"Channels": channels,
	})
}
