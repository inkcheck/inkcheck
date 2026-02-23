package rhetoric

import (
	"strings"

	"github.com/inkcheck/shared"
)

// StanceResult holds the pronoun-based stance analysis of a text.
type StanceResult struct {
	SecondPerson     int     // count of 2nd-person pronouns (you, your, ...)
	FirstPlural      int     // count of 1st-person plural pronouns (we, our, ...)
	FirstSingular    int     // count of 1st-person singular pronouns (I, my, ...)
	ThirdImpersonal  int     // count of 3rd-person / impersonal pronouns (one, they, ...)
	TotalPronouns    int     // sum of all above
	ReaderCentricity float64 // weighted score / total_words; 0 if no pronouns or no words
}

var secondPersonPronouns = map[string]bool{
	"you": true, "your": true, "yours": true, "yourself": true, "yourselves": true,
}

var firstPluralPronouns = map[string]bool{
	"we": true, "our": true, "ours": true, "ourselves": true, "us": true,
}

var firstSingularPronouns = map[string]bool{
	"i": true, "my": true, "mine": true, "myself": true, "me": true,
}

var thirdImpersonalPronouns = map[string]bool{
	"one": true, "they": true, "their": true, "theirs": true,
	"them": true, "themselves": true,
}

// StanceAnalysis analyses the pronoun stance of the given text, computing
// counts for each pronoun category and a reader-centricity score.
// Score formula: (you×1.0 + we×0.6 + I×0.4 + they/one×0.1) / total_words.
// Zero pronouns → ReaderCentricity of 0 (impersonal/institutional stance).
func StanceAnalysis(text string) StanceResult {
	prose := shared.ExtractProseText(text)
	words := shared.ListWords(prose)
	totalWords := len(words)

	if totalWords == 0 {
		return StanceResult{}
	}

	var secondPerson, firstPlural, firstSingular, thirdImpersonal int

	for _, w := range words {
		token := strings.ToLower(shared.StripPunctuation(w))
		if token == "" {
			continue
		}
		switch {
		case secondPersonPronouns[token]:
			secondPerson++
		case firstPluralPronouns[token]:
			firstPlural++
		case firstSingularPronouns[token]:
			firstSingular++
		case thirdImpersonalPronouns[token]:
			thirdImpersonal++
		}
	}

	totalPronouns := secondPerson + firstPlural + firstSingular + thirdImpersonal

	readerCentricity := 0.0
	if totalPronouns > 0 {
		weighted := float64(secondPerson)*1.0 +
			float64(firstPlural)*0.6 +
			float64(firstSingular)*0.4 +
			float64(thirdImpersonal)*0.1
		readerCentricity = weighted / float64(totalWords)
	}

	return StanceResult{
		SecondPerson:     secondPerson,
		FirstPlural:      firstPlural,
		FirstSingular:    firstSingular,
		ThirdImpersonal:  thirdImpersonal,
		TotalPronouns:    totalPronouns,
		ReaderCentricity: readerCentricity,
	}
}
