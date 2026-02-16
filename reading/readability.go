package reading

import (
	"fmt"

	"github.com/inkcheck/config"
	readabilitylib "github.com/inkcheck/readability"
	"github.com/inkcheck/shared"
)

type ReadabilityResult struct {
	Formula string
	Score   float64
	Grade   string // For display (e.g., "8th-9th grade" or "College level")
}

// ReadabilityAnalysis computes the configured readability formula.
func ReadabilityAnalysis(cfg config.Config, text string) ReadabilityResult {
	prose := shared.ExtractProseText(text)
	if prose == "" {
		return ReadabilityResult{}
	}

	analysis := readabilitylib.NewAnalysis(prose)
	formula := readabilitylib.Formula(cfg.ReadabilityFormula)

	// Validate formula, fallback to default if invalid
	if !formula.Valid() {
		formula = readabilitylib.FleschKincaidGrade
	}

	score, _ := analysis.Score(formula)

	return ReadabilityResult{
		Formula: string(formula),
		Score:   score,
		Grade:   formatGrade(formula, score),
	}
}

// ReadabilityAnalysisWithFormula allows specifying formula at runtime.
func ReadabilityAnalysisWithFormula(text string, formula string) ReadabilityResult {
	prose := shared.ExtractProseText(text)
	if prose == "" {
		return ReadabilityResult{}
	}

	analysis := readabilitylib.NewAnalysis(prose)
	f := readabilitylib.Formula(formula)

	if !f.Valid() {
		f = readabilitylib.FleschKincaidGrade
	}

	score, _ := analysis.Score(f)

	return ReadabilityResult{
		Formula: string(f),
		Score:   score,
		Grade:   formatGrade(f, score),
	}
}

// formatGrade converts numeric score to human-readable grade level.
func formatGrade(formula readabilitylib.Formula, score float64) string {
	// Grade-level formulas (flesch_kincaid_grade, gunning_fog, etc.)
	gradeFormulas := map[readabilitylib.Formula]bool{
		readabilitylib.FleschKincaidGrade:          true,
		readabilitylib.GunningFog:                  true,
		readabilitylib.SmogIndex:                   true,
		readabilitylib.ColemanLiauIndex:            true,
		readabilitylib.AutomatedReadabilityIndex:   true,
		readabilitylib.LinsearWriteFormula:         true,
		readabilitylib.SpacheReadability:           true,
		readabilitylib.TextStandard:                true,
	}

	if gradeFormulas[formula] {
		grade := int(score)
		switch {
		case grade < 1:
			return "Kindergarten"
		case grade <= 12:
			return fmt.Sprintf("%dth grade", grade)
		case grade <= 16:
			return "College level"
		default:
			return "Graduate level"
		}
	}

	// flesch_reading_ease uses different scale (0-100, higher = easier)
	if formula == readabilitylib.FleschReadingEase {
		switch {
		case score >= 90:
			return "Very easy (5th grade)"
		case score >= 80:
			return "Easy (6th grade)"
		case score >= 70:
			return "Fairly easy (7th grade)"
		case score >= 60:
			return "Standard (8th-9th grade)"
		case score >= 50:
			return "Fairly difficult (10th-12th grade)"
		case score >= 30:
			return "Difficult (College)"
		default:
			return "Very difficult (Graduate)"
		}
	}

	// Other formulas (lix, rix, mcalpine_eflaw, etc.) - just return score
	return ""
}
