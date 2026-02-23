package shared

import "strings"

// CountOccurrences counts word-boundary-aware occurrences of phrase in text.
// Both text and phrase should be lowercased before calling.
// Word boundaries are defined by ASCII letters, digits, and underscores.
func CountOccurrences(text, phrase string) int {
	count := 0
	tLen := len(text)
	pLen := len(phrase)
	for i := 0; i <= tLen-pLen; i++ {
		if text[i:i+pLen] == phrase {
			if i > 0 && isWordChar(text[i-1]) {
				continue
			}
			end := i + pLen
			if end < tLen && isWordChar(text[end]) {
				continue
			}
			count++
			i += pLen - 1
		}
	}
	return count
}

func isWordChar(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9') || b == '_'
}

// ContainsAny reports whether text contains any of the given phrases.
func ContainsAny(text string, phrases []string) bool {
	for _, p := range phrases {
		if strings.Contains(text, p) {
			return true
		}
	}
	return false
}
