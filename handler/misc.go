package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zjyl1994/livetv/service"
)

func IndexHandler(c *gin.Context) {
	channels, err := service.GetAllChannel()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"ErrMsg": err.Error(),
		})
		return
	}
	c.HTML(http.StatusOK, "index.html", gin.H{
		"Channels": channels,
	})
}
