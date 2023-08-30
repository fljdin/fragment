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

func trimAndAppend(slice []string, element string) []string {
	element = strings.TrimSpace(element)
	if len(element) > 0 {
		return append(slice, element)
	}
	return slice
}

func hasSuffixFold(s, suffix []byte) bool {
	if len(s) < len(suffix) {
		return false
	}
	return bytes.EqualFold(s[len(s)-len(suffix):], suffix)
}

func (lang *Language) Split(input string) []string {
	var fragment bytes.Buffer
	var fragments []string

	scanner := bufio.NewScanner(strings.NewReader(input))
	scanner.Split(bufio.ScanBytes)

	var currentRule *RangeRule

Scan:
	for scanner.Scan() {
		char := scanner.Bytes()
		fragment.Write(char)

		if currentRule == nil {
			// Look for a new rule
			for _, rule := range lang.Rules {
				if hasSuffixFold(fragment.Bytes(), []byte(rule.Start)) {
					currentRule = &rule
					continue Scan
				}
			}
		} else {
			// Look for the end of the current rule
			var stopBytes []byte
			if currentRule.End == "\\n" {
				stopBytes = []byte("\x0a")
			} else {
				stopBytes = []byte(currentRule.End)
			}

			if hasSuffixFold(fragment.Bytes(), stopBytes) {
				currentRule = nil
			}
			continue Scan
		}

		// Look for a delimiter
		for _, delimiter := range lang.Delimiters {
			if hasSuffixFold(fragment.Bytes(), []byte(delimiter)) {
				fragments = trimAndAppend(fragments, fragment.String())
				fragment.Reset()
				break
			}
		}
	}

	if fragment.Len() > 0 {
		fragments = trimAndAppend(fragments, fragment.String())
	}

	return fragments
}
