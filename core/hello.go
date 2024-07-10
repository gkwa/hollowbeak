package core

import (
	"github.com/go-logr/logr"
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
