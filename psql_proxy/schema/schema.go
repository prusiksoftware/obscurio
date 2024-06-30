package schema

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5"
	_ "github.com/lib/pq"
	"github.com/lib/pq/oid"
	"log"
	"strings"
)

type Column struct {
	TableCatalog           sql.NullString `json:"table_catalog"`
	TableSchema            sql.NullString `json:"table_schema"`
	TableName              sql.NullString `json:"table_name"`
	ColumnName             string         `json:"column_name"`
	OrdinalPosition        int            `json:"ordinal_position"`
	ColumnDefault          sql.NullString `json:"column_default"`
	IsNullable             sql.NullString `json:"is_nullable"`
	DataType               sql.NullString `json:"data_type"`
	CharacterMaximumLength *int16         `json:"character_maximum_length,omitempty"`
	CharacterOctetLength   *int           `json:"character_octet_length,omitempty"`
	NumericPrecision       *int           `json:"numeric_precision,omitempty"`
	NumericPrecisionRadix  *int           `json:"numeric_precision_radix,omitempty"`
	NumericScale           *int           `json:"numeric_scale,omitempty"`
	DatetimePrecision      *int           `json:"datetime_precision,omitempty"`
	CharacterSetName       sql.NullString `json:"character_set_name"`
	CollationName          sql.NullString `json:"collation_name"`
	oidType                oid.Oid
}

type Table struct {
	TableCatalog              string         `json:"table_catalog"`
	TableSchema               string         `json:"table_schema"`
	TableName                 string         `json:"table_name"`
	TableType                 sql.NullString `json:"table_type"`
	SelfReferencingColumnName sql.NullString `json:"self_referencing_column_name,omitempty"`
	ReferenceGeneration       sql.NullString `json:"reference_generation,omitempty"`
	UserDefinedTypeCatalog    sql.NullString `json:"user_defined_type_catalog,omitempty"`
	UserDefinedTypeSchema     sql.NullString `json:"user_defined_type_schema,omitempty"`
	UserDefinedTypeName       sql.NullString `json:"user_defined_type_name,omitempty"`
	IsInsertableInto          sql.NullString `json:"is_insertable_into,omitempty"`
	IsTyped                   sql.NullString `json:"is_typed,omitempty"`
	CommitAction              sql.NullString `json:"commit_action,omitempty"`
	Columns                   []Column       `json:"columns,omitempty"`
}

func GetColumns(db *pgx.Conn, tableName string) []Column {

	columnNames := []string{
		"table_catalog",
		"table_schema",
		"table_name",
		"column_name",
		"ordinal_position",
		"column_default",
		"is_nullable",
		"data_type",
		"character_maximum_length",
		"character_octet_length",
		"numeric_precision",
		"numeric_precision_radix",
		"numeric_scale",
		"datetime_precision",
		"character_set_name",
		"collation_name",
		"pgc.oid",
	}
	columnList := strings.Join(columnNames, ", ")
	query := fmt.Sprintf(`
        SELECT %s
        FROM information_schema.columns
        JOIN pg_class pgc ON information_schema.columns.table_name = pgc.relname
        WHERE table_name = '%s'`, columnList, tableName)

	rows, err := db.Query(context.Background(), query)
	defer rows.Close()
	if err != nil {
		panic(err)
	}

	var columns []Column
	for rows.Next() {
		var column Column
		err := rows.Scan(
			&column.TableCatalog,
			&column.TableSchema,
			&column.TableName,
			&column.ColumnName,
			&column.OrdinalPosition,
			&column.ColumnDefault,
			&column.IsNullable,
			&column.DataType,
			&column.CharacterMaximumLength,
			&column.CharacterOctetLength,
			&column.NumericPrecision,
			&column.NumericPrecisionRadix,
			&column.NumericScale,
			&column.DatetimePrecision,
			&column.CharacterSetName,
			&column.CollationName,
			&column.oidType,
		)
		if err != nil {
			log.Fatal(err)
		}

		columns = append(columns, column)
	}
	return columns
}

func GetTables(dburi string) []Table {
	config, err := pgx.Connect(context.Background(), dburi)
	if err != nil {
		log.Fatal(err)
	}
	defer config.Close(context.Background())

	conn, err := pgx.ConnectConfig(context.Background(), config.Config())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(context.Background())

	columnNames := []string{
		"table_catalog",
		"table_schema",
		"table_name",
		"table_type",
		"self_referencing_column_name",
		"reference_generation",
		"user_defined_type_catalog",
		"user_defined_type_schema",
		"user_defined_type_name",
		"is_insertable_into",
		"is_typed",
		"commit_action",
	}

	columnList := strings.Join(columnNames, ", ")
	query := fmt.Sprintf(`
        SELECT %s
        FROM information_schema.tables
        ORDER BY table_name`, columnList)

	rows, err := conn.Query(context.Background(), query)
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}

	var tables []Table
	for rows.Next() {
		var table Table
		err := rows.Scan(
			&table.TableCatalog,
			&table.TableSchema,
			&table.TableName,
			&table.TableType,
			&table.SelfReferencingColumnName,
			&table.ReferenceGeneration,
			&table.UserDefinedTypeCatalog,
			&table.UserDefinedTypeSchema,
			&table.UserDefinedTypeName,
			&table.IsInsertableInto,
			&table.IsTyped,
			&table.CommitAction,
		)
		if err != nil {
			panic(err)
		}
		tables = append(tables, table)
	}

	for i, table := range tables {
		tables[i].Columns = GetColumns(conn, table.TableName)
	}

	return tables
}

func (t *Table) PrettyPrint() {
	b, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
}
