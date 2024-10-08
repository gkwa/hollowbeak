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

func FetchURLTitles(
	logger logr.Logger,
	reader io.Reader,
	outputFormat string,
	fetcherTypes []string,
	noCache bool,
) error {
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

	logger.V(1).Info("Debug: Creating new URLExtractor")
	extractor, err := NewURLExtractor(logger, reader, titleFetchers, noCache)
	if err != nil {
		return fmt.Errorf("failed to create URLExtractor: %w", err)
	}
	logger.V(2).Info("Debug: URLExtractor created")

	defer func() {
		if !noCache {
			logger.V(1).Info("Debug: Cleaning up and saving cache")
			if err := extractor.cache.CleanupAndSave(); err != nil {
				logger.Error(err, "Failed to cleanup and save cache")
			}
		}
	}()

	urlInfoList, err := BuildURLInfoList(logger, extractor)
	if err != nil {
		return fmt.Errorf("failed to build URL info list: %w", err)
	}

	var output string
	switch outputFormat {
	case "markdown":
		output = GenerateMarkdown(urlInfoList)
	case "html":
		output = GenerateHTML(urlInfoList)
	case "space":
		output = GenerateSpaceDelimited(urlInfoList)
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

	titles, err := extractor.GetOrFetchTitles(urls)
	if err != nil {
		return nil, fmt.Errorf("failed to get or fetch titles: %w", err)
	}

	var urlInfoList []URLInfo
	for _, url := range urls {
		title := titles[url.URL]
		logger.V(2).Info("Title", "url", url.URL, "title", title)
		urlInfoList = append(urlInfoList, URLInfo{URL: url.URL, Title: title})
	}

	return urlInfoList, nil
}

func GenerateMarkdown(urlInfoList []URLInfo) string {
	var sb strings.Builder
	for _, info := range urlInfoList {
		sb.WriteString(fmt.Sprintf("[%s](%s)\n\n", info.Title, info.URL))
	}
	return sb.String()
}

func GenerateSpaceDelimited(urlInfoList []URLInfo) string {
	var sb strings.Builder
	for _, info := range urlInfoList {
		sb.WriteString(fmt.Sprintf("%s %s\n", info.URL, info.Title))
	}
	return sb.String()
}

func GenerateHTML(urlInfoList []URLInfo) string {
	var sb strings.Builder
	sb.WriteString("<ul>\n")
	for _, info := range urlInfoList {
		sb.WriteString(fmt.Sprintf("  <li><a href=\"%s\">%s</a></li>\n", info.URL, info.Title))
	}
	sb.WriteString("</ul>")
	return sb.String()
}
