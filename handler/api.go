package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zjyl1994/livetv/model"
	"github.com/zjyl1994/livetv/service"
	"github.com/zjyl1994/livetv/util"
)

func ChannelIndexHandler(c *gin.Context) {
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
			ID:    v.ID,
			Name:  v.Name,
			URL:   v.URL,
			M3U8:  baseUrl + "/live.m3u8?channel=" + strconv.Itoa(int(v.ID)),
			Proxy: v.Proxy,
		}
	}
	c.HTML(http.StatusOK, "index.html", gin.H{
		"Channels": channels,
	})
}

func NewChannelHandler(c *gin.Context) {
	chName := c.PostForm("name")
	chURL := c.PostForm("url")
	chProxy := c.PostForm("proxy") != ""
	mch := model.Channel{
		Name:  chName,
		URL:   chURL,
		Proxy: chProxy,
	}
	err := service.SaveChannel(mch)
	if err != nil {
		log.Println(err.Error())
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"ErrMsg": err.Error(),
		})
		return
	}
	c.Redirect(http.StatusFound, "/")
}

func DeleteChannelHandler(c *gin.Context) {
	chID := util.String2Uint(c.Query("id"))
	if chID == 0 {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"ErrMsg": "empty id",
		})
		return
	}
	err := service.DeleteChannel(chID)
	if err != nil {
		log.Println(err.Error())
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"ErrMsg": err.Error(),
		})
		return
	}
	c.Redirect(http.StatusFound, "/")
}
