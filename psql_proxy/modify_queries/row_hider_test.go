package modify_queries

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRowHider(t *testing.T) {

	t.Run("a,b,c -> where b = '2'", func(t *testing.T) {
		ch := NewRowHider("users", "b", Equal, "2")
		qm, err := NewQueryModifier("SELECT a, b, c FROM users", []ModifierInterface{ch})
		require.NoError(t, err)

		err = qm.Modify()
		query, err2 := qm.Query()

		require.NoError(t, err)
		require.NoError(t, err2)
		require.Equal(t, "SELECT a, b, c FROM users WHERE users.b = '2'", query)
	})

	t.Run("a,b,c where c = '5' -> where b = '2'", func(t *testing.T) {
		ch := NewRowHider("users", "b", Equal, "2")
		qm, err := NewQueryModifier("SELECT a, b, c FROM users WHERE c = '5'", []ModifierInterface{ch})
		require.NoError(t, err)

		err = qm.Modify()
		query, err2 := qm.Query()

		require.NoError(t, err)
		require.NoError(t, err2)
		require.Equal(t, "SELECT a, b, c FROM users WHERE c = '5' AND users.b = '2'", query)
	})

	t.Run("a,b,c where c = '5' -> no change", func(t *testing.T) {
		ch := NewRowHider("other_table", "b", Equal, "2")
		qm, err := NewQueryModifier("SELECT a, b, c FROM users WHERE c = '5'", []ModifierInterface{ch})
		require.NoError(t, err)

		err = qm.Modify()
		query, err2 := qm.Query()

		require.NoError(t, err)
		require.NoError(t, err2)
		require.Equal(t, "SELECT a, b, c FROM users WHERE c = '5'", query)
	})

}
