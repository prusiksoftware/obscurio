package modify_queries

import (
	"errors"
	"fmt"
	pg_query "github.com/pganalyze/pg_query_go/v5"
	"github.com/prusiksoftware/monorepo/obscurio/psql_proxy/schema"
)

type WildcardExpander struct {
	tableDefinitions []schema.Table
	columnsToHide    map[string][]string
}

func NewWildcardExpander(tables []schema.Table, columnsToHide map[string][]string) *WildcardExpander {
	return &WildcardExpander{
		tableDefinitions: tables,
		columnsToHide:    columnsToHide,
	}
}

func (we *WildcardExpander) String() string {
	return "WildcardExpander"
}

func (we *WildcardExpander) _addTarget(parts []string, stmt *pg_query.SelectStmt) {
	var fields []*pg_query.Node
	for _, part := range parts {
		fields = append(fields, &pg_query.Node{
			Node: &pg_query.Node_String_{
				String_: &pg_query.String{Sval: part},
			},
		})
	}

	stmt.TargetList = append(stmt.TargetList, &pg_query.Node{
		Node: &pg_query.Node_ResTarget{
			ResTarget: &pg_query.ResTarget{
				Val: &pg_query.Node{
					Node: &pg_query.Node_ColumnRef{
						ColumnRef: &pg_query.ColumnRef{
							Fields: fields,
						},
					},
				},
			},
		},
	})

}

func (we *WildcardExpander) _getAllColumnNamesForTable(tableName string) []string {
	var columnNames []string
	for _, table := range we.tableDefinitions {
		if table.TableName == tableName {
			for _, column := range table.Columns {
				columnNames = append(columnNames, column.ColumnName)
			}
		}
	}
	return columnNames
}

func (we *WildcardExpander) columnHidden(columnName, tableName string) bool {
	hiddenColumns, ok := we.columnsToHide[tableName]
	if ok {
		for _, hiddenColumn := range hiddenColumns {
			if hiddenColumn == columnName {
				return true
			}
		}
	}
	return false
}

func (we *WildcardExpander) visit(stmt *pg_query.RawStmt) error {
	selectStatement, ok := stmt.Stmt.Node.(*pg_query.Node_SelectStmt)
	if ok {
		hadStar := false
		targetList := selectStatement.SelectStmt.TargetList
		for i, target := range targetList {
			resTarget := target.GetResTarget()
			if resTarget != nil {
				val := resTarget.GetVal()
				if val != nil {
					colRef := val.GetColumnRef()
					if colRef != nil {
						for _, field := range colRef.Fields {
							if field.GetAStar() != nil {
								hadStar = true
								selectStatement.SelectStmt.TargetList = append(targetList[:i], targetList[i+1:]...)
								break
							} else {
								colName := field.GetString_().Sval
								tableName := we.getTableName(stmt)
								isHidden := we.columnHidden(colName, tableName)
								if isHidden {
									return errors.New(fmt.Sprintf("column \"%s\" does not exist", colName))
								}
							}
						}
					}
				}
			}
			if hadStar {
				break
			}
		}

		if hadStar {
			tableName := we.getTableName(stmt)
			columnNames := we._getAllColumnNamesForTable(tableName)
			for _, fieldName := range columnNames {
				we._addTarget([]string{fieldName}, selectStatement.SelectStmt)
			}
		}
	}
	return nil
}

func (we *WildcardExpander) getTableName(stmt *pg_query.RawStmt) string {
	q, ok := stmt.Stmt.Node.(*pg_query.Node_SelectStmt)
	if ok {
		fromClause := q.SelectStmt.FromClause
		if fromClause != nil {
			for _, from := range fromClause {
				rangeVar, ok := from.Node.(*pg_query.Node_RangeVar)
				if ok {
					return rangeVar.RangeVar.Relname
				}
			}
		}
	}
	return ""
}
