package structure

import "github.com/inkcheck/shared"

type PunctuationProfile struct {
	Periods      int
	Commas       int
	Semicolons   int
	Colons       int
	Dashes       int
	Parentheses  int
	Questions    int
	Exclamations int
}

func (p PunctuationProfile) Variety() int {
	count := 0
	if p.Periods > 0 {
		count++
	}
	if p.Commas > 0 {
		count++
	}
	if p.Semicolons > 0 {
		count++
	}
	if p.Colons > 0 {
		count++
	}
	if p.Dashes > 0 {
		count++
	}
	if p.Parentheses > 0 {
		count++
	}
	if p.Questions > 0 {
		count++
	}
	if p.Exclamations > 0 {
		count++
	}
	return count
}

func (p PunctuationProfile) Total() int {
	return p.Periods + p.Commas + p.Semicolons + p.Colons +
		p.Dashes + p.Parentheses + p.Questions + p.Exclamations
}

func PunctuationAnalysis(text string) PunctuationProfile {
	prose := shared.ExtractProseText(text)
	var p PunctuationProfile
	var prev rune
	for _, r := range prose {
		switch r {
		case '.':
			p.Periods++
		case ',':
			p.Commas++
		case ';':
			p.Semicolons++
		case ':':
			p.Colons++
		case '?':
			p.Questions++
		case '!':
			p.Exclamations++
		case '(':
			p.Parentheses++
		case '-':
			if prev == '-' {
				p.Dashes++
			}
		case '\u2013', '\u2014':
			p.Dashes++
		}
		prev = r
	}
	return p
}
