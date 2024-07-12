package core

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/go-logr/logr"
)

func DebugSQLiteQueries(logger logr.Logger, db *sql.DB) error {
	logger.Info("Debug: Setting up SQLite query logging")

	_, err := db.Exec("PRAGMA query_log = on;")
	if err != nil {
		logger.Error(err, "Failed to enable query logging")
		return fmt.Errorf("failed to enable query logging: %w", err)
	}

	_, err = db.Exec(`
   CREATE TABLE IF NOT EXISTS query_log (
   	time DATETIME DEFAULT CURRENT_TIMESTAMP,
   	query TEXT
   )
   `)
	if err != nil {
		logger.Error(err, "Failed to create query_log table")
		return fmt.Errorf("failed to create query_log table: %w", err)
	}

	logger.Info("Debug: Retrieving SQLite query history")

	rows, err := db.Query("SELECT * FROM query_log ORDER BY time DESC LIMIT 10")
	if err != nil {
		return fmt.Errorf("failed to query query_log: %w", err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return fmt.Errorf("failed to get column names: %w", err)
	}

	logger.Info("Query History:")
	logger.Info(strings.Join(columns, "\t"))

	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))

	for rows.Next() {
		for i := range columns {
			valuePtrs[i] = &values[i]
		}
		err := rows.Scan(valuePtrs...)
		if err != nil {
			return fmt.Errorf("failed to scan row: %w", err)
		}

		rowValues := make([]string, len(columns))
		for i, col := range values {
			if col != nil {
				rowValues[i] = fmt.Sprintf("%v", col)
			} else {
				rowValues[i] = "NULL"
			}
		}
		logger.Info(strings.Join(rowValues, "\t"))
	}

	if err = rows.Err(); err != nil {
		return fmt.Errorf("error during row iteration: %w", err)
	}

	return nil
}
