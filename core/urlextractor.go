package core

import (
	"fmt"
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
		return nil, fmt.Errorf("failed to create cache: %w", err)
	}

	return &URLExtractor{
		logger:   logger,
		filePath: filePath,
		cache:    cache,
	}, nil
}

func (urlExtractor *URLExtractor) ExtractURLs() ([]string, error) {
	urlExtractor.logger.V(1).Info("Debug: Extracting URLs from file", "path", urlExtractor.filePath)
	content, err := os.ReadFile(urlExtractor.filePath)
	if err != nil {
		urlExtractor.logger.Error(err, "Failed to read file", "path", urlExtractor.filePath)
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	rx := xurls.Relaxed()
	urls := rx.FindAllString(string(content), -1)

	urlExtractor.logger.V(1).Info("Debug: URLs extracted", "count", len(urls))
	return urls, nil
}

func (urlExtractor *URLExtractor) GetOrFetchTitle(url string) (string, error) {
	if title, ok := urlExtractor.cache.Get(url); ok {
		urlExtractor.logger.V(1).Info("Debug: Title found in cache", "url", url)
		return title, nil
	}

	urlExtractor.logger.V(1).Info("Debug: Fetching title from web", "url", url)
	title, err := getPageTitle(urlExtractor.logger, url)
	if err != nil {
		return "", fmt.Errorf("failed to get page title: %w", err)
	}

	if err := urlExtractor.cache.Set(url, title); err != nil {
		urlExtractor.logger.Error(err, "Failed to cache title", "url", url)
	}
	return title, nil
}
