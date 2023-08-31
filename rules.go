package fragment

import (
	"regexp"
	"strconv"
)

type Rule interface {
	IsStarted(input []byte) bool
	IsStopped(input []byte) bool
}

type StringRule struct {
	Start string
	Stop  string
}

func (s StringRule) IsStarted(input []byte) bool {
	return hasSuffix(input, toBytes(s.Start))
}

func (s StringRule) IsStopped(input []byte) bool {
	return hasSuffix(input, toBytes(s.Stop))
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
