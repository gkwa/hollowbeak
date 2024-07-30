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

func (ue *URLExtractor) ExtractURLs() ([]urlRecord, error) {
	ue.logger.V(1).Info("Debug: Extracting URLs from file", "path", ue.filePath)
	content, err := os.ReadFile(ue.filePath)
	if err != nil {
		ue.logger.Error(err, "Failed to read file", "path", ue.filePath)
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	rx := xurls.Strict()
	rawURLs := rx.FindAllString(string(content), -1)

	urls := make([]urlRecord, len(rawURLs))
	for i, rawURL := range rawURLs {
		urls[i] = newURLRecord(rawURL)
	}

	ue.logger.V(1).Info("Debug: URLs extracted", "count", len(urls))
	return urls, nil
}

func (ue *URLExtractor) GetOrFetchTitles(urls []urlRecord) (map[string]string, error) {
	titles := make(map[string]string)
	urlsToFetch := make([]urlRecord, 0)

	if ue.noCache {
		urlsToFetch = urls
	} else {
		for _, url := range urls {
			if title, ok := ue.cache.Get(url.URL); ok {
				ue.logger.V(1).Info("Debug: Title found in cache", "url", url.URL, "title", title)
				titles[url.URL] = title
			} else {
				urlsToFetch = append(urlsToFetch, url)
			}
		}
	}

	if len(urlsToFetch) > 0 {
		ue.logger.V(1).Info("Debug: Fetching titles from web", "urlCount", len(urlsToFetch))
		var lastErr error
		for _, fetcher := range ue.titleFetchers {
			fetchedTitles, err := fetcher.FetchTitles(urlsToFetch)
			if err == nil {
				for url, title := range fetchedTitles {
					titles[url] = title
					if !ue.noCache {
						if err := ue.cache.Set(url, title); err != nil {
							ue.logger.Error(err, "Failed to cache title", "url", url)
						}
					}
				}
				return titles, nil
			}
			lastErr = err
			ue.logger.V(2).Info("Debug: Fetcher failed, trying next", "error", err.Error())
		}
		return titles, fmt.Errorf("all fetchers failed to fetch titles: %w", lastErr)
	}

	return titles, nil
}
