package core

import (
	"fmt"
	"strings"

	"github.com/go-logr/logr"
)

func Hello(logger logr.Logger) {
	logger.V(1).Info("Debug: Entering Hello function")

	filePath := "/Users/mtm/Documents/Obsidian Vault/2024-07-10.md"

	logger.V(1).Info("Debug: Creating new URLExtractor")
	extractor, err := NewURLExtractor(logger, filePath)
	if err != nil {
		logger.Error(err, "Failed to create URLExtractor")
		return
	}
	logger.V(2).Info("Debug: URLExtractor created")

	logger.V(1).Info("Debug: Extracting URLs from file")
	urls, err := extractor.ExtractURLs()
	if err != nil {
		logger.Error(err, "Failed to extract URLs")
		return
	}
	logger.V(2).Info("Debug: URLs extracted", "count", len(urls))

	var stringBuffer strings.Builder

	for _, url := range urls {
		logger.V(1).Info("Debug: Processing URL", "url", url)

		title, err := extractor.GetOrFetchTitle(url)
		if err != nil {
			logger.Error(err, "Failed to get or fetch title", "url", url)
			continue
		}

		logger.V(2).Info("Title", "url", url, "title", title)

		_, err = fmt.Fprintf(&stringBuffer, "[%s](%s)\n\n", title, url)
		if err != nil {
			logger.Error(err, "Failed to write to string buffer")
		}
	}

	fmt.Println("Markdown-formatted URLs:\n" + stringBuffer.String())

	logger.V(1).Info("Debug: Exiting Hello function")
}
