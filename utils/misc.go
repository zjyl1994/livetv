package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

const (
	MegaByte = 1024 * 1024
)

func DownloadFile(url string) ([]byte, string, error) {
	c := &http.Client{Timeout: 15 * time.Second}
	resp, err := c.Get(url)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("bad status: %s", resp.Status)
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}
	return content, resp.Header.Get("Content-Type"), nil
}

func SplitAtCommas(s string) []string {
	res := []string{}
	var beg int
	var inString bool

	for i := 0; i < len(s); i++ {
		if s[i] == ',' && !inString {
			res = append(res, s[beg:i])
			beg = i + 1
		} else if s[i] == '"' {
			if !inString {
				inString = true
			} else if i > 0 && s[i-1] != '\\' {
				inString = false
			}
		}
	}
	return append(res, s[beg:])
}

func DataDir(path string) string {
	datadir := os.Getenv("LIVETV_DATADIR")
	if datadir == "" {
		datadir = "./data"
	}
	if path == "" {
		return datadir
	} else {
		return filepath.Join(datadir, path)
	}
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func GetUrlExt(sUrl string) (string, error) {
	u, err := url.Parse(sUrl)
	if err != nil {
		return "", err
	}
	return filepath.Ext(u.Path), nil
}
