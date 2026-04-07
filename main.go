package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	dsn     = flag.String("dsn", "username:password@tcp(localhost:3306)/database?charset=utf8mb4&parseTime=True&loc=Local&timeout=10s", "database connection string")
	dir     = flag.String("dir", "./output", "directory to save the file")
	isWhole = flag.Bool("whole", false, "generate whole file (default false)")
	isDDL   = flag.Bool("ddl", false, "generate ddl info (default false)")
)

// run is the main entry point that parses flags and executes the schema documentation generation.
func run() error {
	flag.Parse()

	if *dsn == "" {
		return fmt.Errorf("dsn is empty")
	}

	querier, err := NewGormQuerier(*dsn)
	if err != nil {
		return err
	}
	defer querier.Close()

	return execute(querier, *dir, *isWhole, *isDDL)
}

// execute orchestrates the schema documentation generation process.
// It creates the output directory, queries the database metadata, and writes markdown files.
// When whole is true, a single file is generated; otherwise, separate files per table plus an index.
func execute(querier SchemaQuerier, outputDir string, whole, includeDDL bool) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("creating output directory %s: %w", outputDir, err)
	}

	dbName, err := querier.DatabaseName()
	if err != nil {
		return err
	}

	tables, err := querier.TableList()
	if err != nil {
		return err
	}

	if whole {
		return WriteWholeFile(querier, outputDir, dbName, tables, includeDDL)
	}

	return WriteMultipleFiles(querier, outputDir, dbName, tables, includeDDL)
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}
}
