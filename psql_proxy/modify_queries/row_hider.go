package modify_queries

import (
	pg_query "github.com/pganalyze/pg_query_go/v5"
)

type WhereOperator string

const (
	Equal              WhereOperator = "="
	NotEqual           WhereOperator = "!="
	LessThan           WhereOperator = "<"
	LessThanOrEqual    WhereOperator = "<="
	GreaterThan        WhereOperator = ">"
	GreaterThanOrEqual WhereOperator = ">="
)

type RowHider struct {
	tableName     string
	whereColumn   string
	whereOperator WhereOperator
	whereValue    string
}

func NewRowHider(tableName, col string, op WhereOperator, val string) *RowHider {
	return &RowHider{
		tableName:     tableName,
		whereColumn:   col,
		whereOperator: op,
		whereValue:    val,
	}
}

func (ch *RowHider) String() string {
	return "RowHider"
}

func (ch *RowHider) visit(rawStatement *pg_query.RawStmt) error {
	selectStatement, ok := rawStatement.Stmt.Node.(*pg_query.Node_SelectStmt)
	if ok {
		fromClause := selectStatement.SelectStmt.FromClause
		if fromClause != nil {
			for _, from := range fromClause {
				rangeVar, ok := from.Node.(*pg_query.Node_RangeVar)
				if ok {
					tableName := rangeVar.RangeVar.Relname
					if tableName == ch.tableName {
						whereClause := selectStatement.SelectStmt.WhereClause

						// new where clause
						if whereClause == nil {
							ch.addWhereClause(selectStatement.SelectStmt)
							return nil
						}

						// expand existing where clause
						ch.expandWhereClause(selectStatement.SelectStmt)
						return nil
					}
				}
			}
		}
	}
	return nil
}

func (ch *RowHider) whereExpr() pg_query.A_Expr {
	return pg_query.A_Expr{
		Kind: pg_query.A_Expr_Kind_AEXPR_OP,
		Name: []*pg_query.Node{
			{
				Node: &pg_query.Node_String_{
					String_: &pg_query.String{Sval: string(ch.whereOperator)},
				},
			},
		},
		Lexpr: &pg_query.Node{
			Node: &pg_query.Node_ColumnRef{
				ColumnRef: &pg_query.ColumnRef{
					Fields: []*pg_query.Node{
						{
							Node: &pg_query.Node_String_{
								String_: &pg_query.String{Sval: ch.tableName},
							},
						},
						{
							Node: &pg_query.Node_String_{
								String_: &pg_query.String{Sval: ch.whereColumn},
							},
						},
					},
				},
			},
		},
		Rexpr: &pg_query.Node{
			Node: &pg_query.Node_AConst{
				AConst: &pg_query.A_Const{
					Val: &pg_query.A_Const_Sval{
						Sval: &pg_query.String{Sval: ch.whereValue}},
				},
			},
		},
	}
}

func (ch *RowHider) addWhereClause(stmt *pg_query.SelectStmt) {
	where := ch.whereExpr()
	stmt.WhereClause = &pg_query.Node{
		Node: &pg_query.Node_AExpr{
			AExpr: &where,
		},
	}

}

func (ch *RowHider) expandWhereClause(stmt *pg_query.SelectStmt) {
	existingWhere := stmt.WhereClause.GetAExpr()
	where := ch.whereExpr()

	stmt.WhereClause = &pg_query.Node{
		Node: &pg_query.Node_BoolExpr{
			BoolExpr: &pg_query.BoolExpr{
				Boolop: pg_query.BoolExprType_AND_EXPR,
				Args: []*pg_query.Node{
					{
						Node: &pg_query.Node_AExpr{
							AExpr: existingWhere,
						},
					},
					{
						Node: &pg_query.Node_AExpr{
							AExpr: &where,
						},
					},
				},
			},
		},
	}

}
