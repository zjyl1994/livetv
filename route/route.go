package route

import (
	"github.com/gin-gonic/gin"
	"github.com/zjyl1994/livetv/handler"
)

func Register(r *gin.Engine) {
	r.LoadHTMLGlob("view/*")

	r.GET("/lives.m3u", handler.M3UHandler)
	r.GET("/live.m3u8", handler.LiveHandler)
	r.GET("/live.ts", handler.TsProxyHandler)
	r.GET("/cache.txt", handler.CacheHandler)

	r.GET("/", handler.IndexHandler)
	r.POST("/api/newchannel", handler.NewChannelHandler)
	r.GET("/api/delchannel", handler.DeleteChannelHandler)
	r.POST("/api/updconfig", handler.UpdateConfigHandler)
	r.GET("/log", handler.LogHandler)
	r.GET("/login", handler.LoginViewHandler)
	r.POST("/api/login", handler.LoginActionHandler)
	r.GET("/api/logout", handler.LogoutHandler)
	r.POST("/api/changepwd", handler.ChangePasswordHandler)
}
