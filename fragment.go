package fragment

import (
	"bufio"
	"bytes"
	"regexp"
	"strconv"
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
	Delimiters []Delimiter
	Rules      []Rule
}

func (lang Language) Split(input string) []string {
	var fragment bytes.Buffer
	var fragments []string

	scanner := bufio.NewScanner(strings.NewReader(input))
	scanner.Split(bufio.ScanBytes)

	var currentRule Rule

Scan:
	for scanner.Scan() {
		char := scanner.Bytes()
		fragment.Write(char)

		// Look for a new rule
		if currentRule == nil {
			for _, rule := range lang.Rules {
				if rule.IsStarted(fragment.Bytes()) {
					currentRule = rule
					continue Scan
				}
			}
		}

		// Look for the end of the current rule
		if currentRule != nil {
			if currentRule.UseStopRule() {
				if currentRule.IsStopped(fragment.Bytes()) {
					currentRule = nil
				}
				continue Scan
			}
		}

		// Look for a delimiter
		for _, delimiter := range lang.Delimiters {
			if delimiter.IsDetected(fragment.Bytes()) {
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

type Delimiter struct {
	String string
	Regex  string
}

func (d Delimiter) IsDetected(input []byte) bool {
	if len(d.String) > 0 {
		return bytes.HasSuffix(input, []byte(d.String))
	}
	if len(d.Regex) > 0 {
		re := regexp.MustCompile(d.Regex + `$`)
		return re.Match(input)
	}
	return false
}

type Rule interface {
	IsStarted(input []byte) bool
	IsStopped(input []byte) bool
	UseStopRule() bool
}

type StringRule struct {
	Start       string
	Stop        string
	StopAtDelim bool
}

func (s StringRule) IsStarted(input []byte) bool {
	return bytes.HasSuffix(input, []byte(s.Start))
}

func (s StringRule) IsStopped(input []byte) bool {
	return bytes.HasSuffix(input, []byte(s.Stop))
}

func (s StringRule) UseStopRule() bool {
	return !s.StopAtDelim
}

type RegexRule struct {
	Start   string
	Stop    string
	matches [][]byte
}

func (r *RegexRule) IsStarted(input []byte) bool {
	// append $ to match only the end of the input
	re := regexp.MustCompile(r.Start + `$`)
	r.matches = re.FindSubmatch(input)

	return len(r.matches) > 0
}

func (r RegexRule) IsStopped(input []byte) bool {
	// append $ to match only the end of the input
	re := regexp.MustCompile(r.groupAsRegex() + `$`)
	return re.Match(input)
}

func (r RegexRule) groupAsRegex() string {
	re := regexp.MustCompile(`\\(\d)`)
	result := re.ReplaceAllFunc([]byte(r.Stop), func(match []byte) []byte {
		idxStr := match[1:]
		idx, err := strconv.Atoi(string(idxStr))
		// replace unknown index by empty string
		if err != nil || idx < 1 || idx > len(r.matches) {
			return []byte{}
		}
		return r.matches[idx]
	})
	return string(result)
}

func (r RegexRule) UseStopRule() bool {
	return true
}
