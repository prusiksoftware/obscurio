package modify_queries

import (
	pg_query "github.com/pganalyze/pg_query_go/v5"
)

type QueryModifier struct {
	originalQuery string
	ast           *pg_query.ParseResult
	visitors      []ModifierInterface
}

type ModifierInterface interface {
	visit(*pg_query.RawStmt) error
	String() string
}

func NewQueryModifier(query string, visitors []ModifierInterface) (*QueryModifier, error) {
	parseResult, err := pg_query.Parse(query)
	if err != nil {
		return nil, err
	}
	return &QueryModifier{
		originalQuery: query,
		ast:           parseResult,
		visitors:      visitors,
	}, nil
}

func (qc *QueryModifier) Modify() error {
	for _, visitor := range qc.visitors {
		for _, stmt := range qc.ast.Stmts {
			err := visitor.visit(stmt)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (qc *QueryModifier) Query() (string, error) {
	res, err := pg_query.Deparse(qc.ast)
	if err != nil {
		return "", err
	}
	return res, nil
}
