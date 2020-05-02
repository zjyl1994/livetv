package main

import (
	"context"
	"log"
	"os/exec"
	"strings"
	"time"
)

func getYoutubeLiveM3U8(youtubeURL string) (string, error) {
	liveURL, ok := urlCache.Load(youtubeURL)
	if ok {
		log.Println("cache hit", youtubeURL)
		return liveURL.(string), nil
	} else {
		log.Println("cache miss", youtubeURL)
		liveURL, err := realGetYoutubeLiveM3U8(youtubeURL)
		if err != nil {
			return "", err
		} else {
			urlCache.Store(youtubeURL, liveURL)
			return liveURL, nil
		}
	}
}

func realGetYoutubeLiveM3U8(youtubeURL string) (string, error) {
	ytdlArgs := strings.Fields(cfg.YtdlArgs)
	for i, v := range ytdlArgs {
		if strings.EqualFold(v, "{url}") {
			ytdlArgs[i] = youtubeURL
		}
	}
	_, err := exec.LookPath(cfg.YtdlCmd)
	if err != nil {
		return "", err
	} else {
		ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancelFunc()
		cmd := exec.CommandContext(ctx, cfg.YtdlCmd, ytdlArgs...)
		out, err := cmd.CombinedOutput()
		return strings.TrimSpace(string(out)), err
	}
}
