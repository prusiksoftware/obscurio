package modify_queries

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewColumnReplacer(t *testing.T) {
	cr := NewColumnReplacer("users", map[string]string{"a": "1", "b": "2", "c": "3"})
	qm, err := NewQueryModifier("SELECT a, b, c FROM users", []ModifierInterface{cr})
	require.NoError(t, err)

	err = qm.Modify()
	query, err2 := qm.Query()

	require.NoError(t, err)
	require.NoError(t, err2)
	require.NotNil(t, cr)
	require.Equal(t, "SELECT '1' AS a, '2' AS b, '3' AS c FROM users", query)
}
