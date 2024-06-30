package modify_queries

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestColumnHider(t *testing.T) {

	t.Run("TestColumnHider a,b,c -> a,c", func(t *testing.T) {
		ch := NewColumnHider(map[string][]string{
			"users": {"b"},
		})
		qm, err := NewQueryModifier("SELECT a, b, c FROM users", []ModifierInterface{ch})
		require.NoError(t, err)

		err = qm.Modify()
		query, err2 := qm.Query()

		require.NoError(t, err)
		require.NoError(t, err2)
		require.Equal(t, "SELECT a, c FROM users", query)
	})

	t.Run("TestColumnHider a,b,c -> a,b,c", func(t *testing.T) {
		ch := NewColumnHider(map[string][]string{
			"users":      {"d"},
			"otherTable": {"a"},
		})
		qm, err := NewQueryModifier("SELECT a, b, c FROM users", []ModifierInterface{ch})
		require.NoError(t, err)

		err = qm.Modify()
		query, err2 := qm.Query()

		require.NoError(t, err)
		require.NoError(t, err2)
		require.Equal(t, "SELECT a, b, c FROM users", query)
	})

}
