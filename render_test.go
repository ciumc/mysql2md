package main

import (
	"strings"
	"testing"
)

func TestEscapeMD(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"empty string", "", ""},
		{"no special chars", "hello world", "hello world"},
		{"pipe char", "foo|bar", "foo\\|bar"},
		{"newline", "foo\nbar", "foo bar"},
		{"pipe and newline", "a|b\nc", "a\\|b c"},
		{"multiple pipes", "a|b|c", "a\\|b\\|c"},
		{"chinese with pipe", "用户|表", "用户\\|表"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := escapeMD(tt.input)
			if got != tt.expected {
				t.Errorf("escapeMD(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestRenderTableIndex(t *testing.T) {
	tables := []Table{
		{Name: "users", Engine: "InnoDB", CreateTime: "2024-01-01", Collation: "utf8mb4_general_ci", Comment: "用户表"},
		{Name: "orders", Engine: "InnoDB", CreateTime: "2024-01-02", Collation: "utf8mb4_general_ci", Comment: "订单表"},
	}

	result := RenderTableIndex("mydb", tables)

	// Check header
	if !strings.Contains(result, "# mydb tables list\n") {
		t.Error("missing header")
	}

	// Check table headers
	if !strings.Contains(result, "| Name | Engine | Create_time | Collation | Comment |") {
		t.Error("missing table header row")
	}

	// Check separator
	if !strings.Contains(result, "| ---- | ------ | ----------- | --------- | ------- |") {
		t.Error("missing separator row")
	}

	// Check table rows with links
	if !strings.Contains(result, "| [users](mydb.users.md) | InnoDB | 2024-01-01 | utf8mb4_general_ci | 用户表 |") {
		t.Errorf("missing or incorrect users row, got:\n%s", result)
	}
	if !strings.Contains(result, "| [orders](mydb.orders.md) | InnoDB | 2024-01-02 | utf8mb4_general_ci | 订单表 |") {
		t.Errorf("missing or incorrect orders row, got:\n%s", result)
	}
}

func TestRenderTableIndexEmpty(t *testing.T) {
	result := RenderTableIndex("mydb", []Table{})

	if !strings.Contains(result, "# mydb tables list\n") {
		t.Error("missing header for empty table list")
	}

	lines := strings.Split(strings.TrimSpace(result), "\n")
	if len(lines) != 3 {
		t.Errorf("expected 3 lines (header + table header + separator), got %d", len(lines))
	}
}

func TestRenderTableIndexWithSpecialChars(t *testing.T) {
	tables := []Table{
		{Name: "test", Engine: "InnoDB", CreateTime: "2024-01-01", Collation: "utf8mb4", Comment: "has|pipe\nand newline"},
	}

	result := RenderTableIndex("mydb", tables)

	if strings.Contains(result, "has|pipe") {
		t.Error("pipe character in comment was not escaped")
	}
	if !strings.Contains(result, "has\\|pipe and newline") {
		t.Errorf("expected escaped comment, got:\n%s", result)
	}
}

func TestRenderTableDetail(t *testing.T) {
	table := Table{Name: "users", Comment: "用户表"}
	columns := []TableColumn{
		{
			ColumnName:    "id",
			ColumnDefault: "",
			IsNullable:    "NO",
			CollationName: "",
			ColumnType:    "bigint",
			ColumnKey:     "PRI",
			Extra:         "auto_increment",
			ColumnComment: "主键ID",
		},
		{
			ColumnName:    "name",
			ColumnDefault: "",
			IsNullable:    "YES",
			CollationName: "utf8mb4_general_ci",
			ColumnType:    "varchar(255)",
			ColumnKey:     "",
			Extra:         "",
			ColumnComment: "用户名",
		},
	}
	ddl := &TableDDL{
		Table:       "users",
		CreateTable: "CREATE TABLE `users` (\n  `id` bigint NOT NULL AUTO_INCREMENT,\n  `name` varchar(255) DEFAULT NULL,\n  PRIMARY KEY (`id`)\n) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4",
	}

	t.Run("without DDL", func(t *testing.T) {
		result := RenderTableDetail("mydb", table, columns, ddl, false)

		if !strings.Contains(result, "# mydb.users\n> 用户表\n") {
			t.Errorf("missing or incorrect header, got:\n%s", result)
		}
		if !strings.Contains(result, "### COLUMNS\n") {
			t.Error("missing COLUMNS header")
		}
		if !strings.Contains(result, "| id |  | NO |  | bigint | PRI | auto_increment | 主键ID |") {
			t.Errorf("missing id column row, got:\n%s", result)
		}
		if !strings.Contains(result, "| name |  | YES | utf8mb4_general_ci | varchar(255) |  |  | 用户名 |") {
			t.Errorf("missing name column row, got:\n%s", result)
		}
		if strings.Contains(result, "### DDL") {
			t.Error("DDL section should not be present when includeDDL=false")
		}
	})

	t.Run("with DDL", func(t *testing.T) {
		result := RenderTableDetail("mydb", table, columns, ddl, true)

		if !strings.Contains(result, "### DDL\n") {
			t.Error("missing DDL header")
		}
		if !strings.Contains(result, "```sql\n") {
			t.Error("missing SQL code block")
		}
		if !strings.Contains(result, "CREATE TABLE `users`") {
			t.Error("missing DDL content")
		}
	})

	t.Run("with nil DDL and includeDDL true", func(t *testing.T) {
		result := RenderTableDetail("mydb", table, columns, nil, true)

		if strings.Contains(result, "### DDL") {
			t.Error("DDL section should not be present when ddl is nil")
		}
	})
}

func TestRenderTableDetailEmptyColumns(t *testing.T) {
	table := Table{Name: "empty_table", Comment: ""}
	result := RenderTableDetail("mydb", table, []TableColumn{}, nil, false)

	if !strings.Contains(result, "# mydb.empty_table\n") {
		t.Error("missing header")
	}

	// Should have header rows but no data rows
	lines := strings.Split(strings.TrimSpace(result), "\n")
	expectedLines := 5 // title, comment, section header, table header, separator
	if len(lines) != expectedLines {
		t.Errorf("expected %d lines, got %d:\n%s", expectedLines, len(lines), result)
	}
}

func TestRenderTableIndexURLEscape(t *testing.T) {
	tables := []Table{
		{Name: "my table", Engine: "InnoDB", CreateTime: "2024-01-01", Collation: "utf8mb4", Comment: ""},
		{Name: "foo(bar)", Engine: "InnoDB", CreateTime: "2024-01-01", Collation: "utf8mb4", Comment: ""},
	}

	result := RenderTableIndex("mydb", tables)

	// URL should have spaces encoded
	if strings.Contains(result, "(mydb.my table.md)") {
		t.Error("space in URL should be escaped")
	}
	// URL should have parentheses encoded
	if strings.Contains(result, "(mydb.foo(bar).md)") {
		t.Error("parentheses in URL should be escaped")
	}
	// Display name should still show the original (with pipe escaping)
	if !strings.Contains(result, "[my table]") {
		t.Errorf("display name should contain original table name, got:\n%s", result)
	}
}

func TestRenderTableDetailWithSpecialCharsInComment(t *testing.T) {
	table := Table{Name: "test", Comment: "表|注释\n换行"}
	columns := []TableColumn{
		{
			ColumnName:    "col1",
			ColumnComment: "字段|注释\n换行",
		},
	}

	result := RenderTableDetail("mydb", table, columns, nil, false)

	if strings.Contains(result, "表|注释") {
		t.Error("pipe in table comment was not escaped")
	}
	if !strings.Contains(result, "表\\|注释 换行") {
		t.Errorf("table comment not correctly escaped, got:\n%s", result)
	}
	if !strings.Contains(result, "字段\\|注释 换行") {
		t.Errorf("column comment not correctly escaped, got:\n%s", result)
	}
}
