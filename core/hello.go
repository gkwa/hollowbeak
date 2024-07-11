package core

import (
	"fmt"
	"strings"

	"github.com/go-logr/logr"
)

type URLInfo struct {
	URL   string
	Title string
}

func Hello(logger logr.Logger, filePath string, outputFormat string) error {
	logger.V(1).Info("Debug: Entering Hello function")

	titleFetcher := NewCollyTitleFetcher(logger)

	logger.V(1).Info("Debug: Creating new URLExtractor", "filePath", filePath)
	extractor, err := NewURLExtractor(logger, filePath, titleFetcher)
	if err != nil {
		return fmt.Errorf("failed to create URLExtractor: %w", err)
	}
	logger.V(2).Info("Debug: URLExtractor created")

	urlInfoList, err := BuildURLInfoList(logger, extractor)
	if err != nil {
		return fmt.Errorf("failed to build URL info list: %w", err)
	}

	switch outputFormat {
	case "markdown":
		PrintMarkdown(urlInfoList)
	case "html":
		PrintHTML(urlInfoList)
	default:
		return fmt.Errorf("invalid output format: %s", outputFormat)
	}

	logger.V(1).Info("Debug: Exiting Hello function")
	return nil
}

func BuildURLInfoList(logger logr.Logger, extractor *URLExtractor) ([]URLInfo, error) {
	logger.V(1).Info("Debug: Extracting URLs from file")
	urls, err := extractor.ExtractURLs()
	if err != nil {
		return nil, fmt.Errorf("failed to extract URLs: %w", err)
	}
	logger.V(2).Info("Debug: URLs extracted", "count", len(urls))

	var urlInfoList []URLInfo

	for _, url := range urls {
		logger.V(1).Info("Debug: Processing URL", "url", url)

		title, err := extractor.GetOrFetchTitle(url)
		if err != nil {
			logger.Error(err, "Failed to get or fetch title", "url", url)
			continue
		}

		logger.V(2).Info("Title", "url", url, "title", title)

		urlInfoList = append(urlInfoList, URLInfo{URL: url, Title: title})
	}

	return urlInfoList, nil
}

func PrintMarkdown(urlInfoList []URLInfo) {
	var sb strings.Builder
	for _, info := range urlInfoList {
		sb.WriteString(fmt.Sprintf("[%s](%s)\n\n", info.Title, info.URL))
	}
	fmt.Print(sb.String())
}

func PrintHTML(urlInfoList []URLInfo) {
	var sb strings.Builder
	sb.WriteString("<ul>\n")
	for _, info := range urlInfoList {
		sb.WriteString(fmt.Sprintf("  <li><a href=\"%s\">%s</a></li>\n", info.URL, info.Title))
	}
	sb.WriteString("</ul>")
	fmt.Println(sb.String())
}
