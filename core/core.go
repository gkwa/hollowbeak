package core

import (
	"io"
	"net/http"
	"strings"

	"github.com/go-logr/logr"
	"golang.org/x/net/html"
)

func Hello(logger logr.Logger) {
	logger.V(1).Info("Debug: Entering Hello function")
	logger.Info("Hello, World!")
	logger.V(1).Info("Debug: Exiting Hello function")

	url := "https://open.substack.com/pub/systemdesignone/p/saga-design-pattern?utm_source=share&utm_medium=android&r=21036"
	logger.V(1).Info("Debug: About to fetch page title", "url", url)
	title, err := getPageTitle(logger, url)
	if err != nil {
		logger.Error(err, "Failed to get page title")
		return
	}
	logger.Info("Fetched page title", "title", title)
}

func getPageTitle(logger logr.Logger, url string) (string, error) {
	logger.V(1).Info("Debug: Entering getPageTitle function", "url", url)

	resp, err := http.Get(url)
	if err != nil {
		logger.Error(err, "Failed to make HTTP request")
		return "", err
	}
	defer resp.Body.Close()

	logger.V(1).Info("Debug: HTTP request successful", "status", resp.Status)

	tokenizer := html.NewTokenizer(resp.Body)
	logger.V(1).Info("Debug: Created HTML tokenizer")

	for {
		tokenType := tokenizer.Next()
		logger.V(2).Info("Debug: Processing token", "type", tokenType)

		switch tokenType {
		case html.ErrorToken:
			logger.V(1).Info("Debug: Reached end of HTML document")
			return "", io.EOF
		case html.StartTagToken, html.SelfClosingTagToken:
			token := tokenizer.Token()
			logger.V(2).Info("Debug: Found start/self-closing tag", "tag", token.Data)

			if token.Data == "title" {
				logger.V(1).Info("Debug: Found title tag")
				tokenType = tokenizer.Next()
				if tokenType == html.TextToken {
					title := strings.TrimSpace(tokenizer.Token().Data)
					logger.V(1).Info("Debug: Extracted title", "title", title)
					return title, nil
				}
			}
		}
	}
}
