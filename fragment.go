package fragment

import (
	"bufio"
	"bytes"
	"io"
	"regexp"
	"strconv"
	"strings"
)

var TrimOption bool = false

type Language struct {
	Delimiters []Delimiter
	Rules      []Rule
}

// Split() splits the input string into fragments.
func (lang *Language) Split(input string) (fragments []string) {
	ch := make(chan string)
	go lang.Read(ch, strings.NewReader(input))

	for fragment := range ch {
		if len(fragment) == 0 {
			continue
		}
		fragments = append(fragments, fragment)
	}
	return
}

// Read() reads the input stream and pushes to the fragments channel.
func (lang *Language) Read(ch chan string, input io.Reader) {
	defer close(ch)

	var fragment bytes.Buffer

	scanner := bufio.NewScanner(input)
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
				ch <- trim(fragment.String())
				fragment.Reset()
				break
			}
		}
	}

	ch <- trim(fragment.String())
}

func trim(s string) string {
	if TrimOption {
		return strings.TrimSpace(s)
	}
	return s
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
