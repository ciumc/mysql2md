package main

import (
	"fmt"
	"net/url"
	"strings"
)

// escapeMD escapes special markdown characters in table cell content.
func escapeMD(s string) string {
	s = strings.ReplaceAll(s, "|", "\\|")
	s = strings.ReplaceAll(s, "\n", " ")
	return s
}

// RenderTableIndex generates a markdown table-of-contents listing all tables.
func RenderTableIndex(dbName string, tables []Table) string {
	var b strings.Builder

	fmt.Fprintf(&b, "# %s tables list\n", dbName)
	b.WriteString("| Name | Engine | Create_time | Collation | Comment |\n")
	b.WriteString("| ---- | ------ | ----------- | --------- | ------- |\n")

	for _, table := range tables {
		fileName := fmt.Sprintf("%s.%s.md", dbName, table.Name)
		fmt.Fprintf(&b, "| [%s](%s) | %s | %s | %s | %s |\n",
			escapeMD(table.Name),
			url.PathEscape(fileName),
			escapeMD(table.Engine),
			escapeMD(table.CreateTime),
			escapeMD(table.Collation),
			escapeMD(table.Comment),
		)
	}

	return b.String()
}

// RenderTableDetail generates the markdown documentation for a single table.
func RenderTableDetail(dbName string, table Table, columns []TableColumn, ddl *TableDDL, includeDDL bool) string {
	var b strings.Builder

	fmt.Fprintf(&b, "# %s.%s\n> %s\n", dbName, table.Name, escapeMD(table.Comment))
	b.WriteString("### COLUMNS\n")
	b.WriteString("| COLUMN_NAME | COLUMN_DEFAULT | IS_NULLABLE | COLLATION_NAME | COLUMN_TYPE | COLUMN_KEY | EXTRA | COLUMN_COMMENT |\n")
	b.WriteString("| ----------- | -------------- | ----------- | -------------- | ----------- | ---------- | ----- | -------------- |\n")

	for _, col := range columns {
		fmt.Fprintf(&b, "| %s | %s | %s | %s | %s | %s | %s | %s |\n",
			escapeMD(col.ColumnName),
			escapeMD(col.ColumnDefault),
			escapeMD(col.IsNullable),
			escapeMD(col.CollationName),
			escapeMD(col.ColumnType),
			escapeMD(col.ColumnKey),
			escapeMD(col.Extra),
			escapeMD(col.ColumnComment),
		)
	}

	if includeDDL && ddl != nil {
		b.WriteString("### DDL\n")
		fmt.Fprintf(&b, "```sql\n%s\n```\n", ddl.CreateTable)
	}

	return b.String()
}
