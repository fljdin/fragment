package fragment

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var psql = Language{
	Delimiters: []string{`;`},
	Rules: []RangeRule{
		{Start: `--`, End: `\n`},
		{Start: `/*`, End: `*/`},
		{Start: `BEGIN`, End: `END`},
		{Start: `BEGIN`, End: `COMMIT`},
		{Start: `BEGIN`, End: `ROLLBACK`},
		{Start: `'`, End: `'`},
		{Start: `"`, End: `"`},
	},
}

func TestSemicolonDelimiter(t *testing.T) {
	input := "SELECT 1; SELECT 2; SELECT 3;"
	expected := []string{"SELECT 1;", "SELECT 2;", "SELECT 3;"}
	fragments := psql.Split(input)

	require.Equal(t, expected, fragments)
}
