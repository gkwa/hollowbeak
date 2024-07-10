package core

import (
	"path/filepath"

	"github.com/go-logr/logr"
)

func Hello(logger logr.Logger) {
	logger.V(1).Info("Debug: Entering Hello function")
	logger.Info("Hello, World!")

	filePath := "/Users/mtm/Documents/Obsidian Vault/2024-07-10.md"
	cacheFile := filepath.Join(filepath.Dir(filePath), "data.json")

	extractor := NewURLExtractor(logger, filePath, cacheFile, false)
	err := extractor.LoadCache()
	if err != nil {
		logger.Error(err, "Failed to load cache")
		return
	}

	urls, err := extractor.ExtractURLs()
	if err != nil {
		logger.Error(err, "Failed to extract URLs")
		return
	}

	for _, url := range urls {
		logger.V(1).Info("Debug: Processing URL", "url", url)

		if title, ok := extractor.GetTitle(url); ok {
			logger.Info("Found cached title", "url", url, "title", title)
			continue
		}

		title, err := getPageTitle(logger, url)
		if err != nil {
			logger.Error(err, "Failed to get page title", "url", url)
			continue
		}

		extractor.SetTitle(url, title)
		logger.Info("Fetched new title", "url", url, "title", title)
	}

	err = extractor.SaveCache()
	if err != nil {
		logger.Error(err, "Failed to save cache")
	}

	logger.V(1).Info("Debug: Exiting Hello function")
}
