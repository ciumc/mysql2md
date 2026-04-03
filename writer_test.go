package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteMultipleFiles(t *testing.T) {
	querier := newTestQuerier(t)
	dir := t.TempDir()

	tables, err := querier.TableList()
	if err != nil {
		t.Fatalf("failed to list tables: %v", err)
	}

	err = WriteMultipleFiles(querier, dir, "testdb", tables, false)
	if err != nil {
		t.Fatalf("WriteMultipleFiles failed: %v", err)
	}

	// Check index file exists
	indexPath := filepath.Join(dir, "testdb.tables.md")
	indexContent, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatalf("failed to read index file: %v", err)
	}

	indexStr := string(indexContent)
	if !strings.Contains(indexStr, "# testdb tables list") {
		t.Error("index file missing header")
	}
	if !strings.Contains(indexStr, "[orders](testdb.orders.md)") {
		t.Error("index file missing orders link")
	}
	if !strings.Contains(indexStr, "[users](testdb.users.md)") {
		t.Error("index file missing users link")
	}

	// Check individual table files
	for _, table := range tables {
		filePath := filepath.Join(dir, "testdb."+table.Name+".md")
		content, err := os.ReadFile(filePath)
		if err != nil {
			t.Fatalf("failed to read file for table %s: %v", table.Name, err)
		}
		contentStr := string(content)
		if !strings.Contains(contentStr, "# testdb."+table.Name) {
			t.Errorf("table file %s missing header", table.Name)
		}
		if !strings.Contains(contentStr, "### COLUMNS") {
			t.Errorf("table file %s missing COLUMNS section", table.Name)
		}
		// Should not have DDL
		if strings.Contains(contentStr, "### DDL") {
			t.Errorf("table file %s should not have DDL section", table.Name)
		}
	}
}

func TestWriteMultipleFilesWithDDL(t *testing.T) {
	querier := newTestQuerier(t)
	dir := t.TempDir()

	tables, err := querier.TableList()
	if err != nil {
		t.Fatalf("failed to list tables: %v", err)
	}

	err = WriteMultipleFiles(querier, dir, "testdb", tables, true)
	if err != nil {
		t.Fatalf("WriteMultipleFiles with DDL failed: %v", err)
	}

	for _, table := range tables {
		filePath := filepath.Join(dir, "testdb."+table.Name+".md")
		content, err := os.ReadFile(filePath)
		if err != nil {
			t.Fatalf("failed to read file for table %s: %v", table.Name, err)
		}
		contentStr := string(content)
		if !strings.Contains(contentStr, "### DDL") {
			t.Errorf("table file %s missing DDL section", table.Name)
		}
		if !strings.Contains(contentStr, "```sql") {
			t.Errorf("table file %s missing SQL code block", table.Name)
		}
	}
}

func TestWriteWholeFile(t *testing.T) {
	querier := newTestQuerier(t)
	dir := t.TempDir()

	tables, err := querier.TableList()
	if err != nil {
		t.Fatalf("failed to list tables: %v", err)
	}

	err = WriteWholeFile(querier, dir, "testdb", tables, false)
	if err != nil {
		t.Fatalf("WriteWholeFile failed: %v", err)
	}

	filePath := filepath.Join(dir, "testdb.tables.md")
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read whole file: %v", err)
	}

	contentStr := string(content)

	// Should contain all tables
	for _, table := range tables {
		if !strings.Contains(contentStr, "# testdb."+table.Name) {
			t.Errorf("whole file missing table %s", table.Name)
		}
	}

	// Tables should appear in the same order as input
	ordersIdx := strings.Index(contentStr, "# testdb.orders")
	usersIdx := strings.Index(contentStr, "# testdb.users")
	if ordersIdx == -1 || usersIdx == -1 {
		t.Fatal("missing table entries in whole file")
	}
	if ordersIdx > usersIdx {
		t.Error("tables not in expected order (orders should come before users)")
	}

	// Should not have DDL
	if strings.Contains(contentStr, "### DDL") {
		t.Error("whole file should not have DDL section")
	}
}

func TestWriteWholeFileWithDDL(t *testing.T) {
	querier := newTestQuerier(t)
	dir := t.TempDir()

	tables, err := querier.TableList()
	if err != nil {
		t.Fatalf("failed to list tables: %v", err)
	}

	err = WriteWholeFile(querier, dir, "testdb", tables, true)
	if err != nil {
		t.Fatalf("WriteWholeFile with DDL failed: %v", err)
	}

	filePath := filepath.Join(dir, "testdb.tables.md")
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read whole file: %v", err)
	}

	contentStr := string(content)
	ddlCount := strings.Count(contentStr, "### DDL")
	if ddlCount != len(tables) {
		t.Errorf("expected %d DDL sections, got %d", len(tables), ddlCount)
	}
}

func TestWriteWholeFileRace(t *testing.T) {
	querier := newTestQuerier(t)
	dir := t.TempDir()

	// Create many tables for a more meaningful race test
	for i := 0; i < 20; i++ {
		tableName := fmt.Sprintf("race_table_%c", 'a'+i)
		_ = querier.db.Exec("CREATE TABLE " + tableName + " (id INTEGER PRIMARY KEY, val TEXT)").Error
	}

	tables, err := querier.TableList()
	if err != nil {
		t.Fatalf("failed to list tables: %v", err)
	}

	// Run with -race flag to detect data races
	err = WriteWholeFile(querier, dir, "testdb", tables, true)
	if err != nil {
		t.Fatalf("WriteWholeFile failed: %v", err)
	}

	filePath := filepath.Join(dir, "testdb.tables.md")
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read whole file: %v", err)
	}

	// Verify all tables are present
	for _, table := range tables {
		if !strings.Contains(string(content), "# testdb."+table.Name) {
			t.Errorf("whole file missing table %s", table.Name)
		}
	}
}

func TestWriteErrorPropagation(t *testing.T) {
	querier := newTestQuerier(t)
	// Use a non-existent directory to trigger an error
	dir := filepath.Join(t.TempDir(), "nonexistent", "subdir")

	tables, err := querier.TableList()
	if err != nil {
		t.Fatalf("failed to list tables: %v", err)
	}

	err = WriteMultipleFiles(querier, dir, "testdb", tables, false)
	if err == nil {
		t.Error("expected error for non-existent directory, got nil")
	}

	err = WriteWholeFile(querier, dir, "testdb", tables, false)
	if err == nil {
		t.Error("expected error for non-existent directory, got nil")
	}
}
