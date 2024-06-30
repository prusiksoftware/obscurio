package psql_proxy

import (
	wire "github.com/jeroenrinzema/psql-wire"
	"github.com/stretchr/testify/mock"
)

type MockQueryWriter struct {
	mock.Mock
	defineCalls      []wire.Columns
	rowCalls         [][]interface{}
	descriptionCalls []string
}

func (m *MockQueryWriter) Define(columns wire.Columns) error {
	m.defineCalls = append(m.defineCalls, columns)
	return nil
}

func (m *MockQueryWriter) Row(i []interface{}) error {
	m.rowCalls = append(m.rowCalls, i)
	return nil
}

func (m *MockQueryWriter) Empty() error {
	return nil
}

func (m *MockQueryWriter) Complete(description string) error {
	m.descriptionCalls = append(m.descriptionCalls, description)
	return nil
}
