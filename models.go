package main

// SchemaQuerier defines the interface for querying database schema information.
type SchemaQuerier interface {
	DatabaseName() (string, error)
	TableList() ([]Table, error)
	TableColumns(dbName, tableName string) ([]TableColumn, error)
	TableDDL(tableName string) (*TableDDL, error)
}

// Table represents a database table's metadata from SHOW TABLE STATUS.
type Table struct {
	Name       string `gorm:"column:Name"`
	Engine     string `gorm:"column:Engine"`
	Collation  string `gorm:"column:Collation"`
	Comment    string `gorm:"column:Comment"`
	CreateTime string `gorm:"column:Create_time"`
}

// TableDDL represents a table's DDL (CREATE TABLE statement).
type TableDDL struct {
	Table       string `gorm:"column:Table"`
	CreateTable string `gorm:"column:Create Table"`
}

// TableColumn represents a column's metadata from information_schema.columns.
type TableColumn struct {
	TableCatalog           string `gorm:"column:TABLE_CATALOG"`
	TableSchema            string `gorm:"column:TABLE_SCHEMA"`
	TableName              string `gorm:"column:TABLE_NAME"`
	ColumnName             string `gorm:"column:COLUMN_NAME"`
	OrdinalPosition        int    `gorm:"column:ORDINAL_POSITION"`
	ColumnDefault          string `gorm:"column:COLUMN_DEFAULT"`
	IsNullable             string `gorm:"column:IS_NULLABLE"`
	DataType               string `gorm:"column:DATA_TYPE"`
	CharacterMaximumLength int    `gorm:"column:CHARACTER_MAXIMUM_LENGTH"`
	CharacterOctetLength   int    `gorm:"column:CHARACTER_OCTET_LENGTH"`
	NumericPrecision       int    `gorm:"column:NUMERIC_PRECISION"`
	NumericScale           int    `gorm:"column:NUMERIC_SCALE"`
	DatetimePrecision      int    `gorm:"column:DATETIME_PRECISION"`
	CharacterSetName       string `gorm:"column:CHARACTER_SET_NAME"`
	CollationName          string `gorm:"column:COLLATION_NAME"`
	ColumnType             string `gorm:"column:COLUMN_TYPE"`
	ColumnKey              string `gorm:"column:COLUMN_KEY"`
	Extra                  string `gorm:"column:EXTRA"`
	Privileges             string `gorm:"column:PRIVILEGES"`
	ColumnComment          string `gorm:"column:COLUMN_COMMENT"`
	GenerationExpression   string `gorm:"column:GENERATION_EXPRESSION"`
}
