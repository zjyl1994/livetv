package parser

import (
	"bufio"
	"bytes"
	"errors"
	"strings"

	"github.com/zjyl1994/livetv/utils"
)

/*
Support URL:
RTHK31 https://www.rthk.hk/feeds/dtt/rthktv31_https.m3u8
RTHK32 https://www.rthk.hk/feeds/dtt/rthktv32_https.m3u8
*/

type RTHKLiveParser struct {
}

func NewRTHKLiveParser() RTHKLiveParser {
	return RTHKLiveParser{}
}

func (RTHKLiveParser) GetLiveStream(url string) (info LiveInfo, err error) {
	content, _, err := utils.DownloadFile(url)
	if err != nil {
		return LiveInfo{}, err
	}
	var meta map[string]string
	scanner := bufio.NewScanner(bytes.NewReader(content))
	firstLine := true
	result := make([]LiveFormat, 0)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if firstLine && line != "#EXTM3U" {
			return LiveInfo{}, errors.New("not support")
		}
		firstLine = false
		if strings.HasPrefix(line, "#EXT-X-STREAM-INF:") {
			meta = make(map[string]string)
			for _, v := range utils.SplitAtCommas(strings.TrimPrefix(line, "#EXT-X-STREAM-INF:")) {
				kvPair := strings.Split(v, "=")
				if len(kvPair) == 2 {
					meta[kvPair[0]] = kvPair[1]
				}
			}
			for _, v := range []string{"BANDWIDTH", "RESOLUTION"} {
				if _, ok := meta[v]; !ok {
					return LiveInfo{}, errors.New("not support")
				}
			}
		}
		if (strings.HasPrefix(line, "https://") || strings.HasPrefix(line, "http://")) && meta != nil {
			result = append(result, LiveFormat{
				ID:     meta["BANDWIDTH"],
				URL:    line,
				Format: meta["RESOLUTION"],
			})
			meta = nil
		}
	}

	return LiveInfo{
		Title:   "Live",
		Formats: liveFormatRemoveDuplicateByFormateID(result),
	}, nil
}

func (RTHKLiveParser) ParserUpgradeable() bool {
	return false
}
func (RTHKLiveParser) ParserNeedUpdate() (bool, error) {
	return false, nil
}
func (RTHKLiveParser) ParserUpdate() error {
	return nil
}

func liveFormatRemoveDuplicateByFormateID(slc []LiveFormat) []LiveFormat {
	result := []LiveFormat{}
	for i := range slc {
		flag := true
		for j := range result {
			if slc[i].ID == result[j].ID {
				flag = false
				break
			}
		}
		if flag {
			result = append(result, slc[i])
		}
	}
	return result
}
