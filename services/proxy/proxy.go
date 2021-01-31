package proxy

import (
	"bufio"
	"bytes"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"
	uuid "github.com/satori/go.uuid"
	"github.com/zjyl1994/livetv/database"
	"github.com/zjyl1994/livetv/services/parser"
	"github.com/zjyl1994/livetv/utils"
)

var (
	tsFileCache *utils.FileLRU
	tsUrlCache  *cache.Cache
	m3u8Cache   *cache.Cache
	ErrNotFound = errors.New("proxy:not found")
)

func Init() error {
	tsUrlCache = cache.New(30*time.Second, 5*time.Minute)
	m3u8Cache = cache.New(30*time.Second, 5*time.Minute)
	var err error
	tsFileCache, err = utils.NewFileLRU(utils.DataDir("./cache"), 100*utils.MegaByte)
	if err != nil {
		return err
	}
	return nil
}

func TSProxy(tsClipUUID string) (contentType string, data []byte, err error) {
	if tsFileCache.CheckFileExist(tsClipUUID) {
		data, meta, err := tsFileCache.GetFile(tsClipUUID)
		return meta.(string), data, err
	}
	if realUrl, ok := tsUrlCache.Get(tsClipUUID); !ok {
		return "", nil, ErrNotFound
	} else {
		data, contentType, err = utils.DownloadFile(realUrl.(string))
		if err != nil {
			return "", nil, err
		}
		err = tsFileCache.NewFile(tsClipUUID, data, contentType)
		if err != nil {
			return "", nil, err
		}
		return
	}
}

func M3U8Proxy(channelID string) (contentType string, data []byte, err error) {
	cacheContent, ok := m3u8Cache.Get(channelID)
	if ok {
		return "application/x-mpegURL", cacheContent.([]byte), nil
	} else {
		iChannelID, err := strconv.Atoi(channelID)
		if err != nil {
			return "", nil, err
		}
		chInfo, err := database.GetChannel(uint(iChannelID))
		if err != nil {
			return "", nil, err
		}
		parser, err := parser.PickParser(chInfo.URL)
		if err != nil {
			return "", nil, err
		}
		liveInfo, err := parser.GetLiveStream(chInfo.URL)
		if err != nil {
			return "", nil, err
		}
		maxFormatID := 0
		maxURL := ""
		for _, v := range liveInfo.Formats {
			id, err := strconv.Atoi(v.ID)
			if err != nil {
				continue
			}
			if id > maxFormatID {
				maxURL = v.URL
			}
		}
		if maxURL == "" {
			return "", nil, ErrNotFound
		}
		baseUrl, err := database.GetConfig("base_url")
		if err != nil {
			return "", nil, err
		}
		data, contentType, err := utils.DownloadFile(maxURL)
		if err != nil {
			return "", nil, err
		}
		if chInfo.Proxy {
			newData, err := realM3U8Trans(data, tsUrlCache, baseUrl+"/ts?id=")
			if err != nil {
				return "", nil, err
			}
			return contentType, newData, nil
		} else {
			return contentType, data, nil
		}
	}
}

func realM3U8Trans(data []byte, urlCache *cache.Cache, prefixURL string) ([]byte, error) {
	var sb strings.Builder
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		l := scanner.Text()
		if strings.HasPrefix(l, "#") || !isTSLink(l) {
			sb.WriteString(l)
		} else {
			sb.WriteString(prefixURL)
			tsClipUUID := uuid.NewV4().String()
			urlCache.Set(tsClipUUID, l, cache.DefaultExpiration)
			sb.WriteString(tsClipUUID)
		}
		sb.WriteString("\n")
	}
	return []byte(sb.String()), nil
}

func isTSLink(sUrl string) bool {
	ext, err := utils.GetUrlExt(sUrl)
	if err != nil {
		return false
	}
	return strings.EqualFold(".ts", ext)
}
