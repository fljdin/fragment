package languages

import "github.com/fljdin/fragment"

var XML = fragment.Language{
	Delimiters: []fragment.Delimiter{
		// Delimiter for line endings, indicating the end of a line
		{String: "\n"},
	},
	Rules: []fragment.Rule{
		// XML tags in XML documents, both opening and closing tags This
		// configuration is case-insensitive and dynamically matches the tag
		// names
		&fragment.RegexRule{Start: `(?i)<(\w+)>`, Stop: `(?i)</\1>`},
	},
}
