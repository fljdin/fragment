package fragment

import (
	"bufio"
	"bytes"
	"strings"
)

type RangeRule struct {
	Start string
	End   string
}

type Language struct {
	Delimiters []string
	Rules      []RangeRule
}

func (lang *Language) Split(input string) []string {
	var fragments []string
	var fragment bytes.Buffer

	scanner := bufio.NewScanner(strings.NewReader(input))
	scanner.Split(bufio.ScanBytes)

	for scanner.Scan() {
		char := scanner.Bytes()
		fragment.Write(char)

		for _, delimiter := range lang.Delimiters {
			if bytes.HasSuffix(fragment.Bytes(), []byte(delimiter)) {
				fragments = append(fragments, strings.TrimSpace(fragment.String()))
				fragment.Reset()
				break
			}
		}
	}

	if fragment.Len() > 0 {
		fragments = append(fragments, strings.TrimSpace(fragment.String()))
	}

	return fragments
}
