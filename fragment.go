package fragment

import (
	"bytes"
)

func hasSuffix(input, suffix []byte) bool {
	if len(input) < len(suffix) {
		return false
	}
	return bytes.HasSuffix(input, suffix)
}

func toBytes(input string) []byte {
	if input == "\\n" {
		return []byte("\x0a")
	}
	return []byte(input)
}
