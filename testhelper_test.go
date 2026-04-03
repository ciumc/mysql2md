package main

import (
	"fmt"

	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// sqliteQuerier implements SchemaQuerier using SQLite in-memory database for testing.
type sqliteQuerier struct {
	db     *gorm.DB
	dbName string
}

func (q *sqliteQuerier) DatabaseName() (string, error) {
	return q.dbName, nil
}

func (q *sqliteQuerier) TableList() ([]Table, error) {
	type sqliteMaster struct {
		Name string `gorm:"column:name"`
	}
	var masters []sqliteMaster
	err := q.db.Raw("SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%' ORDER BY name").Scan(&masters).Error
	if err != nil {
		return nil, fmt.Errorf("listing tables: %w", err)
	}

	tables := make([]Table, len(masters))
	for i, m := range masters {
		tables[i] = Table{
			Name:       m.Name,
			Engine:     "SQLite",
			Collation:  "utf8",
			Comment:    "",
			CreateTime: "2024-01-01 00:00:00",
		}
	}
	return tables, nil
}

func (q *sqliteQuerier) TableColumns(dbName, tableName string) ([]TableColumn, error) {
	type pragmaColumn struct {
		CID        int    `gorm:"column:cid"`
		Name       string `gorm:"column:name"`
		Type       string `gorm:"column:type"`
		NotNull    int    `gorm:"column:notnull"`
		DefaultVal *string `gorm:"column:dflt_value"`
		PK         int    `gorm:"column:pk"`
	}
	var pragmaCols []pragmaColumn
	err := q.db.Raw(fmt.Sprintf("PRAGMA table_info(`%s`)", tableName)).Scan(&pragmaCols).Error
	if err != nil {
		return nil, fmt.Errorf("querying columns for table %s: %w", tableName, err)
	}

	columns := make([]TableColumn, len(pragmaCols))
	for i, pc := range pragmaCols {
		isNullable := "YES"
		if pc.NotNull == 1 || pc.PK > 0 {
			isNullable = "NO"
		}
		columnKey := ""
		if pc.PK > 0 {
			columnKey = "PRI"
		}
		defaultVal := ""
		if pc.DefaultVal != nil {
			defaultVal = *pc.DefaultVal
		}
		columns[i] = TableColumn{
			ColumnName:    pc.Name,
			ColumnType:    pc.Type,
			ColumnDefault: defaultVal,
			IsNullable:    isNullable,
			ColumnKey:     columnKey,
			CollationName: "",
			Extra:         "",
			ColumnComment: "",
		}
	}
	return columns, nil
}

func (q *sqliteQuerier) TableDDL(tableName string) (*TableDDL, error) {
	type sqliteMaster struct {
		SQL string `gorm:"column:sql"`
	}
	var master sqliteMaster
	err := q.db.Raw("SELECT sql FROM sqlite_master WHERE type='table' AND name=?", tableName).Scan(&master).Error
	if err != nil {
		return nil, fmt.Errorf("querying DDL for table %s: %w", tableName, err)
	}
	return &TableDDL{
		Table:       tableName,
		CreateTable: master.SQL,
	}, nil
}

// newTestQuerier creates a SQLite in-memory database with test tables and returns a sqliteQuerier.
func newTestQuerier(t *testing.T) *sqliteQuerier {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to open sqlite: %v", err)
	}

	// Create test tables
	sqls := []string{
		`CREATE TABLE users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL DEFAULT '',
			email TEXT,
			created_at TEXT
		)`,
		`CREATE TABLE orders (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			amount REAL DEFAULT 0,
			status TEXT DEFAULT 'pending'
		)`,
	}

	for _, sql := range sqls {
		if err := db.Exec(sql).Error; err != nil {
			t.Fatalf("failed to create table: %v", err)
		}
	}

	return &sqliteQuerier{db: db, dbName: "testdb"}
}
