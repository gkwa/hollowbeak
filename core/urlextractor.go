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
	return &URLExtractor{
		logger:     logger,
		filePath:   filePath,
		cache:      make(map[string]string),
		cacheFile:  cacheFile,
		strictMode: strictMode,
	}
}

func (ue *URLExtractor) ExtractURLs() ([]string, error) {
	content, err := os.ReadFile(ue.filePath)
	if err != nil {
		ue.logger.Error(err, "Failed to read file", "path", ue.filePath)
		return nil, err
	}

	rxStrict := xurls.Strict()
	rxRelaxed := xurls.Relaxed()

	var urls []string
	if ue.strictMode {
		urls = rxStrict.FindAllString(string(content), -1)
	} else {
		urls = rxRelaxed.FindAllString(string(content), -1)
	}

	return urls, nil
}

func (ue *URLExtractor) LoadCache() error {
	data, err := os.ReadFile(ue.cacheFile)
	if err != nil {
		if os.IsNotExist(err) {
			ue.logger.Info("Cache file not found, starting with empty cache", "path", ue.cacheFile)
			return nil
		}
		return err
	}

	return json.Unmarshal(data, &ue.cache)
}

func (ue *URLExtractor) SaveCache() error {
	data, err := json.MarshalIndent(ue.cache, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(ue.cacheFile, data, 0o644)
}

func (ue *URLExtractor) GetTitle(url string) (string, bool) {
	title, ok := ue.cache[url]
	return title, ok
}

func (ue *URLExtractor) SetTitle(url, title string) {
	ue.cache[url] = title
}
