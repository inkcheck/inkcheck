package rhetoric_test

import (
	"testing"

	"github.com/inkcheck/config"
	"github.com/inkcheck/rhetoric"
)

func TestVocabSophistication_LexicalDensity_ContentHeavy(t *testing.T) {
	// Technical prose with many content words.
	text := "Photosynthesis converts sunlight, water, and carbon dioxide into glucose and oxygen. " +
		"Chlorophyll absorbs red and blue wavelengths while reflecting green. " +
		"Mitochondria generate adenosine triphosphate through oxidative phosphorylation."
	cfg := config.DefaultConfig()
	result := rhetoric.VocabSophisticationDistribution(cfg, text)

	if result.LexicalDensity <= 0 {
		t.Errorf("expected LexicalDensity > 0, got %v", result.LexicalDensity)
	}
	if result.LexicalDensity > 1.0 {
		t.Errorf("expected LexicalDensity <= 1.0, got %v", result.LexicalDensity)
	}
	if result.LexicalDensity < 0.4 {
		t.Errorf("expected high LexicalDensity for content-heavy text, got %v", result.LexicalDensity)
	}
}

func TestVocabSophistication_LexicalDensity_FunctionWordHeavy(t *testing.T) {
	// Function-word-heavy sentence lowers lexical density.
	text := "It is what it is. This is a thing that we do because it is how it has to be. " +
		"We are all in this together and it will be as it was before. " +
		"That is all that we have and it is enough for us."
	cfg := config.DefaultConfig()
	result := rhetoric.VocabSophisticationDistribution(cfg, text)

	if result.LexicalDensity > 0.5 {
		t.Errorf("expected low LexicalDensity for function-word-heavy text, got %v", result.LexicalDensity)
	}
}

func TestVocabSophistication_LowFreqWordRatio(t *testing.T) {
	// Text with rare/unknown words should have a non-zero low-freq ratio.
	text := "The xenophobic demagogue obfuscated parliamentary deliberations with " +
		"tautological circumlocution and ephemeral sophistry. " +
		"Jurisprudential hermeneutics underpinned the bibliographic exegesis. " +
		"Epistemological presuppositions pervaded the phenomenological discourse."
	cfg := config.DefaultConfig()
	result := rhetoric.VocabSophisticationDistribution(cfg, text)

	if result.LowFreqWordRatio <= 0 {
		t.Errorf("expected LowFreqWordRatio > 0 for rare-word text, got %v", result.LowFreqWordRatio)
	}
	if result.LowFreqWordRatio > 1.0 {
		t.Errorf("expected LowFreqWordRatio <= 1.0, got %v", result.LowFreqWordRatio)
	}
}

func TestVocabSophistication_LowFreqWordRatio_CommonWords(t *testing.T) {
	// Very common words → most words in band 0 → low LowFreqWordRatio.
	text := "The big dog and the small cat and the old man and the new day and the good work."
	cfg := config.DefaultConfig()
	result := rhetoric.VocabSophisticationDistribution(cfg, text)

	if result.LowFreqWordRatio > 0.3 {
		t.Errorf("expected low LowFreqWordRatio for common-word text, got %v", result.LowFreqWordRatio)
	}
}

func TestVocabSophistication_LowFreqRatioIsBandSum(t *testing.T) {
	text := "The study of linguistics encompasses morphology, syntax, semantics, and pragmatics. " +
		"Psycholinguistics investigates how humans acquire, produce, and comprehend language. " +
		"Computational approaches leverage statistical models and neural architectures."
	cfg := config.DefaultConfig()
	result := rhetoric.VocabSophisticationDistribution(cfg, text)

	expected := result.BandRatios[2] + result.BandRatios[3] + result.BandRatios[4]
	if result.LowFreqWordRatio != expected {
		t.Errorf("LowFreqWordRatio should equal BandRatios[2]+[3]+[4]: got %v, want %v",
			result.LowFreqWordRatio, expected)
	}
}

func TestVocabSophistication_Empty(t *testing.T) {
	cfg := config.DefaultConfig()
	result := rhetoric.VocabSophisticationDistribution(cfg, "")

	if result.LexicalDensity != 0 {
		t.Errorf("expected LexicalDensity = 0 for empty text, got %v", result.LexicalDensity)
	}
	if result.LowFreqWordRatio != 0 {
		t.Errorf("expected LowFreqWordRatio = 0 for empty text, got %v", result.LowFreqWordRatio)
	}
}
