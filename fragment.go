package fragment

import (
	"bufio"
	"bytes"
	"strings"
)

type RangeRule struct {
	Start any
	End   any
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

func hasSuffixFold(input, suffix []byte) bool {
	if len(input) < len(suffix) {
		return false
	}
	return bytes.EqualFold(input[len(input)-len(suffix):], suffix)
}

func toBytes(input string) []byte {
	if input == "\\n" {
		return []byte("\x0a")
	}
	return []byte(input)
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
				if rule.IsStartDetected(fragment.Bytes()) {
					currentRule = &rule
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

func (r *RangeRule) IsStartDetected(input []byte) bool {
	switch start := r.Start.(type) {
	case string:
		if hasSuffixFold(input, toBytes(start)) {
			return true
		}
	case []string:
		for _, start := range start {
			if hasSuffixFold(input, toBytes(start)) {
				return true
			}
		}
	}
	return false
}

func (r *RangeRule) IsEndDetected(input []byte) bool {
	switch end := r.End.(type) {
	case string:
		if hasSuffixFold(input, toBytes(end)) {
			return true
		}
	case []string:
		for _, end := range end {
			if hasSuffixFold(input, toBytes(end)) {
				return true
			}
		}
	}
	return false
}
