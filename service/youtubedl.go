package service

import (
	"context"
	"log"
	"os/exec"
	"strings"

	"github.com/zjyl1994/livetv/global"
)

func GetYoutubeLiveM3U8(youtubeURL string) (string, error) {
	liveURL, ok := global.URLCache.Load(youtubeURL)
	if ok {
		return liveURL.(string), nil
	} else {
		log.Println("cache miss", youtubeURL)
		liveURL, err := RealGetYoutubeLiveM3U8(youtubeURL)
		if err != nil {
			log.Println(err)
			log.Println("[YTDL]",liveURL)
			return "", err
		} else {
			global.URLCache.Store(youtubeURL, liveURL)
			return liveURL, nil
		}
	}
}

func RealGetYoutubeLiveM3U8(youtubeURL string) (string, error) {
	YtdlCmd, err := GetConfig("ytdl_cmd")
	if err != nil {
		log.Println(err)
		return "", err
	}
	YtdlArgs, err := GetConfig("ytdl_args")
	if err != nil {
		log.Println(err)
		return "", err
	}
	ytdlArgs := strings.Fields(YtdlArgs)
	for i, v := range ytdlArgs {
		if strings.EqualFold(v, "{url}") {
			ytdlArgs[i] = youtubeURL
		}
	}
	_, err = exec.LookPath(YtdlCmd)
	if err != nil {
		log.Println(err)
		return "", err
	} else {
		ctx, cancelFunc := context.WithTimeout(context.Background(), global.HttpClientTimeout)
		defer cancelFunc()
		log.Println(YtdlCmd, ytdlArgs)
		cmd := exec.CommandContext(ctx, YtdlCmd, ytdlArgs...)
		out, err := cmd.CombinedOutput()
		return strings.TrimSpace(string(out)), err
	}
}
