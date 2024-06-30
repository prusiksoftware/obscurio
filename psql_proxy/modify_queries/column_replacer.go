package modify_queries

import (
	pg_query "github.com/pganalyze/pg_query_go/v5"
)

type ColumnReplacer struct {
	tableName string
	values    map[string]string
}

func NewColumnReplacer(tableName string, values map[string]string) *ColumnReplacer {
	return &ColumnReplacer{
		tableName: tableName,
		values:    values,
	}
}

func (cr *ColumnReplacer) String() string {
	return "ColumnReplacer"
}

func (cr *ColumnReplacer) visit(stmt *pg_query.RawStmt) error {
	selectStmt := stmt.Stmt.GetSelectStmt()
	if selectStmt == nil {
		return nil
	}

	for _, target := range selectStmt.TargetList {
		resTarget, ok := target.Node.(*pg_query.Node_ResTarget)
		if ok {
			node := resTarget.ResTarget.Val.Node
			columnRef, ok := node.(*pg_query.Node_ColumnRef)
			if ok {
				fields := columnRef.ColumnRef.Fields
				fieldName := fields[len(fields)-1].GetString_().GetSval()
				if fieldName == "" {
					continue
				}

				tableName := ""
				if len(fields) == 1 {
					tableName = cr.tableName
				}

				if tableName == cr.tableName {
					if newValue, ok := cr.values[fieldName]; ok {
						resTarget.ResTarget.Name = fieldName
						resTarget.ResTarget.Val = &pg_query.Node{
							Node: &pg_query.Node_AConst{
								AConst: &pg_query.A_Const{
									Val: &pg_query.A_Const_Sval{
										Sval: &pg_query.String{Sval: newValue},
									},
								},
							},
						}
					}
				}
			}
		}
	}
	return nil
}
