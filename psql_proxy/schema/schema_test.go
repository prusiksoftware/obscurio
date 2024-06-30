package schema

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetTables(t *testing.T) {
	tests := []struct {
		name string
		want []Table
	}{
		//{
		//	name: "test1",
		//	want: []Table{},
		//},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := PrettyPrints(GetTables("postgres://user:password@proxy:5432/chinook?sslmode=disable"))
			filename := fmt.Sprintf("./golden_files/schema_test/testGetTable_%s.txt", tt.name)
			matchesGoldenFile(t, actual, filename)
		})
	}
}

func matchesGoldenFile(t *testing.T, actual string, filename string) {
	//writetoFile(actual, filename)
	goldenContent := readFromFile(filename)
	require.Equal(t, goldenContent, actual)
}
