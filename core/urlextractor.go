package core

import (
	"os"

	"github.com/go-logr/logr"
	"mvdan.cc/xurls/v2"
)

type URLExtractor struct {
	logger   logr.Logger
	filePath string
	cache    *Cache
}

func NewURLExtractor(logger logr.Logger, filePath string) (*URLExtractor, error) {
	cache, err := NewCache(logger)
	if err != nil {
		return nil, err
	}

	return &URLExtractor{
		logger:   logger,
		filePath: filePath,
		cache:    cache,
	}, nil
}

func (ue *URLExtractor) ExtractURLs() ([]string, error) {
	ue.logger.V(1).Info("Debug: Extracting URLs from file", "path", ue.filePath)
	content, err := os.ReadFile(ue.filePath)
	if err != nil {
		ue.logger.Error(err, "Failed to read file", "path", ue.filePath)
		return nil, err
	}

	rx := xurls.Relaxed()
	urls := rx.FindAllString(string(content), -1)

	ue.logger.V(1).Info("Debug: URLs extracted", "count", len(urls))
	return urls, nil
}

func (ue *URLExtractor) GetOrFetchTitle(url string) (string, error) {
	if title, ok := ue.cache.Get(url); ok {
		ue.logger.V(1).Info("Debug: Title found in cache", "url", url)
		return title, nil
	}

	ue.logger.V(1).Info("Debug: Fetching title from web", "url", url)
	title, err := getPageTitle(ue.logger, url)
	if err != nil {
		return "", err
	}

	if err := ue.cache.Set(url, title); err != nil {
		ue.logger.Error(err, "Failed to cache title", "url", url)
		// We don't return here because we still want to return the title even if caching failed
	}
	return title, nil
}
