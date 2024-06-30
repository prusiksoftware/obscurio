package modify_queries

import (
	pg_query "github.com/pganalyze/pg_query_go/v5"
)

type ColumnHider struct {
	columnsToHide map[string][]string
}

func NewColumnHider(columnsToHide map[string][]string) *ColumnHider {
	return &ColumnHider{
		columnsToHide: columnsToHide,
	}
}

func (ch *ColumnHider) String() string {
	return "ColumnHider"
}

func (ch *ColumnHider) visit(rawStatement *pg_query.RawStmt) error {
	selectStatement, ok := rawStatement.Stmt.Node.(*pg_query.Node_SelectStmt)
	if ok {
		tableNameAliases := map[string]string{}
		fromClause := selectStatement.SelectStmt.FromClause
		if fromClause != nil {
			for _, from := range fromClause {
				rangeVar, ok := from.Node.(*pg_query.Node_RangeVar)
				if ok {
					tableName := rangeVar.RangeVar.Relname
					tableNameAliases[tableName] = tableName
					// TODO: aliases
				}
			}
		}

		var newTargetList []*pg_query.Node
		for _, target := range selectStatement.SelectStmt.TargetList {
			shouldHide := false
			resTarget, ok := target.Node.(*pg_query.Node_ResTarget)
			if ok {
				node := resTarget.ResTarget.Val.Node
				columnRef, ok := node.(*pg_query.Node_ColumnRef)
				if ok {
					fields := columnRef.ColumnRef.Fields
					fieldName := fields[len(fields)-1].GetString_().GetSval()
					// TODO: handle aliases and table names

					tableName := ""
					if len(fields) == 1 {
						for _, fullName := range tableNameAliases {
							tableName = fullName
						}
					} else {
					}

					for hideTableName, hideColumnNames := range ch.columnsToHide {
						if hideTableName == tableName {
							for _, hideColumnName := range hideColumnNames {
								if hideColumnName == fieldName {
									shouldHide = true
								}
							}
						}
					}

				}
			}
			if !shouldHide {
				newTargetList = append(newTargetList, target)
			}
		}
		selectStatement.SelectStmt.TargetList = newTargetList
	}
	return nil
}
