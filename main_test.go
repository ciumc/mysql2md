package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunEmptyDSN(t *testing.T) {
	oldDSN := *dsn
	defer func() { *dsn = oldDSN }()
	*dsn = ""

	err := run()
	if err == nil {
		t.Fatal("expected error for empty DSN, got nil")
	}
	if !strings.Contains(err.Error(), "dsn is empty") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestRunInvalidDSN(t *testing.T) {
	oldDSN := *dsn
	defer func() { *dsn = oldDSN }()
	*dsn = "invalid-dsn"

	err := run()
	if err == nil {
		t.Fatal("expected error for invalid DSN, got nil")
	}
}

func TestExecuteMultipleFiles(t *testing.T) {
	querier := newTestQuerier(t)
	dir := t.TempDir()

	err := execute(querier, dir, false, false)
	if err != nil {
		t.Fatalf("execute failed: %v", err)
	}

	// Verify index file
	indexPath := filepath.Join(dir, "testdb.tables.md")
	content, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatalf("failed to read index file: %v", err)
	}
	if !strings.Contains(string(content), "# testdb tables list") {
		t.Error("index file missing header")
	}

	// Verify table files exist
	files, _ := filepath.Glob(filepath.Join(dir, "testdb.*.md"))
	// Should have index + 2 table files = 3
	if len(files) < 3 {
		t.Errorf("expected at least 3 files, got %d", len(files))
	}
}

func TestExecuteWholeFile(t *testing.T) {
	querier := newTestQuerier(t)
	dir := t.TempDir()

	err := execute(querier, dir, true, false)
	if err != nil {
		t.Fatalf("execute whole failed: %v", err)
	}

	filePath := filepath.Join(dir, "testdb.tables.md")
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read whole file: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "# testdb.users") {
		t.Error("whole file missing users table")
	}
	if !strings.Contains(contentStr, "# testdb.orders") {
		t.Error("whole file missing orders table")
	}
}

func TestExecuteWholeFileWithDDL(t *testing.T) {
	querier := newTestQuerier(t)
	dir := t.TempDir()

	err := execute(querier, dir, true, true)
	if err != nil {
		t.Fatalf("execute whole with DDL failed: %v", err)
	}

	filePath := filepath.Join(dir, "testdb.tables.md")
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read whole file: %v", err)
	}
	if !strings.Contains(string(content), "### DDL") {
		t.Error("whole file missing DDL sections")
	}
}

func TestExecuteMultipleFilesWithDDL(t *testing.T) {
	querier := newTestQuerier(t)
	dir := t.TempDir()

	err := execute(querier, dir, false, true)
	if err != nil {
		t.Fatalf("execute multiple with DDL failed: %v", err)
	}

	// Check that a table file has DDL
	content, err := os.ReadFile(filepath.Join(dir, "testdb.users.md"))
	if err != nil {
		t.Fatalf("failed to read users file: %v", err)
	}
	if !strings.Contains(string(content), "### DDL") {
		t.Error("users file missing DDL section")
	}
}

func TestExecuteInvalidDir(t *testing.T) {
	querier := newTestQuerier(t)
	// /dev/null is not a directory, MkdirAll should fail
	err := execute(querier, "/dev/null/impossible", false, false)
	if err == nil {
		t.Error("expected error for invalid directory, got nil")
	}
}
