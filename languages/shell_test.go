package languages_test

import (
	"testing"

	"github.com/fljdin/fragment"
	. "github.com/fljdin/fragment/languages"
	"github.com/lithammer/dedent"
	"github.com/stretchr/testify/require"
)

func init() {
	fragment.TrimOption = true
}

func TestNewlineDelimiter(t *testing.T) {
	fragments := Shell.Split("true\nfalse")

	require.Equal(t, "true", fragments[0])
	require.Equal(t, "false", fragments[1])
}

func TestShellSemicolonDelimiter(t *testing.T) {
	fragments := Shell.Split("true; false")

	require.Equal(t, "true;", fragments[0])
	require.Equal(t, "false", fragments[1])
}

func TestIgnoreEmptyLines(t *testing.T) {
	input := dedent.Dedent(`
		true

		false
	`)
	fragments := Shell.Split(input)

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
	fragments := Shell.Split(input)

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
	fragments := Shell.Split(input)

	require.Equal(t, "true \\\n  && false", fragments[0])
	require.Equal(t, "false", fragments[1])
}

func TestCommentRule(t *testing.T) {
	input := dedent.Dedent(`
		true # comment \
		false
	`)
	fragments := Shell.Split(input)

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
	fragments := Shell.Split(input)

	for i, fragment := range fragments {
		require.Equal(t, expected[i], fragment)
	}
}
