package main

import (
	"context"
	"os/exec"
	"strings"
	"time"
)

func getYoutubeLiveM3U8(youtubeURL string) (string, error) {
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
