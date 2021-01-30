package parser

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/zjyl1994/livetv/utils"
)

var (
	ErrParserError = errors.New("parser error")
)

type YoutubeLiveParser struct {
	datadir string
	updator utils.GithubUpdator
}

type ytdlOutput struct {
	Title   string `json:"title"`
	Formats []struct {
		ID     string `json:"format_id"`
		Format string `json:"format"`
		Url    string `json:"url"`
	} `json:"formats"`
}

func (o ytdlOutput) ToLiveInfo() LiveInfo {
	fmts := make([]LiveFormat, len(o.Formats))
	for i, v := range o.Formats {
		fmts[i] = LiveFormat{
			ID:     v.ID,
			Format: v.Format,
			URL:    v.Url,
		}
	}
	return LiveInfo{
		Title:   o.Title,
		Formats: fmts,
	}
}

func NewYoutubeLiveParser(datadir string) *YoutubeLiveParser {
	return &YoutubeLiveParser{
		datadir: datadir,
		updator: utils.NewGithubUpdator("ytdl-org/youtube-dl", datadir, "youtube-dl"),
	}
}

func (p *YoutubeLiveParser) GetLiveStream(url string) (info LiveInfo, err error) {
	jsonOutput, err := p.callYoutubedl(url)
	if err != nil {
		return LiveInfo{}, err
	}
	var output ytdlOutput
	if !json.Valid(jsonOutput) {
		return LiveInfo{}, fmt.Errorf("%w:%s", ErrParserError, string(jsonOutput))
	}
	err = json.Unmarshal(jsonOutput, &output)
	if err != nil {
		return LiveInfo{}, err
	}
	return output.ToLiveInfo(), nil
}

func (p *YoutubeLiveParser) callYoutubedl(url string) (output []byte, err error) {
	execPath := filepath.Join(p.datadir, "youtube-dl")
	_, err = exec.LookPath(execPath)
	if err != nil {
		return nil, err
	} else {
		ctx, cancelFunc := context.WithTimeout(context.Background(), 8*time.Second)
		defer cancelFunc()
		cmd := exec.CommandContext(ctx, "python3", execPath, "-j", url)
		return cmd.CombinedOutput()
	}
}

// Placeholder
func (p *YoutubeLiveParser) ParserUpgradeable() bool {
	return true
}

func (p *YoutubeLiveParser) ParserNeedUpdate() (bool, error) {
	return p.updator.CheckUpdate()
}
func (p *YoutubeLiveParser) ParserUpdate() error {
	return p.updator.Update()
}
