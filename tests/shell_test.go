package tests

import (
	"testing"

	. "github.com/fljdin/fragment"
	"github.com/lithammer/dedent"
	"github.com/stretchr/testify/require"
)

var shell = Language{
	Delimiters: []string{"\n"},
	Rules: []Rule{
		StringRule{Start: `'`, Stop: `'`},
		StringRule{Start: `"`, Stop: `"`},
		StringRule{Start: "\\", Stop: "\n"},
		StringRule{Start: "#", StopAtDelim: true},
		&RegexRule{
			Start: `<<-?\s*"?'?([^"'<>\s\n]+).*\n`,
			Stop:  `\n\1`,
		},
	},
}

func TestNewlineDelimiter(t *testing.T) {
	fragments := shell.Split("true\nfalse")

	require.Equal(t, "true", fragments[0])
	require.Equal(t, "false", fragments[1])
}

func TestIgnoreEmptyLines(t *testing.T) {
	input := dedent.Dedent(`
		true

		false
	`)
	fragments := shell.Split(input)

	require.Equal(t, "true", fragments[0])
	require.Equal(t, "false", fragments[1])
}

func TestMultilineStringRule(t *testing.T) {
	input := dedent.Dedent(`
		echo "hello '
		world"
		echo 'hello "
		world'
	`)
	expected := []string{
		"echo \"hello '\nworld\"",
		"echo 'hello \"\nworld'",
	}
	fragments := shell.Split(input)

	for i, fragment := range fragments {
		require.Equal(t, expected[i], fragment)
	}
}

func TestNewlineEscapeRule(t *testing.T) {
	input := dedent.Dedent(`
		true \
		  && false
		false
	`)
	fragments := shell.Split(input)

	require.Equal(t, "true \\\n  && false", fragments[0])
	require.Equal(t, "false", fragments[1])
}

func TestCommentRule(t *testing.T) {
	input := dedent.Dedent(`
		true # comment \
		false
	`)
	fragments := shell.Split(input)

	require.Equal(t, "true # comment \\", fragments[0])
	require.Equal(t, "false", fragments[1])
}

func TestHereDocRule(t *testing.T) {
	input := dedent.Dedent(`
		<<HERE
		hello world
		HERE
		cat <<- "EOF"
		hello world
		EOF
		cat << 'content' > content.txt
		hello world
		content
	`)
	expected := []string{
		"<<HERE\nhello world\nHERE",
		"cat <<- \"EOF\"\nhello world\nEOF",
		"cat << 'content' > content.txt\nhello world\ncontent",
	}
	fragments := shell.Split(input)

	for i, fragment := range fragments {
		require.Equal(t, expected[i], fragment)
	}
}
