package fragment

import (
	"regexp"
	"strconv"
)

type Rule interface {
	IsStartDetected(input []byte) bool
	IsEndDetected(input []byte) bool
}

type StringRule struct {
	Start string
	End   string
}

func (s StringRule) IsStartDetected(input []byte) bool {
	return hasSuffixFold(input, toBytes(s.Start))
}

func (s StringRule) IsEndDetected(input []byte) bool {
	return hasSuffixFold(input, toBytes(s.End))
}

type RegexRule struct {
	Start   string
	End     string
	matches [][]byte
}

func (r *RegexRule) IsStartDetected(input []byte) bool {
	// append $ to match only the end of the input
	re := regexp.MustCompile(r.Start + `$`)
	r.matches = re.FindSubmatch(input)

	return len(r.matches) > 0
}

func (r RegexRule) IsEndDetected(input []byte) bool {
	// append $ to match only the end of the input
	re := regexp.MustCompile(r.groupAsRegex() + `$`)
	return re.Match(input)
}

func (r RegexRule) groupAsRegex() string {
	re := regexp.MustCompile(`\\(\d)`)
	result := re.ReplaceAllFunc([]byte(r.End), func(match []byte) []byte {
		idxStr := match[1:]
		idx, err := strconv.Atoi(string(idxStr))
		// replace unkown index by empty string
		if err != nil || idx < 1 || idx > len(r.matches) {
			return []byte{}
		}
		return r.matches[idx]
	})
	return string(result)
}
