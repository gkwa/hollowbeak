package core

import (
	"net/http"

	"github.com/go-logr/logr"
)

func getPageTitle(logger logr.Logger, url string) (string, error) {
	logger.V(1).Info("Debug: Entering getPageTitle function", "url", url)

	resp, err := http.Get(url)
	if err != nil {
		logger.Error(err, "Failed to make HTTP request")
		return "", err
	}
	defer resp.Body.Close()

	logger.V(1).Info("Debug: HTTP request successful", "status", resp.Status)

	return extractTitle(logger, resp.Body)
}
