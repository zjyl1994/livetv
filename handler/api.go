package handler

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zjyl1994/livetv/model"
	"github.com/zjyl1994/livetv/service"
	"github.com/zjyl1994/livetv/util"

	"golang.org/x/text/language"
)

var langMatcher = language.NewMatcher([]language.Tag{
	language.English,
	language.Chinese,
})

func IndexHandler(c *gin.Context) {
	acceptLang := c.Request.Header.Get("Accept-Language")
	langTag, _ := language.MatchStrings(langMatcher, acceptLang)

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
	var m3uName string
	if langTag == language.Chinese {
		m3uName = "M3U 頻道列表"
	} else {
		m3uName = "M3U File"
	}
	channels := make([]Channel, len(channelModels)+1)
	channels[0] = Channel{
		ID:   0,
		Name: m3uName,
		M3U8: baseUrl + "/lives.m3u",
	}
	for i, v := range channelModels {
		channels[i+1] = Channel{
			ID:    v.ID,
			Name:  v.Name,
			URL:   v.URL,
			M3U8:  baseUrl + "/live.m3u8?channel=" + strconv.Itoa(int(v.ID)),
			Proxy: v.Proxy,
		}
	}
	conf, err := loadConfig()
	if err != nil {
		log.Println(err.Error())
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"ErrMsg": err.Error(),
		})
		return
	}

	var templateFilename string
	if langTag == language.Chinese {
		templateFilename = "index-zh.html"
	} else {
		templateFilename = "index.html"
	}
	c.HTML(http.StatusOK, templateFilename, gin.H{
		"Channels": channels,
		"Configs":  conf,
	})
}

func loadConfig() (Config, error) {
	var conf Config
	if cmd, err := service.GetConfig("ytdl_cmd"); err != nil {
		return conf, err
	} else {
		conf.Cmd = cmd
	}
	if args, err := service.GetConfig("ytdl_args"); err != nil {
		return conf, err
	} else {
		conf.Args = args
	}
	if burl, err := service.GetConfig("base_url"); err != nil {
		return conf, err
	} else {
		conf.BaseURL = burl
	}
	return conf, nil
}

func NewChannelHandler(c *gin.Context) {
	chName := c.PostForm("name")
	chURL := c.PostForm("url")
	if chName == "" || chURL == "" {
		c.Redirect(http.StatusFound, "/")
		return
	}
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

func UpdateConfigHandler(c *gin.Context) {
	ytdlCmd := c.PostForm("cmd")
	ytdlArgs := c.PostForm("args")
	baseUrl := strings.TrimSuffix(c.PostForm("baseurl"), "/")
	if len(ytdlCmd) > 0 {
		err := service.SetConfig("ytdl_cmd", ytdlCmd)
		if err != nil {
			log.Println(err.Error())
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{
				"ErrMsg": err.Error(),
			})
			return
		}
	}
	if len(ytdlArgs) > 0 {
		err := service.SetConfig("ytdl_args", ytdlArgs)
		if err != nil {
			log.Println(err.Error())
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{
				"ErrMsg": err.Error(),
			})
			return
		}
	}
	if len(baseUrl) > 0 {
		err := service.SetConfig("base_url", baseUrl)
		if err != nil {
			log.Println(err.Error())
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{
				"ErrMsg": err.Error(),
			})
			return
		}
	}
	c.Redirect(http.StatusFound, "/")
}

func LogHandler(c *gin.Context) {
	c.File(os.Getenv("LIVETV_DATADIR") + "/livetv.log")
}
