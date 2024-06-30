package modify_queries

import (
	"github.com/prusiksoftware/monorepo/obscurio/psql_proxy/schema"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestQueryModifier_Query(t *testing.T) {
	t.Run("TestQueryModifier_Query", func(t *testing.T) {
		abcUserTable := schema.Table{
			TableName: "users",
			Columns: []schema.Column{
				{
					ColumnName: "a",
				},
				{
					ColumnName: "b",
				},
				{
					ColumnName: "c",
				},
			},
		}

		we := NewWildcardExpander([]schema.Table{
			abcUserTable,
		}, map[string][]string{})
		qm, err := NewQueryModifier("SELECT * FROM users", []ModifierInterface{we})
		require.NoError(t, err)

		err = qm.Modify()
		query, err2 := qm.Query()

		require.NoError(t, err)
		require.NoError(t, err2)
		require.Equal(t, "SELECT a, b, c FROM users", query)
	})
}
