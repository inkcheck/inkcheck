package structure_test

import (
	"testing"

	"github.com/inkcheck/structure"
)

func TestSentenceTypeDistribution_Mixed(t *testing.T) {
	text := "The sky is blue. Is that always true? Consider the clouds. What a beautiful day!"
	result := structure.SentenceTypeDistribution(text)

	if result.Total == 0 {
		t.Fatal("expected Total > 0")
	}
	if result.Interrogative == 0 {
		t.Error("expected at least one interrogative sentence (ending with ?)")
	}
	if result.Exclamatory == 0 {
		t.Error("expected at least one exclamatory sentence (ending with !)")
	}
	if result.Imperative == 0 {
		t.Error("expected at least one imperative sentence (starting with 'consider')")
	}
	if result.Entropy <= 0 {
		t.Errorf("expected Entropy > 0 for mixed sentence types, got %v", result.Entropy)
	}
	if result.Entropy > 2.01 {
		t.Errorf("Entropy should be <= log2(4)=2.0, got %v", result.Entropy)
	}
}

func TestSentenceTypeDistribution_AllDeclarative(t *testing.T) {
	text := "The project was completed on time. The team worked hard. Results were excellent."
	result := structure.SentenceTypeDistribution(text)

	if result.Total == 0 {
		t.Fatal("expected Total > 0")
	}
	if result.Declarative != result.Total {
		t.Errorf("expected all sentences to be declarative: Declarative=%d, Total=%d",
			result.Declarative, result.Total)
	}
	if result.Entropy != 0 {
		t.Errorf("expected Entropy = 0 for all-declarative text, got %v", result.Entropy)
	}
}

func TestSentenceTypeDistribution_AmbiguousFirstWord(t *testing.T) {
	// "Set" is in the imperative verb list; tests word-level heuristic with ambiguity.
	text := "Set the value to zero. The set contains three elements."
	result := structure.SentenceTypeDistribution(text)

	if result.Total != 2 {
		t.Errorf("expected 2 sentences, got %d", result.Total)
	}
	// "Set" at the start triggers imperative detection
	if result.Imperative < 1 {
		t.Error("expected at least 1 imperative sentence starting with 'Set'")
	}
	if result.Declarative < 1 {
		t.Error("expected at least 1 declarative sentence starting with 'The'")
	}
}

func TestSentenceTypeDistribution_SingleSentence(t *testing.T) {
	text := "The cat sat on the mat."
	result := structure.SentenceTypeDistribution(text)

	if result.Total != 1 {
		t.Errorf("expected Total = 1, got %d", result.Total)
	}
	if result.Declarative != 1 {
		t.Errorf("expected Declarative = 1, got %d", result.Declarative)
	}
	if result.Entropy != 0 {
		t.Errorf("expected Entropy = 0 for single sentence, got %v", result.Entropy)
	}
}

func TestSentenceTypeDistribution_ExclamatoryImperative(t *testing.T) {
	// "Go!" ends with ! so it should be classified as exclamatory, not imperative.
	text := "Go!"
	result := structure.SentenceTypeDistribution(text)

	if result.Total != 1 {
		t.Errorf("expected Total = 1, got %d", result.Total)
	}
	if result.Exclamatory != 1 {
		t.Errorf("expected Exclamatory = 1 for 'Go!', got %d", result.Exclamatory)
	}
}

func TestSentenceTypeDistribution_Empty(t *testing.T) {
	result := structure.SentenceTypeDistribution("")

	if result.Total != 0 {
		t.Errorf("expected Total = 0, got %d", result.Total)
	}
	if result.Entropy != 0 {
		t.Errorf("expected Entropy = 0, got %v", result.Entropy)
	}
}
