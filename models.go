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
	TableSchema            string `gorm:"column:TABLE_SCHEMA"`              // Database name
	TableName              string `gorm:"column:TABLE_NAME"`                // Table name
	ColumnName             string `gorm:"column:COLUMN_NAME"`               // Column name
	OrdinalPosition        int    `gorm:"column:ORDINAL_POSITION"`           // Column position in table
	ColumnDefault          string `gorm:"column:COLUMN_DEFAULT"`             // Default value
	IsNullable             string `gorm:"column:IS_NULLABLE"`                // YES or NO
	DataType               string `gorm:"column:DATA_TYPE"`                  // Data type name
	CharacterMaximumLength int    `gorm:"column:CHARACTER_MAXIMUM_LENGTH"`   // Max length for string types
	CharacterOctetLength   int    `gorm:"column:CHARACTER_OCTET_LENGTH"`     // Max length in bytes
	NumericPrecision       int    `gorm:"column:NUMERIC_PRECISION"`          // Precision for numeric types
	NumericScale           int    `gorm:"column:NUMERIC_SCALE"`              // Scale for numeric types
	DatetimePrecision      int    `gorm:"column:DATETIME_PRECISION"`        // Precision for datetime types
	CharacterSetName       string `gorm:"column:CHARACTER_SET_NAME"`
	CollationName          string `gorm:"column:COLLATION_NAME"`            // Collation name
	ColumnType             string `gorm:"column:COLUMN_TYPE"`                // Full column type definition
	ColumnKey              string `gorm:"column:COLUMN_KEY"`                 // Key type: PRI, UNI, MUL
	Extra                  string `gorm:"column:EXTRA"`                     // Additional info like auto_increment
	Privileges             string `gorm:"column:PRIVILEGES"`                 // Column privileges
	ColumnComment          string `gorm:"column:COLUMN_COMMENT"`             // Column comment
	GenerationExpression   string `gorm:"column:GENERATION_EXPRESSION"`      // Expression for generated columns
}
