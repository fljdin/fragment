package fragment

import (
	"bufio"
	"bytes"
	"strings"
)

func trimAndAppend(slice []string, element string) []string {
	element = strings.TrimSpace(element)
	if len(element) > 0 {
		return append(slice, element)
	}
	return slice
}

type Language struct {
	Delimiters []string
	Rules      []RangeRule
}

func (lang *Language) Split(input string) []string {
	var fragment bytes.Buffer
	var fragments []string

	scanner := bufio.NewScanner(strings.NewReader(input))
	scanner.Split(bufio.ScanBytes)

	var currentRule RangeRule

Scan:
	for scanner.Scan() {
		char := scanner.Bytes()
		fragment.Write(char)

		if currentRule == nil {
			// Look for a new rule
			for _, rule := range lang.Rules {
				if rule.IsStartDetected(fragment.Bytes()) {
					currentRule = rule
					continue Scan
				}
			}
		} else {
			// Look for the end of the current rule
			if currentRule.IsEndDetected(fragment.Bytes()) {
				currentRule = nil
			}
			continue Scan
		}

		// Look for a delimiter
		for _, delimiter := range lang.Delimiters {
			if hasSuffixFold(fragment.Bytes(), toBytes(delimiter)) {
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
