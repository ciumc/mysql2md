package main

import (
	"fmt"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// GormQuerier implements SchemaQuerier using GORM with a MySQL backend.
type GormQuerier struct {
	db *gorm.DB
}

// NewGormQuerier creates a new GormQuerier by opening a MySQL connection with the given DSN.
func NewGormQuerier(dsn string) (*GormQuerier, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("connecting to database: %w", err)
	}
	return &GormQuerier{db: db}, nil
}

// escapeIdentifier wraps a SQL identifier in backticks, escaping any internal backticks.
func escapeIdentifier(name string) string {
	return "`" + strings.ReplaceAll(name, "`", "``") + "`"
}

// Close closes the underlying database connection.
func (q *GormQuerier) Close() error {
	sqlDB, err := q.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// DatabaseName returns the name of the currently connected database.
func (q *GormQuerier) DatabaseName() (string, error) {
	var name string
	err := q.db.Raw("SELECT DATABASE()").Scan(&name).Error
	if err != nil {
		return "", fmt.Errorf("querying database name: %w", err)
	}
	return name, nil
}

// TableList returns metadata for all tables in the current database.
func (q *GormQuerier) TableList() ([]Table, error) {
	var tables []Table
	err := q.db.Raw("SHOW TABLE STATUS").Scan(&tables).Error
	if err != nil {
		return nil, fmt.Errorf("listing tables: %w", err)
	}
	return tables, nil
}

// TableColumns returns column metadata for a specific table.
func (q *GormQuerier) TableColumns(dbName, tableName string) ([]TableColumn, error) {
	var columns []TableColumn
	err := q.db.Raw(
		"SELECT * FROM information_schema.columns WHERE table_schema = ? AND table_name = ?",
		dbName, tableName,
	).Scan(&columns).Error
	if err != nil {
		return nil, fmt.Errorf("querying columns for table %s: %w", tableName, err)
	}
	return columns, nil
}

// TableDDL retrieves the CREATE TABLE statement for a given table.
// NOTE: We must use string concatenation here instead of parameterized queries
// because MySQL's "SHOW CREATE TABLE" statement does not support placeholders.
// The table name is escaped via escapeIdentifier() to prevent SQL injection.
func (q *GormQuerier) TableDDL(tableName string) (*TableDDL, error) {
	ddl := &TableDDL{}
	err := q.db.Raw(fmt.Sprintf("SHOW CREATE TABLE %s", escapeIdentifier(tableName))).Scan(ddl).Error
	if err != nil {
		return nil, fmt.Errorf("querying DDL for table %s: %w", tableName, err)
	}
	return ddl, nil
}
