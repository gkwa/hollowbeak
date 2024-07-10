package core

import (
	"encoding/json"
	"os"

	"github.com/go-logr/logr"
	"mvdan.cc/xurls/v2"
)

type URLExtractor struct {
	logger     logr.Logger
	filePath   string
	cache      map[string]string
	cacheFile  string
	strictMode bool
}

func NewURLExtractor(logger logr.Logger, filePath string, cacheFile string, strictMode bool) *URLExtractor {
	logger.V(1).Info("Debug: Creating new URLExtractor", "filePath", filePath, "cacheFile", cacheFile, "strictMode", strictMode)
	return &URLExtractor{
		logger:     logger,
		filePath:   filePath,
		cache:      make(map[string]string),
		cacheFile:  cacheFile,
		strictMode: strictMode,
	}
}

func (ue *URLExtractor) ExtractURLs() ([]string, error) {
	ue.logger.V(1).Info("Debug: Extracting URLs from file", "path", ue.filePath)
	content, err := os.ReadFile(ue.filePath)
	if err != nil {
		ue.logger.Error(err, "Failed to read file", "path", ue.filePath)
		return nil, err
	}
	ue.logger.V(2).Info("Debug: File read successfully", "bytes", len(content))

	ue.logger.V(2).Info("Debug: Initializing URL finder", "strictMode", ue.strictMode)
	rxStrict := xurls.Strict()
	rxRelaxed := xurls.Relaxed()

	var urls []string
	if ue.strictMode {
		ue.logger.V(2).Info("Debug: Using strict URL matching")
		urls = rxStrict.FindAllString(string(content), -1)
	} else {
		ue.logger.V(2).Info("Debug: Using relaxed URL matching")
		urls = rxRelaxed.FindAllString(string(content), -1)
	}
	ue.logger.V(1).Info("Debug: URLs extracted", "count", len(urls))

	return urls, nil
}

func (ue *URLExtractor) LoadCache() error {
	ue.logger.V(1).Info("Debug: Loading cache", "path", ue.cacheFile)
	data, err := os.ReadFile(ue.cacheFile)
	if err != nil {
		if os.IsNotExist(err) {
			ue.logger.Info("Cache file not found, starting with empty cache", "path", ue.cacheFile)
			return nil
		}
		ue.logger.Error(err, "Failed to read cache file")
		return err
	}
	ue.logger.V(2).Info("Debug: Cache file read successfully", "bytes", len(data))

	err = json.Unmarshal(data, &ue.cache)
	if err != nil {
		ue.logger.Error(err, "Failed to unmarshal cache data")
		return err
	}
	ue.logger.V(1).Info("Debug: Cache loaded successfully", "entries", len(ue.cache))

	return nil
}

func (ue *URLExtractor) SaveCache() error {
	ue.logger.V(1).Info("Debug: Saving cache", "path", ue.cacheFile, "entries", len(ue.cache))
	data, err := json.MarshalIndent(ue.cache, "", "  ")
	if err != nil {
		ue.logger.Error(err, "Failed to marshal cache data")
		return err
	}
	ue.logger.V(2).Info("Debug: Cache data marshaled successfully", "bytes", len(data))

	err = os.WriteFile(ue.cacheFile, data, 0o644)
	if err != nil {
		ue.logger.Error(err, "Failed to write cache file")
		return err
	}
	ue.logger.V(1).Info("Debug: Cache saved successfully")

	return nil
}

func (ue *URLExtractor) GetTitle(url string) (string, bool) {
	ue.logger.V(2).Info("Debug: Getting title from cache", "url", url)
	title, ok := ue.cache[url]
	if ok {
		ue.logger.V(2).Info("Debug: Title found in cache", "url", url, "title", title)
	} else {
		ue.logger.V(2).Info("Debug: Title not found in cache", "url", url)
	}
	return title, ok
}

func (ue *URLExtractor) SetTitle(url, title string) {
	ue.logger.V(2).Info("Debug: Setting title in cache", "url", url, "title", title)
	ue.cache[url] = title
}
