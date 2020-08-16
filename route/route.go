package route

import (
	"github.com/gin-gonic/gin"
	"github.com/zjyl1994/livetv/handler"
)

func Register(r *gin.Engine) {
	r.LoadHTMLGlob("view/*")
	r.GET("/lives.m3u", handler.M3UHandler)
	r.GET("/live.m3u8", handler.LiveHandler)
	r.GET("/p/live.m3u8", handler.M3U8ProxyHandler)
	r.GET("/p/live.ts", handler.TsProxyHandler)
}
