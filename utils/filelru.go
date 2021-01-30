package utils

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

type FileLRU struct {
	maxCacheSize     int64
	currentCacheSize int64
	cacheDir         string
	lastReadTime     sync.Map
	fileSize         sync.Map
	cleanMutex       sync.Mutex
	isCleanRunning   bool
}

func NewFileLRU(cacheDir string, maxCacheSize int64) (*FileLRU, error) {
	err := os.RemoveAll(cacheDir)
	if err != nil {
		return nil, err
	}
	err = os.MkdirAll(cacheDir, 0644)
	if err != nil {
		return nil, err
	}
	return &FileLRU{cacheDir: cacheDir, maxCacheSize: maxCacheSize}, nil
}

func (lru *FileLRU) NewFile(key string, data []byte) error {
	realFilepath := filepath.Join(lru.cacheDir, key)
	dataSize := int64(len(data))
	err := ioutil.WriteFile(realFilepath, data, 0644)
	if err != nil {
		return err
	}
	lru.fileSize.Store(key, dataSize)
	lru.lastReadTime.Store(key, time.Now().Unix())
	lru.currentCacheSize += dataSize
	if lru.currentCacheSize > lru.maxCacheSize {
		go lru.Clean()
	}
	return nil
}

func (lru *FileLRU) GetFile(key string) ([]byte, error) {
	realFilepath := filepath.Join(lru.cacheDir, key)
	lru.lastReadTime.Store(key, time.Now().Unix())
	return ioutil.ReadFile(realFilepath)
}

type lruEntity struct {
	Key      string
	ReadTime int64
}

type lruEntityArray []lruEntity

func (e lruEntityArray) Len() int           { return len(e) }
func (e lruEntityArray) Less(i, j int) bool { return e[i].ReadTime < e[j].ReadTime }
func (e lruEntityArray) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }

func (lru *FileLRU) Clean() {
	if lru.tryCleanLock() {
		defer lru.finishCleanLock()
		entity := make(lruEntityArray, 0)
		lru.lastReadTime.Range(func(key interface{}, value interface{}) bool {
			cacheKey := key.(string)
			readTime := value.(int64)
			entity = append(entity, lruEntity{Key: cacheKey, ReadTime: readTime})
			return true
		})
		sort.Sort(entity)
		for _, v := range entity {
			realFilepath := filepath.Join(lru.cacheDir, v.Key)
			os.Remove(realFilepath)
			var deletedSize int64
			if dSize, ok := lru.fileSize.Load(v.Key); ok {
				deletedSize = dSize.(int64)
			}
			lru.lastReadTime.Delete(v.Key)
			lru.fileSize.Delete(v.Key)
			lru.currentCacheSize -= deletedSize
			if lru.currentCacheSize < lru.maxCacheSize {
				break
			}
		}
	}
}

func (lru *FileLRU) tryCleanLock() bool {
	lru.cleanMutex.Lock()
	defer lru.cleanMutex.Unlock()
	if lru.isCleanRunning {
		return false
	}
	lru.isCleanRunning = true
	return true
}

func (lru *FileLRU) finishCleanLock() bool {
	lru.cleanMutex.Lock()
	defer lru.cleanMutex.Unlock()
	if !lru.isCleanRunning {
		return false
	}
	lru.isCleanRunning = false
	return true
}
