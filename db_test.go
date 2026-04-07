package main

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestEscapeIdentifier(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"normal name", "users", "`users`"},
		{"with backtick", "foo`bar", "`foo``bar`"},
		{"multiple backticks", "a`b`c", "`a``b``c`"},
		{"empty string", "", "``"},
		{"chinese name", "用户表", "`用户表`"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := escapeIdentifier(tt.input)
			if got != tt.expected {
				t.Errorf("escapeIdentifier(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestNewGormQuerierInvalidDSN(t *testing.T) {
	_, err := NewGormQuerier("invalid-dsn-that-will-fail")
	if err == nil {
		t.Error("expected error for invalid DSN, got nil")
	}
}

func TestSQLiteQuerierDatabaseName(t *testing.T) {
	q := newTestQuerier(t)
	name, err := q.DatabaseName()
	if err != nil {
		t.Fatalf("DatabaseName failed: %v", err)
	}
	if name != "testdb" {
		t.Errorf("expected 'testdb', got %q", name)
	}
}

func TestSQLiteQuerierTableList(t *testing.T) {
	q := newTestQuerier(t)
	tables, err := q.TableList()
	if err != nil {
		t.Fatalf("TableList failed: %v", err)
	}
	if len(tables) != 2 {
		t.Fatalf("expected 2 tables, got %d", len(tables))
	}

	names := map[string]bool{}
	for _, tbl := range tables {
		names[tbl.Name] = true
	}
	if !names["users"] || !names["orders"] {
		t.Errorf("expected tables 'users' and 'orders', got %v", names)
	}
}

func TestSQLiteQuerierTableColumns(t *testing.T) {
	q := newTestQuerier(t)
	columns, err := q.TableColumns("testdb", "users")
	if err != nil {
		t.Fatalf("TableColumns failed: %v", err)
	}
	if len(columns) != 4 {
		t.Fatalf("expected 4 columns for users, got %d", len(columns))
	}

	// Verify first column (id) is PRIMARY KEY
	if columns[0].ColumnName != "id" {
		t.Errorf("expected first column 'id', got %q", columns[0].ColumnName)
	}
	if columns[0].ColumnKey != "PRI" {
		t.Errorf("expected PRI key for id, got %q", columns[0].ColumnKey)
	}
	if columns[0].IsNullable != "NO" {
		t.Errorf("expected NO nullable for id (PK implies NOT NULL in SQLite), got %q", columns[0].IsNullable)
	}

	// Verify nullable column
	if columns[2].ColumnName != "email" {
		t.Errorf("expected third column 'email', got %q", columns[2].ColumnName)
	}
	if columns[2].IsNullable != "YES" {
		t.Errorf("expected YES nullable for email, got %q", columns[2].IsNullable)
	}
}

func TestSQLiteQuerierTableDDL(t *testing.T) {
	q := newTestQuerier(t)
	ddl, err := q.TableDDL("users")
	if err != nil {
		t.Fatalf("TableDDL failed: %v", err)
	}
	if ddl.Table != "users" {
		t.Errorf("expected table name 'users', got %q", ddl.Table)
	}
	if ddl.CreateTable == "" {
		t.Error("expected non-empty DDL")
	}
}

func TestSQLiteQuerierInterfaceCompliance(t *testing.T) {
	q := newTestQuerier(t)
	// Verify sqliteQuerier implements SchemaQuerier
	var _ SchemaQuerier = q
}

// newSQLiteGormQuerier creates a GormQuerier backed by SQLite for testing Close() and error paths.
func newSQLiteGormQuerier(t *testing.T) *GormQuerier {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to open sqlite: %v", err)
	}
	return &GormQuerier{db: db}
}

func TestGormQuerierClose(t *testing.T) {
	q := newSQLiteGormQuerier(t)
	err := q.Close()
	if err != nil {
		t.Fatalf("Close() failed: %v", err)
	}
}

