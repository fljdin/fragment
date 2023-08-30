package fragment

import (
	"bytes"
)

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
