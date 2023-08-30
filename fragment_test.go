package fragment

import (
	"testing"

	"github.com/lithammer/dedent"
	"github.com/stretchr/testify/require"
)

var psql = Language{
	Delimiters: []string{`;`},
	Rules: []RangeRule{
		{Start: `--`, End: `\n`},
		{Start: `/*`, End: `*/`},
		{Start: `'`, End: `'`},
		{Start: `"`, End: `"`},
		// {Start: `BEGIN`, End: `END`},
		// {Start: `BEGIN`, End: `COMMIT`},
		// {Start: `BEGIN`, End: `ROLLBACK`},
	},
}

func TestIgnoreEmptyFragments(t *testing.T) {
	input := " "
	fragments := psql.Split(input)
	require.Equal(t, 0, len(fragments))
}

func TestSemicolonDelimiter(t *testing.T) {
	input := "SELECT 1; SELECT 2; SELECT 3;"
	expected := []string{"SELECT 1;", "SELECT 2;", "SELECT 3;"}
	fragments := psql.Split(input)

	for i, fragment := range fragments {
		require.Equal(t, expected[i], fragment)
	}
}

func TestCommentRules(t *testing.T) {
	input := dedent.Dedent(`
		SELECT 1 -- comment;
		 + 1;
		SELECT 1 /* comment; */ + 1;
		SELECT 1,
		 /* multi-line
		  * comment ;
		  */
		 2;
		SELECT 1 /* -- comment ; */ + 1;
		SELECT 1 -- /* comment ;
		 + 1;
	`)
	expected := []string{
		"SELECT 1 -- comment;\n + 1;",
		"SELECT 1 /* comment; */ + 1;",
		"SELECT 1,\n /* multi-line\n  * comment ;\n  */\n 2;",
		"SELECT 1 /* -- comment ; */ + 1;",
		"SELECT 1 -- /* comment ;\n + 1;",
	}
	fragments := psql.Split(input)

	for i, fragment := range fragments {
		require.Equal(t, expected[i], fragment)
	}
}

func TestStringRules(t *testing.T) {
	input := dedent.Dedent(`
		SELECT ';"';
		SELECT 1 ";'";
		SELECT /*'*/ 1"';";
	`)
	expected := []string{
		`SELECT ';"';`,
		`SELECT 1 ";'";`,
		`SELECT /*'*/ 1"';";`,
	}
	fragments := psql.Split(input)

	for i, fragment := range fragments {
		require.Equal(t, expected[i], fragment)
	}
}

// func TestTransactionRules(t *testing.T) {
// 	input := dedent.Dedent(`
// 		begin; SELECT 1; end;
// 		BEGIN; SELECT 1; COMMIT;
// 		begin; SELECT 1; rollback;
// 		BEGIN; SELECT 'END'; END;
// 	`)
// 	expected := []string{
// 		`begin; SELECT 1; end;`,
// 		`BEGIN; SELECT 1; COMMIT;`,
// 		`begin; SELECT 1; rollback;`,
// 		`BEGIN; SELECT 'END'; END;`,
// 	}
// 	fragments := psql.Split(input)
//
// 	for i, fragment := range fragments {
// 		require.Equal(t, expected[i], fragment)
// 	}
// }
