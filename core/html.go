package core

import (
	"io"
	"strings"

	"github.com/go-logr/logr"
	"golang.org/x/net/html"
)

func extractTitle(logger logr.Logger, reader io.Reader) (string, error) {
	tokenizer := html.NewTokenizer(reader)
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
