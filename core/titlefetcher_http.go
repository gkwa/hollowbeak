package core

import (
	"fmt"
	"net/http"

	"github.com/go-logr/logr"
)

type TitleFetcher interface {
	FetchTitle(url string) (string, error)
}

type HTTPTitleFetcher struct {
	logger logr.Logger
	client *http.Client
}

func NewHTTPTitleFetcher(logger logr.Logger) *HTTPTitleFetcher {
	logger.V(1).Info("Debug: Creating new HTTPTitleFetcher")
	return &HTTPTitleFetcher{
		logger: logger,
		client: &http.Client{},
	}
}

func (f *HTTPTitleFetcher) FetchTitle(url string) (string, error) {
	f.logger.V(1).Info("Debug: Fetching title", "url", url)

	f.logger.V(2).Info("Debug: Creating HTTP request", "url", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		f.logger.Error(err, "Failed to create HTTP request", "url", url)
		return "", fmt.Errorf("failed to create HTTP request: %w", err)
	}

	f.logger.V(2).Info("Debug: Setting User-Agent header")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	f.logger.V(2).Info("Debug: Sending HTTP request", "url", url)
	resp, err := f.client.Do(req)
	if err != nil {
		f.logger.Error(err, "Failed to make HTTP request", "url", url)
		return "", fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	f.logger.V(2).Info("Debug: HTTP request successful", "url", url, "status", resp.Status)

	if resp.StatusCode >= 300 && resp.StatusCode < 400 {
		f.logger.V(2).Info("Debug: Encountered redirect", "url", url, "status", resp.Status, "location", resp.Header.Get("Location"))
	}

	f.logger.V(2).Info("Debug: Extracting title from response body", "url", url)
	title, err := extractTitle(f.logger, resp.Body)
	if err != nil {
		f.logger.Error(err, "Failed to extract title", "url", url)
		return "", fmt.Errorf("failed to extract title: %w", err)
	}

	f.logger.V(1).Info("Debug: Successfully fetched title", "url", url, "title", title)
	return title, nil
}
