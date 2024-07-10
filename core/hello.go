package core

import (
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/go-logr/logr"
)

func Hello(logger logr.Logger) {
	logger.V(1).Info("Debug: Entering Hello function")
	logger.Info("Hello, World!")

	logger.V(1).Info("Debug: Setting up file paths")
	filePath := "/Users/mtm/Documents/Obsidian Vault/2024-07-10.md"
	cacheFile, err := getCachePath("hollowbeak/data.json")
	if err != nil {
		logger.Error(err, "Failed to get cache path")
		return
	}
	logger.V(2).Info("Debug: File paths set", "filePath", filePath, "cacheFile", cacheFile)

	logger.V(1).Info("Debug: Creating new URLExtractor")
	extractor := NewURLExtractor(logger, filePath, cacheFile, false)
	logger.V(2).Info("Debug: URLExtractor created", "strictMode", false)

	logger.V(1).Info("Debug: Loading cache")
	err = extractor.LoadCache()
	if err != nil {
		logger.Error(err, "Failed to load cache")
		return
	}
	logger.V(2).Info("Debug: Cache loaded successfully")

	logger.V(1).Info("Debug: Extracting URLs from file")
	urls, err := extractor.ExtractURLs()
	if err != nil {
		logger.Error(err, "Failed to extract URLs")
		return
	}
	logger.V(2).Info("Debug: URLs extracted", "count", len(urls))

	for _, url := range urls {
		logger.V(1).Info("Debug: Processing URL", "url", url)

		logger.V(2).Info("Debug: Checking cache for URL")
		if title, ok := extractor.GetTitle(url); ok {
			logger.Info("Found cached title", "url", url, "title", title)
			logger.V(2).Info("Debug: Title found in cache", "url", url, "title", title)
			continue
		}
		logger.V(2).Info("Debug: Title not found in cache, fetching from web", "url", url)

		logger.V(2).Info("Debug: Fetching page title", "url", url)
		title, err := getPageTitle(logger, url)
		if err != nil {
			logger.Error(err, "Failed to get page title", "url", url)
			continue
		}
		logger.V(2).Info("Debug: Page title fetched successfully", "url", url, "title", title)

		logger.V(2).Info("Debug: Updating cache with new title", "url", url, "title", title)
		extractor.SetTitle(url, title)
		logger.Info("Fetched new title", "url", url, "title", title)
	}

	logger.V(1).Info("Debug: Saving updated cache")
	err = extractor.SaveCache()
	if err != nil {
		logger.Error(err, "Failed to save cache")
	}
	logger.V(2).Info("Debug: Cache saved successfully")

	logger.V(1).Info("Debug: Exiting Hello function")
}

func getCachePath(configRelPath string) (string, error) {
	configFilePath, err := xdg.ConfigFile(configRelPath)
	if err != nil {
		return "", err
	}

	dirPerm := os.FileMode(0o700)

	d := filepath.Dir(configFilePath)

	if err := os.MkdirAll(d, dirPerm); err != nil {
		return "", err
	}

	return configFilePath, nil
}
