package structure

import (
	"strings"

	"github.com/inkcheck/shared"
)

// SentenceTypeResult holds the distribution and entropy of sentence types.
type SentenceTypeResult struct {
	Declarative   int
	Interrogative int
	Imperative    int
	Exclamatory   int
	Total         int
	Entropy       float64 // Shannon entropy of the 4-type distribution (bits); max = log2(4) ≈ 2.0
}

// imperativeVerbs is a set of common English verbs used to open imperative sentences.
var imperativeVerbs = map[string]bool{
	"consider": true, "think": true, "imagine": true, "note": true, "remember": true,
	"see": true, "look": true, "try": true, "make": true, "get": true, "take": true,
	"use": true, "find": true, "go": true, "come": true, "keep": true, "let": true,
	"give": true, "ask": true, "tell": true, "help": true, "show": true, "check": true,
	"read": true, "click": true, "download": true, "sign": true, "learn": true,
	"discover": true, "explore": true, "start": true, "stop": true, "avoid": true,
	"ensure": true, "follow": true, "choose": true, "create": true, "build": true,
	"add": true, "set": true, "run": true, "install": true, "write": true,
	"open": true, "close": true, "move": true, "change": true, "update": true,
	"delete": true, "save": true, "share": true, "send": true, "buy": true,
	"order": true, "contact": true, "call": true, "visit": true, "press": true,
	"select": true, "enter": true, "type": true, "scroll": true, "watch": true,
	"listen": true, "speak": true, "focus": true, "assume": true, "suppose": true,
	"notice": true, "observe": true, "compare": true, "define": true, "describe": true,
	"explain": true, "identify": true, "include": true, "provide": true, "refer": true,
	"review": true, "specify": true, "state": true, "summarise": true, "summarize": true,
	"understand": true, "verify": true, "confirm": true, "enable": true, "disable": true,
	"configure": true, "deploy": true, "push": true, "pull": true, "merge": true,
	"test": true, "debug": true, "fix": true, "validate": true, "export": true,
	"import": true, "repeat": true, "replace": true, "combine": true, "connect": true,
	"disconnect": true, "reset": true, "restart": true, "launch": true,
	"be": true,
}

// SentenceTypeDistribution classifies sentences in text by type and computes
// the Shannon entropy of the resulting distribution.
func SentenceTypeDistribution(text string) SentenceTypeResult {
	prose := shared.ExtractProseText(text)
	sentences := shared.SplitSentences(prose)

	if len(sentences) == 0 {
		return SentenceTypeResult{}
	}

	var declarative, interrogative, imperative, exclamatory int

	for _, s := range sentences {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		switch {
		case strings.HasSuffix(s, "?"):
			interrogative++
		case strings.HasSuffix(s, "!"):
			exclamatory++
		default:
			// Check for imperative: first word is a known imperative verb
			firstWord := firstToken(s)
			if imperativeVerbs[strings.ToLower(firstWord)] {
				imperative++
			} else {
				declarative++
			}
		}
	}

	total := declarative + interrogative + imperative + exclamatory
	if total == 0 {
		return SentenceTypeResult{}
	}

	counts := [4]int{declarative, interrogative, imperative, exclamatory}
	dist := make([]float64, len(counts))
	for i, c := range counts {
		dist[i] = float64(c) / float64(total)
	}
	entropy := shared.Entropy(dist)

	return SentenceTypeResult{
		Declarative:   declarative,
		Interrogative: interrogative,
		Imperative:    imperative,
		Exclamatory:   exclamatory,
		Total:         total,
		Entropy:       entropy,
	}
}

// firstToken returns the first whitespace-separated token from s,
// with leading/trailing punctuation stripped.
func firstToken(s string) string {
	fields := strings.Fields(s)
	if len(fields) == 0 {
		return ""
	}
	return shared.StripPunctuation(fields[0])
}
