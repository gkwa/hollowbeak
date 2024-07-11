package core

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/go-logr/logr"
)

type URLInfo struct {
	URL   string
	Title string
}

func Hello(logger logr.Logger, filePath string, outputFormat string, fetcherTypes []string, noCache bool) error {
	logger.V(1).Info("Debug: Entering Hello function")

	var titleFetchers []TitleFetcher
	for _, fetcherType := range fetcherTypes {
		switch fetcherType {
		case "http":
			titleFetchers = append(titleFetchers, NewHTTPTitleFetcher(logger))
		case "colly":
			titleFetchers = append(titleFetchers, NewCollyTitleFetcher(logger))
		case "sql":
			titleFetchers = append(titleFetchers, NewSQLTitleFetcher(logger))
		default:
			return fmt.Errorf("invalid fetcher type: %s", fetcherType)
		}
	}

	if len(titleFetchers) == 0 {
		return fmt.Errorf("no valid fetcher types specified")
	}

	logger.V(1).Info("Debug: Creating new URLExtractor", "filePath", filePath)
	extractor, err := NewURLExtractor(logger, filePath, titleFetchers, noCache)
	if err != nil {
		return fmt.Errorf("failed to create URLExtractor: %w", err)
	}
	logger.V(2).Info("Debug: URLExtractor created")

	urlInfoList, err := BuildURLInfoList(logger, extractor)
	if err != nil {
		return fmt.Errorf("failed to build URL info list: %w", err)
	}

	var output string
	switch outputFormat {
	case "markdown":
		output = PrintMarkdown(urlInfoList)
	case "html":
		output = PrintHTML(urlInfoList)
	default:
		return fmt.Errorf("invalid output format: %s", outputFormat)
	}

	_, err = io.WriteString(os.Stdout, output)
	if err != nil {
		return fmt.Errorf("failed to write output: %w", err)
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

func PrintMarkdown(urlInfoList []URLInfo) string {
	var sb strings.Builder
	for _, info := range urlInfoList {
		sb.WriteString(fmt.Sprintf("[%s](%s)\n\n", info.Title, info.URL))
	}
	return sb.String()
}

func PrintHTML(urlInfoList []URLInfo) string {
	var sb strings.Builder
	sb.WriteString("<ul>\n")
	for _, info := range urlInfoList {
		sb.WriteString(fmt.Sprintf("  <li><a href=\"%s\">%s</a></li>\n", info.URL, info.Title))
	}
	sb.WriteString("</ul>")
	return sb.String()
}
