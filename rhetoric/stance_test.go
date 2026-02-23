package rhetoric_test

import (
	"testing"

	"github.com/inkcheck/rhetoric"
)

func TestStanceAnalysis_ReaderCentric(t *testing.T) {
	text := "You should try this. Your experience matters to us. We built this for you. " +
		"You will find it helpful. Your feedback shapes what we do."
	result := rhetoric.StanceAnalysis(text)

	if result.SecondPerson == 0 {
		t.Error("expected SecondPerson count > 0 for you/your-heavy text")
	}
	if result.SecondPerson <= result.FirstPlural {
		t.Errorf("expected SecondPerson (%d) to dominate FirstPlural (%d)", result.SecondPerson, result.FirstPlural)
	}
	if result.ReaderCentricity <= 0 {
		t.Errorf("expected ReaderCentricity > 0, got %v", result.ReaderCentricity)
	}
	if result.TotalPronouns == 0 {
		t.Error("expected TotalPronouns > 0")
	}
}

func TestStanceAnalysis_Impersonal(t *testing.T) {
	text := "The report was submitted. Analysis has been completed. " +
		"Results are presented in the following section. Data was collected over six months."
	result := rhetoric.StanceAnalysis(text)

	if result.TotalPronouns != 0 {
		t.Errorf("expected TotalPronouns = 0 for impersonal text, got %d", result.TotalPronouns)
	}
	if result.ReaderCentricity != 0 {
		t.Errorf("expected ReaderCentricity = 0 for impersonal text, got %v", result.ReaderCentricity)
	}
}

func TestStanceAnalysis_WeAnalysis(t *testing.T) {
	text := "We believe in building great products. Our team works hard every day. " +
		"We have built something we are proud of. Our values guide us in everything we do."
	result := rhetoric.StanceAnalysis(text)

	if result.FirstPlural == 0 {
		t.Error("expected FirstPlural count > 0 for we/our-heavy text")
	}
	if result.FirstPlural <= result.SecondPerson {
		t.Errorf("expected FirstPlural (%d) to dominate SecondPerson (%d)", result.FirstPlural, result.SecondPerson)
	}
	if result.ReaderCentricity <= 0 {
		t.Errorf("expected ReaderCentricity > 0 for first-plural text, got %v", result.ReaderCentricity)
	}
}

func TestStanceAnalysis_FirstSingular(t *testing.T) {
	text := "I think my approach is correct. I have been working on my project."
	result := rhetoric.StanceAnalysis(text)

	if result.FirstSingular == 0 {
		t.Error("expected FirstSingular > 0 for I/my text")
	}
	if result.TotalPronouns == 0 {
		t.Error("expected TotalPronouns > 0")
	}
}

func TestStanceAnalysis_SingleWord(t *testing.T) {
	result := rhetoric.StanceAnalysis("Hello")

	if result.TotalPronouns != 0 {
		t.Errorf("expected TotalPronouns = 0 for single non-pronoun word, got %d", result.TotalPronouns)
	}
	if result.ReaderCentricity != 0 {
		t.Errorf("expected ReaderCentricity = 0, got %v", result.ReaderCentricity)
	}
}

func TestStanceAnalysis_Empty(t *testing.T) {
	result := rhetoric.StanceAnalysis("")

	if result.SecondPerson != 0 || result.FirstPlural != 0 ||
		result.FirstSingular != 0 || result.ThirdImpersonal != 0 {
		t.Error("expected all counts to be 0 for empty text")
	}
	if result.TotalPronouns != 0 {
		t.Errorf("expected TotalPronouns = 0, got %d", result.TotalPronouns)
	}
	if result.ReaderCentricity != 0 {
		t.Errorf("expected ReaderCentricity = 0, got %v", result.ReaderCentricity)
	}
}
