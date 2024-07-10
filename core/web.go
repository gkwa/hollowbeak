package core

import (
	"net/http"

	"github.com/go-logr/logr"
)

func getPageTitle(logger logr.Logger, url string) (string, error) {
	logger.V(1).Info("Debug: Entering getPageTitle function", "url", url)

	logger.V(2).Info("Debug: Creating HTTP client")
	client := &http.Client{}

	logger.V(2).Info("Debug: Creating HTTP request", "url", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Error(err, "Failed to create HTTP request")
		return "", err
	}

	logger.V(2).Info("Debug: Setting User-Agent header")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	logger.V(2).Info("Debug: Sending HTTP request")
	resp, err := client.Do(req)
	if err != nil {
		logger.Error(err, "Failed to make HTTP request")
		return "", err
	}
	defer resp.Body.Close()

	logger.V(1).Info("Debug: HTTP request successful", "status", resp.Status)

	logger.V(2).Info("Debug: Extracting title from response body")
	return extractTitle(logger, resp.Body)
}
