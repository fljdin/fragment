package fragment

type RangeRule interface {
	IsStartDetected(input []byte) bool
	IsEndDetected(input []byte) bool
}

type StringRule struct {
	Start any
	End   any
}

func (s StringRule) IsStartDetected(input []byte) bool {
	switch start := s.Start.(type) {
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

func (s StringRule) IsEndDetected(input []byte) bool {
	switch end := s.End.(type) {
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

type RegexRule struct {
	Start string
	End   string
	group []string
}

func (r RegexRule) IsStartDetected(input []byte) bool {
	return false
}

func (r RegexRule) IsEndDetected(input []byte) bool {
	return false
}
