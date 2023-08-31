package languages

import "github.com/fljdin/fragment"

var PgSQL = fragment.Language{
	Delimiters: []fragment.Delimiter{
		// Delimiter for SQL statements, indicating the end of a statement
		{String: ";"},

		// Delimiter for special PostgreSQL commands (\g, \gdesc, \gexec, \gx,
		// crosstabview) followed by any characters and a newline
		{Regex: `\\(g|gdesc|gexec|gx|crosstabview).*\n`},
	},
	Rules: []fragment.Rule{
		// Single-line comments in PostgreSQL starting with '--' and ending at
		// the newline
		fragment.StringRule{Start: "--", Stop: "\n"},

		// Multi-line comments in PostgreSQL enclosed within '/*' and '*/'
		fragment.StringRule{Start: "/*", Stop: "*/"},

		// Single-quoted string literals
		fragment.StringRule{Start: "'", Stop: "'"},

		// Double-quoted string literals
		fragment.StringRule{Start: `"`, Stop: `"`},

		// Transaction block rules for PostgreSQL (case-insensitive)
		&fragment.RegexRule{Start: `(?i)BEGIN`, Stop: `(?i)END|COMMIT|ROLLBACK`},

		// Block rules for PostgreSQL dollar-quoted strings (e.g., $tag$ ...
		// $tag$)
		&fragment.RegexRule{Start: `(\$([a-zA-Z0-9_]*)\$)`, Stop: `\$\2\$`},
	},
}
