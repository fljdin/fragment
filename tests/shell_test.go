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
		StringRule{Start: "\\", Stop: "\n"},
		StringRule{Start: "#", StopAtDelim: true},
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
