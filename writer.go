package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/schollz/progressbar/v3"
	"golang.org/x/sync/errgroup"
)

// tableResult holds the query result for a single table, used by the gather-then-write pattern.
type tableResult struct {
	table   Table
	columns []TableColumn
	ddl     *TableDDL
}

// WriteMultipleFiles generates one markdown file per table plus an index file.
func WriteMultipleFiles(querier SchemaQuerier, dir, dbName string, tables []Table, includeDDL bool) error {
	// Write index file
	indexContent := RenderTableIndex(dbName, tables)
	indexPath := filepath.Join(dir, fmt.Sprintf("%s.tables.md", dbName))
	if err := os.WriteFile(indexPath, []byte(indexContent), 0644); err != nil {
		return fmt.Errorf("writing index file %s: %w", indexPath, err)
	}

	bar := progressbar.Default(int64(len(tables)))
	defer bar.Close()

	wg := errgroup.Group{}
	wg.SetLimit(runtime.NumCPU())

	for _, table := range tables {
		table := table
		wg.Go(func() error {
			columns, err := querier.TableColumns(dbName, table.Name)
			if err != nil {
				return err
			}

			var ddl *TableDDL
			if includeDDL {
				ddl, err = querier.TableDDL(table.Name)
				if err != nil {
					return err
				}
			}

			content := RenderTableDetail(dbName, table, columns, ddl, includeDDL)
			filePath := filepath.Join(dir, fmt.Sprintf("%s.%s.md", dbName, table.Name))
			if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
				return fmt.Errorf("writing table file %s: %w", filePath, err)
			}

			_ = bar.Add(1)
			return nil
		})
	}

	return wg.Wait()
}

// WriteWholeFile generates a single markdown file containing all tables.
// Uses the gather-then-write pattern: concurrent DB queries, sequential file write.
func WriteWholeFile(querier SchemaQuerier, dir, dbName string, tables []Table, includeDDL bool) error {
	bar := progressbar.Default(int64(len(tables)))
	defer bar.Close()

	// Gather: concurrent DB queries, results stored by index (no race)
	results := make([]tableResult, len(tables))

	wg := errgroup.Group{}
	wg.SetLimit(runtime.NumCPU())

	for i, table := range tables {
		i, table := i, table
		wg.Go(func() error {
			columns, err := querier.TableColumns(dbName, table.Name)
			if err != nil {
				return err
			}

			var ddl *TableDDL
			if includeDDL {
				ddl, err = querier.TableDDL(table.Name)
				if err != nil {
					return err
				}
			}

			results[i] = tableResult{table: table, columns: columns, ddl: ddl}
			_ = bar.Add(1)
			return nil
		})
	}

	if err := wg.Wait(); err != nil {
		return err
	}

	// Write: sequential file write (no race)
	filePath := filepath.Join(dir, fmt.Sprintf("%s.tables.md", dbName))
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("creating file %s: %w", filePath, err)
	}
	defer f.Close()

	for _, r := range results {
		content := RenderTableDetail(dbName, r.table, r.columns, r.ddl, includeDDL)
		if _, err := f.WriteString(content); err != nil {
			return fmt.Errorf("writing to file %s: %w", filePath, err)
		}
	}

	return nil
}
