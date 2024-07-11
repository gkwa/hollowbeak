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
	logger.V(1).Info("Debug: Creating new SQLTitleFetcher")
	return &SQLTitleFetcher{
		logger: logger,
	}
}

func (f *SQLTitleFetcher) FetchTitles(urls []urlRecord) (map[string]string, error) {
	f.logger.V(1).Info("Debug: Fetching titles from SQL database", "urlCount", len(urls))

	f.logger.V(2).Info("Debug: Getting titles for URLs")
	historyItems, query, err := f.getTitlesForURLs(urls)
	if err != nil {
		f.logger.Error(err, "Failed to fetch titles")
		return nil, fmt.Errorf("failed to fetch titles: %w", err)
	}
	f.logger.V(3).Info("Debug: SQL query executed", "query", query)

	titles := make(map[string]string)
	for _, url := range urls {
		if item, ok := historyItems[url.URL]; ok {
			f.logger.V(2).Info("Debug: Found title in database", "url", url.URL, "title", item.Title)
			titles[url.URL] = item.Title
		} else {
			f.logger.V(2).Info("Debug: No title found in database", "url", url.URL)
			titles[url.URL] = ""
		}
	}

	return titles, nil
}

func (f *SQLTitleFetcher) getTitlesForURLs(urls []urlRecord) (map[string]HistoryItem, string, error) {
	f.logger.V(2).Info("Debug: Getting titles for URLs", "urlCount", len(urls))

	historyFilePath, err := homedir.Expand("~/Library/Application Support/Google/Chrome/Default/History")
	if err != nil {
		f.logger.Error(err, "Failed to expand history file path")
		return nil, "", fmt.Errorf("failed to expand history file path: %v", err)
	}
	f.logger.V(3).Info("Debug: Chrome history file path", "path", historyFilePath)

	backupFile := filepath.Join(os.TempDir(), "history_backup.db")
	f.logger.V(3).Info("Debug: Creating history backup", "backupPath", backupFile)
	err = f.createHistoryBackup(historyFilePath, backupFile)
	if err != nil {
		f.logger.Error(err, "Failed to create history backup")
		return nil, "", err
	}
	defer os.Remove(backupFile)

	f.logger.V(3).Info("Debug: Opening SQLite database", "path", backupFile)
	db, err := sql.Open("sqlite3", backupFile+"?mode=ro")
	if err != nil {
		f.logger.Error(err, "Failed to open SQLite database")
		return nil, "", err
	}
	defer db.Close()

	// Debug SQLite queries if verbose mode is enabled
	if f.logger.V(4).Enabled() {
		err = DebugSQLiteQueries(f.logger, backupFile)
		if err != nil {
			f.logger.Error(err, "Failed to debug SQLite queries")
		}
	}

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

	f.logger.V(1).Info("Debug: Prepared SQL query", "query", query)

	args := make([]interface{}, len(urls))
	for i, url := range urls {
		args[i] = url.URL
	}

	f.logger.V(3).Info("Debug: Executing SQL query")
	rows, err := db.Query(query, args...)
	if err != nil {
		f.logger.Error(err, "Failed to execute SQL query")
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
			f.logger.Error(err, "Failed to scan SQL row")
			return nil, query, err
		}
		item.LastVisit, err = time.ParseInLocation("2006-01-02 15:04:05", visitTimeStr, time.Local)
		if err != nil {
			f.logger.Error(err, "Failed to parse visit time")
			return nil, query, err
		}
		item.RelativeVisit = f.formatRelativeTime(currentTime, item.LastVisit)
		item.Count = count
		count++

		historyItems[item.URL] = item
		f.logger.V(3).Info("Debug: Processed history item", "url", item.URL, "title", item.Title)
	}

	f.logger.V(2).Info("Debug: Finished getting titles for URLs", "itemCount", len(historyItems))
	return historyItems, query, nil
}

func (f *SQLTitleFetcher) createHistoryBackup(src, dst string) error {
	f.logger.V(2).Info("Debug: Creating history backup", "src", src, "dst", dst)
	input, err := os.ReadFile(src)
	if err != nil {
		f.logger.Error(err, "Failed to read source file")
		return err
	}

	err = os.WriteFile(dst, input, 0o644)
	if err != nil {
		f.logger.Error(err, "Failed to write destination file")
		return err
	}

	f.logger.V(2).Info("Debug: Successfully created history backup")
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
