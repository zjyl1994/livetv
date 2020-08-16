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
	"github.com/robfig/cron/v3"
	"github.com/zjyl1994/livetv/global"
	"github.com/zjyl1994/livetv/route"
	"github.com/zjyl1994/livetv/service"
)

func main() {
	err := global.InitDB("./data/livetv.db")
	if err != nil {
		log.Panicf("init: %s\n", err)
	}
	log.Println("LiveTV starting...")
	go service.LoadChannelCache()
	c := cron.New()
	_, err = c.AddFunc("0 */4 * * *", service.UpdateURLCache)
	if err != nil {
		log.Panicf("preloadCron: %s\n", err)
	}
	c.Start()
	gin.SetMode(gin.DebugMode)
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, global.VersionString)
	})
	route.Register(router)
	srv := &http.Server{
		Addr:    ":9000",
		Handler: router,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	quit := make(chan os.Signal, 1)
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
