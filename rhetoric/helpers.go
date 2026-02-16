package rhetoric

// countOccurrences counts word-boundary-aware occurrences of phrase in text.
// Both text and phrase should be lowercased before calling.
func countOccurrences(text, phrase string) int {
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
