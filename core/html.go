package core

import (
	"fmt"
	"io"
	"strings"

	"github.com/go-logr/logr"
	"golang.org/x/net/html"
)

func extractTitle(logger logr.Logger, reader io.Reader) (string, error) {
	logger.V(1).Info("Debug: Entering extractTitle function")

	logger.V(2).Info("Debug: Creating HTML tokenizer")
	tokenizer := html.NewTokenizer(reader)
	logger.V(2).Info("Debug: HTML tokenizer created")

	for {
		logger.V(3).Info("Debug: Processing next token")
		tokenType := tokenizer.Next()
		logger.V(3).Info("Debug: Token type", "type", tokenType)

		switch tokenType {
		case html.ErrorToken:
			logger.V(2).Info("Debug: Reached end of HTML document or encountered an error")
			err := tokenizer.Err()
			if err == io.EOF {
				logger.V(2).Info("Debug: Reached end of HTML document without finding title")
				return "", fmt.Errorf("reached end of HTML document without finding title: %w", io.EOF)
			}
			logger.Error(err, "Error while tokenizing HTML")
			return "", fmt.Errorf("error while tokenizing HTML: %w", err)

		case html.StartTagToken, html.SelfClosingTagToken:
			token := tokenizer.Token()
			logger.V(3).Info("Debug: Found start/self-closing tag", "tag", token.Data)

			if token.Data == "title" {
				logger.V(2).Info("Debug: Found title tag")
				tokenType = tokenizer.Next()
				if tokenType == html.TextToken {
					title := strings.TrimSpace(tokenizer.Token().Data)
					logger.V(1).Info("Debug: Extracted title", "title", title)
					return title, nil
				}
				logger.V(2).Info("Debug: Title tag was empty or contained non-text content")
				return "", fmt.Errorf("title tag was empty or contained non-text content")
			}
		}
	}
}
