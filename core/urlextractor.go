package core

import (
	"fmt"
	"os"

	"github.com/go-logr/logr"
	"mvdan.cc/xurls/v2"
)

type URLExtractor struct {
	logger        logr.Logger
	filePath      string
	cache         *Cache
	titleFetchers []TitleFetcher
	noCache       bool
}

func NewURLExtractor(logger logr.Logger, filePath string, titleFetchers []TitleFetcher, noCache bool) (*URLExtractor, error) {
	var cache *Cache
	var err error
	if !noCache {
		cache, err = NewCache(logger)
		if err != nil {
			return nil, fmt.Errorf("failed to create cache: %w", err)
		}
	}

	return &URLExtractor{
		logger:        logger,
		filePath:      filePath,
		cache:         cache,
		titleFetchers: titleFetchers,
		noCache:       noCache,
	}, nil
}

func (ue *URLExtractor) ExtractURLs() ([]string, error) {
	ue.logger.V(1).Info("Debug: Extracting URLs from file", "path", ue.filePath)
	content, err := os.ReadFile(ue.filePath)
	if err != nil {
		ue.logger.Error(err, "Failed to read file", "path", ue.filePath)
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	rx := xurls.Strict()
	urls := rx.FindAllString(string(content), -1)

	ue.logger.V(1).Info("Debug: URLs extracted", "count", len(urls))
	return urls, nil
}

func (ue *URLExtractor) GetOrFetchTitle(url string) (string, error) {
	if !ue.noCache {
		if title, ok := ue.cache.Get(url); ok {
			ue.logger.V(1).Info("Debug: Title found in cache", "url", url)
			return title, nil
		}
	}

	ue.logger.V(1).Info("Debug: Fetching title from web", "url", url)
	var lastErr error
	for _, fetcher := range ue.titleFetchers {
		title, err := fetcher.FetchTitle(url)
		if err == nil {
			if !ue.noCache {
				if err := ue.cache.Set(url, title); err != nil {
					ue.logger.Error(err, "Failed to cache title", "url", url)
				}
			}
			return title, nil
		}
		lastErr = err
		ue.logger.V(2).Info("Debug: Fetcher failed, trying next", "url", url, "error", err.Error())
	}

	return "", fmt.Errorf("all fetchers failed to fetch title: %w", lastErr)
}
