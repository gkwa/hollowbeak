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
	logger.V(1).Info("Debug: Creating new CollyTitleFetcher")
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
	f.logger.V(2).Info("Debug: Created Colly collector", "maxDepth", 5)

	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"
	f.logger.V(2).Info("Debug: Set User-Agent for Colly collector", "userAgent", c.UserAgent)

	var title string
	var finalURL string

	c.OnHTML("title", func(e *colly.HTMLElement) {
		title = e.Text
		finalURL = e.Request.URL.String()
		f.logger.V(2).Info("Debug: Found title", "title", title, "url", finalURL)
	})

	c.OnRequest(func(r *colly.Request) {
		f.logger.V(3).Info("Debug: Colly making request", "url", r.URL.String())
	})

	c.OnResponse(func(r *colly.Response) {
		f.logger.V(3).Info("Debug: Colly received response", "url", r.Request.URL.String(), "statusCode", r.StatusCode)
		if r.Request.URL.String() != url {
			f.logger.V(2).Info("Debug: Followed redirect", "from", url, "to", r.Request.URL.String())
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		f.logger.Error(err, "Colly encountered an error", "url", r.Request.URL.String(), "statusCode", r.StatusCode)
	})

	f.logger.V(2).Info("Debug: Starting Colly visit", "url", url)
	err := c.Visit(url)
	if err != nil {
		f.logger.Error(err, "Failed to visit URL with Colly", "url", url)
		return "", fmt.Errorf("failed to visit URL: %w", err)
	}

	if title == "" {
		f.logger.V(2).Info("Debug: No title found", "url", url)
		return "", fmt.Errorf("no title found for URL: %s", url)
	}

	f.logger.V(1).Info("Debug: Successfully fetched title", "originalURL", url, "finalURL", finalURL, "title", title)
	return title, nil
}
