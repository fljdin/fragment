package languages

import "github.com/fljdin/fragment"

var Shell = fragment.Language{
	Delimiters: []fragment.Delimiter{
		// Delimiter for inlined commands
		{String: ";"},

		// Delimiter for line endings, indicating the end of a line
		{String: "\n"},
	},
	Rules: []fragment.Rule{
		// Single-quoted string literals in Shell, starting and ending with
		// single quotes
		fragment.StringRule{Start: `'`, Stop: `'`},

		// Double-quoted string literals in Shell, starting and ending with
		// double quotes
		fragment.StringRule{Start: `"`, Stop: `"`},

		// Handling of line continuation with backslashes ('\') in Shell
		fragment.StringRule{Start: "\\", Stop: "\n"},

		// Shell-style comments starting with '#' and stopping at the end of the
		// line (including the delimiter)
		fragment.StringRule{Start: "#", StopAtDelim: true},

		// Handling of Here-Documents (Heredocs) in Shell, which start with
		// '<<-' or '<<' followed by a delimiter. The Stop rule dynamically
		// matches the specified delimiter, preserving leading whitespace
		&fragment.RegexRule{
			Start: `<<-?\s*"?'?([^"'<>\s\n]+).*\n`,
			Stop:  `\n\1`,
		},
	},
}
