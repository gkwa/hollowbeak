package core

import (
	"fmt"

	"github.com/go-logr/logr"
	"github.com/gocolly/colly/v2"
)

type CollyTitleFetcher struct {
	logger logr.Logger
}

func NewCollyTitleFetcher(logger logr.Logger) *CollyTitleFetcher {
	return &CollyTitleFetcher{
		logger: logger,
	}
}

func (f *CollyTitleFetcher) FetchTitle(url string) (string, error) {
	f.logger.V(1).Info("Debug: Fetching title with Colly", "url", url)

	c := colly.NewCollector(
		colly.AllowURLRevisit(),
		colly.MaxDepth(5),
	)

	// Set Chrome-like User-Agent
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"

	var title string
	var finalURL string

	c.OnHTML("title", func(e *colly.HTMLElement) {
		title = e.Text
		finalURL = e.Request.URL.String()
		f.logger.V(2).Info("Debug: Found title", "title", title, "url", finalURL)
	})

	c.OnError(func(r *colly.Response, err error) {
		f.logger.Error(err, "Colly encountered an error", "url", r.Request.URL.String(), "statusCode", r.StatusCode)
	})

	c.OnResponse(func(r *colly.Response) {
		if r.Request.URL.String() != url {
			f.logger.V(2).Info("Debug: Followed redirect", "from", url, "to", r.Request.URL.String())
		}
	})

	err := c.Visit(url)
	if err != nil {
		return "", fmt.Errorf("failed to visit URL: %w", err)
	}

	if title == "" {
		return "", fmt.Errorf("no title found for URL: %s", url)
	}

	f.logger.V(1).Info("Debug: Successfully fetched title", "originalURL", url, "finalURL", finalURL, "title", title)
	return title, nil
}
