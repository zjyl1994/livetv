package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	err := initProc()
	if err != nil {
		log.Fatalf("init: %s\n", err)
	}
	defer removePidFile()
	log.Println("LiveTV starting...")
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, versionString)
	})
	router.GET("/lives.m3u", m3uHandler)
	router.GET("/live.m3u8", liveHandler)
	router.GET("/p/live.m3u8", m3u8ProxyHandler)
	router.GET("/p/live.ts", tsProxyHandler)
	srv := &http.Server{
		Addr:    cfg.ListenOn,
		Handler: router,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shuting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Println("Server exiting")
}
