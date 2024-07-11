package core

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-logr/logr"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mitchellh/go-homedir"
)

type SQLTitleFetcher struct {
	logger logr.Logger
}

func NewSQLTitleFetcher(logger logr.Logger) *SQLTitleFetcher {
	return &SQLTitleFetcher{
		logger: logger,
	}
}

func (f *SQLTitleFetcher) FetchTitle(url string) (string, error) {
	f.logger.V(1).Info("Debug: Fetching title from SQL database", "url", url)

	historyItems, _, err := f.getTitlesForURLs([]string{url})
	if err != nil {
		return "", fmt.Errorf("failed to fetch titles: %w", err)
	}

	if item, ok := historyItems[url]; ok {
		f.logger.V(2).Info("Debug: Found title in database", "url", url, "title", item.Title)
		return item.Title, nil
	}

	return "", fmt.Errorf("no title found for URL: %s", url)
}

func (f *SQLTitleFetcher) getTitlesForURLs(urls []string) (map[string]HistoryItem, string, error) {
	historyFilePath, err := homedir.Expand("~/Library/Application Support/Google/Chrome/Default/History")
	if err != nil {
		return nil, "", fmt.Errorf("failed to expand history file path: %v", err)
	}

	backupFile := filepath.Join(os.TempDir(), "history_backup.db")
	err = f.createHistoryBackup(historyFilePath, backupFile)
	if err != nil {
		return nil, "", err
	}
	defer os.Remove(backupFile)

	db, err := sql.Open("sqlite3", backupFile+"?mode=ro")
	if err != nil {
		return nil, "", err
	}
	defer db.Close()

	placeholders := make([]string, len(urls))
	for i := range placeholders {
		placeholders[i] = "?"
	}

	query := fmt.Sprintf(`
   SELECT
   	datetime(visits.visit_time/1000000-11644473600, 'unixepoch', 'localtime') as visit_time,
   	urls.url,
   	urls.title
   FROM
   	visits INNER JOIN urls ON visits.url = urls.id
   WHERE
   	urls.url IN (%s)
   ORDER BY
   	visit_time DESC
  `, strings.Join(placeholders, ","))

	args := make([]interface{}, len(urls))
	for i, url := range urls {
		args[i] = url
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, query, err
	}
	defer rows.Close()

	historyItems := make(map[string]HistoryItem)
	count := 1

	currentTime := time.Now()
	for rows.Next() {
		var item HistoryItem
		var visitTimeStr string
		err := rows.Scan(&visitTimeStr, &item.URL, &item.Title)
		if err != nil {
			return nil, query, err
		}
		item.LastVisit, err = time.ParseInLocation("2006-01-02 15:04:05", visitTimeStr, time.Local)
		if err != nil {
			return nil, query, err
		}
		item.RelativeVisit = f.formatRelativeTime(currentTime, item.LastVisit)
		item.Count = count
		count++

		historyItems[item.URL] = item
	}

	return historyItems, query, nil
}

func (f *SQLTitleFetcher) createHistoryBackup(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	err = os.WriteFile(dst, input, 0o644)
	if err != nil {
		return err
	}

	return nil
}

func (f *SQLTitleFetcher) formatRelativeTime(currentTime, visitTime time.Time) string {
	duration := currentTime.Sub(visitTime)

	if duration < time.Minute {
		return fmt.Sprintf("%d seconds ago", int(duration.Seconds()))
	} else if duration < time.Hour {
		minutes := int(duration.Minutes())
		return fmt.Sprintf("%d %s ago", minutes, f.pluralize("minute", minutes))
	} else if duration < 24*time.Hour {
		hours := int(duration.Hours())
		return fmt.Sprintf("%d %s ago", hours, f.pluralize("hour", hours))
	} else {
		days := int(duration.Hours() / 24)
		return fmt.Sprintf("%d %s ago", days, f.pluralize("day", days))
	}
}

func (f *SQLTitleFetcher) pluralize(word string, count int) string {
	if count == 1 {
		return word
	}
	return word + "s"
}

type HistoryItem struct {
	URL           string
	Title         string
	LastVisit     time.Time
	RelativeVisit string
	Count         int
}
