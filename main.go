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
