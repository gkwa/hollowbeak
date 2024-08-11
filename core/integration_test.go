package core

import (
	"os"
	"testing"

	"github.com/go-logr/logr/testr"
)

func TestIntegration(t *testing.T) {
	logger := testr.New(t)

	markdown := `
[Token consumption in Microsoft’s Graph RAG – baeke.info](https://blog.baeke.info/2024/07/11/token-consumption-in-microsofts-graph-rag/)
`

	tempFile, err := createTempFileWithURL(markdown)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	fetchers := []string{"sql", "colly", "http"}
	noCache := true

	err = FetchURLTitles(logger, tempFile, "markdown", fetchers, noCache)
	if err != nil {
		t.Fatalf("Hello function failed: %v", err)
	}
}

func createTempFileWithURL(url string) (*os.File, error) {
	tempFile, err := os.CreateTemp("", "test_urls_*.txt")
	if err != nil {
		return nil, err
	}

	_, err = tempFile.WriteString(url)
	if err != nil {
		tempFile.Close()
		os.Remove(tempFile.Name())
		return nil, err
	}

	return tempFile, nil
}
