package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"github.com/robfig/cron/v3"
	"github.com/zjyl1994/livetv/global"
	"github.com/zjyl1994/livetv/route"
	"github.com/zjyl1994/livetv/service"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Server listen", os.Getenv("LIVETV_LISTEN"))
	log.Println("Server datadir", os.Getenv("LIVETV_DATADIR"))
	logFile, err := os.OpenFile(os.Getenv("LIVETV_DATADIR")+"/livetv.log", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Panicln(err)
	}
	log.SetOutput(io.MultiWriter(os.Stderr, logFile))
	err = global.InitDB(os.Getenv("LIVETV_DATADIR") + "/livetv.db")
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
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Static("/assert", "./assert")
	router.GET("/version", func(c *gin.Context) {
		c.String(http.StatusOK, global.VersionString)
	})
	route.Register(router)
	srv := &http.Server{
		Addr:    os.Getenv("LIVETV_LISTEN"),
		Handler: router,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Panicf("listen: %s\n", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shuting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Panicf("Server forced to shutdown: %s\n", err)
	}
	log.Println("Server exiting")
	logFile.Close()
}
