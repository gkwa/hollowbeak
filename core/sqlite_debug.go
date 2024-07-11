package core

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/go-logr/logr"
)

func DebugSQLiteQueries(logger logr.Logger, dbPath string) error {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	var tableName string
	err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='sqlite_stat4'").Scan(&tableName)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Info("sqlite_stat4 table does not exist. Query history might not be available.")
			return nil
		}
		return fmt.Errorf("failed to check for sqlite_stat4 table: %w", err)
	}

	rows, err := db.Query("SELECT * FROM sqlite_stat4")
	if err != nil {
		return fmt.Errorf("failed to query sqlite_stat4 table: %w", err)
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
