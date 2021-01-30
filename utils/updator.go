package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/hashicorp/go-version"
)

type GithubUpdator struct {
	repoName                string
	localPath               string
	wantFile                string
	releaseJsonCache        githubReleaseJSON
	releaseJsonLastReadTime time.Time
}

type githubReleaseJSON struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name string `json:"name"`
		URL  string `json:"browser_download_url"`
	} `json:"assets"`
}

const (
	githubReleaseAPIUrl = "https://api.github.com/repos/%s/releases/latest"
)

func NewGithubUpdator(repo, localPath string, wantFile string) GithubUpdator {
	return GithubUpdator{
		repoName:  repo,
		localPath: localPath,
		wantFile:  wantFile,
	}
}

func (gu GithubUpdator) CheckUpdate() (bool, error) {
	currentVersion, err := gu.getVersionLock()
	if err != nil {
		return false, err
	}
	if currentVersion == "" {
		return true, nil
	}

	rjson, err := gu.loadReleaseJSON()
	if err != nil {
		return false, err
	}
	releaseVersion := rjson.TagName
	if releaseVersion == "" {
		return false, errors.New("read github release failed")
	}
	cv, err := version.NewVersion(currentVersion)
	if err != nil {
		return false, err
	}
	rv, err := version.NewVersion(releaseVersion)
	if err != nil {
		return false, err
	}
	needUpdate := rv.GreaterThan(cv)
	return needUpdate, nil
}

func (gu GithubUpdator) Update() error {
	rjson, err := gu.loadReleaseJSON()
	if err != nil {
		return err
	}
	releaseVersion := rjson.TagName
	for _, v := range rjson.Assets {
		if v.Name == gu.wantFile {
			out, err := os.Create(filepath.Join(gu.localPath, gu.wantFile))
			if err != nil {
				return err
			}
			defer out.Close()
			resp, err := http.Get(v.URL)
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			_, err = io.Copy(out, resp.Body)
			if err != nil {
				return err
			}
			return gu.setVersionLock(releaseVersion)
		}
	}
	return nil
}

func (gu GithubUpdator) getVersionLock() (string, error) {
	bVersion, err := ioutil.ReadFile(filepath.Join(gu.localPath, "version.lock"))
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		} else {
			return "", err
		}
	} else {
		return string(bVersion), nil
	}
}

func (gu GithubUpdator) setVersionLock(version string) error {
	return ioutil.WriteFile(filepath.Join(gu.localPath, "version.lock"), []byte(version), 0644)
}

func (gu GithubUpdator) loadReleaseJSON() (githubReleaseJSON, error) {
	if time.Since(gu.releaseJsonLastReadTime).Minutes() < 5 {
		return gu.releaseJsonCache, nil
	} else {
		resp, err := http.Get(fmt.Sprintf(githubReleaseAPIUrl, gu.repoName))
		if err != nil {
			return githubReleaseJSON{}, err
		}
		defer resp.Body.Close()
		bJson, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return githubReleaseJSON{}, err
		}
		var releases githubReleaseJSON
		err = json.Unmarshal(bJson, &releases)
		if err != nil {
			return githubReleaseJSON{}, err
		}
		gu.releaseJsonLastReadTime = time.Now()
		gu.releaseJsonCache = releases
		return releases, nil
	}
}
