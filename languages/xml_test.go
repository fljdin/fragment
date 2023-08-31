package languages_test

import (
	"testing"

	. "github.com/fljdin/fragment/languages"
	"github.com/lithammer/dedent"
	"github.com/stretchr/testify/require"
)

func TestMixedCaseMarkup(t *testing.T) {
	input := dedent.Dedent(`
		<HTML>
		  <BODY>
		    <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed id maximus
		    augue, ut tincidunt elit. Vivamus leo est, finibus egestas lobortis non,
		    interdum a ipsum.</p>
		  </body>
		</html>
		<NOTE>
		  <TO>Tove</TO>
		  <FROM>Jani</FROM>
		  <HEADING>Reminder</HEADING>
		  <BODY>Don't forget me this weekend!</BODY>
		</NOTE>
	`)

	fragments := XML.Split(input)
	require.Equal(t, 2, len(fragments))
}
