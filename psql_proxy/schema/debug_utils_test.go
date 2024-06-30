package schema

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetColumns(t *testing.T) {
	t.Skip("skipping test")
	t.Run("default", func(t *testing.T) {
		tables := GetTables("postgres://user:password@localhost:5432/chinook?sslmode=disable")

		expected := []Table{}
		require.Equal(t, 11, len(tables))
		require.Equal(t, expected, tables)
	})
}
