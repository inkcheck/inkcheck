# Inkcheck

A Go library and CLI tool for analyzing text through structural, rhetorical,
and semantic analysis. It computes 27 metrics that measure writing variability,
argument structure, and topic coherence.

## Features

- **27 metrics** across three categories: structural, rhetorical, and semantic
- **Heuristic-first** approach using POS tagging, keyword matching, and
  statistical variance — no ML models required for 22 of 27 metrics
- **Word2vec semantic analysis** for topic coherence, redundancy,
  information novelty, and emotional tone (5 metrics)
- **Library and CLI** for use in Go programs or from the command line
- **Markdown-aware** text extraction strips headings, code, lists, and other
  non-prose elements before analysis

## Installation

### Library

```bash
go get github.com/inkcheck
```

### CLI

Using Homebrew:

```bash
brew install inkcheck/tap/inkcheck
```

Using Go:

```bash
go install github.com/inkcheck/cmd/inkcheck@latest
```

Or build from source:

```bash
go build -o inkcheck ./cmd/inkcheck
```

## Quick Start

### As a Library

Run all structure and rhetoric metrics at once:

```go
package main

import (
    "fmt"
    "github.com/inkcheck"
    "github.com/inkcheck/config"
)

func main() {
    cfg := config.DefaultConfig()
    text := `# My Essay

    First paragraph with a short opener.

    Second paragraph is longer and more detailed, covering multiple aspects
    of the topic with varied sentence lengths and structure.

    Third paragraph wraps up.`

    // Run structure + rhetoric metrics (no model needed)
    result := inkcheck.Analyze(cfg, text)
    fmt.Printf("Paragraph variance: %.4f\n", result.Structure.ParagraphVariance)
    fmt.Printf("Sentence length CV: %.4f\n", result.Structure.SentenceLengthVariance)
    fmt.Printf("Hedging density: %.2f/100w\n", result.Rhetoric.Hedging.Density)
    fmt.Printf("Formal words: %d\n", result.Rhetoric.VocabSophistication.FormalWordCount)
}
```

Run all metrics including semantic analysis:

```go
import (
    "github.com/inkcheck"
    "github.com/inkcheck/config"
    "github.com/inkcheck/semantic"
)

cfg := config.DefaultConfig()
model, err := semantic.LoadModel(cfg)
if err != nil {
    log.Fatal(err)
}
result := inkcheck.AnalyzeAll(cfg, text, model)
fmt.Printf("Topic coherence: %.4f\n", result.Semantic.TopicCoherence.MeanSimilarity)
```

Or use sub-packages directly for individual metrics:

```go
import (
    "github.com/inkcheck/config"
    "github.com/inkcheck/structure"
    "github.com/inkcheck/rhetoric"
)

cfg := config.DefaultConfig()
cv := structure.ParagraphVariance(text)
h := rhetoric.HedgingAnalysis(text)
t := rhetoric.TransitionWordDensity(cfg, text)
```

### As a CLI Tool

```bash
# Run all metrics on a file
inkcheck essay.md

# Run a single metric
inkcheck -m paragraph_variance essay.md

# Score all text files in a directory
inkcheck -m sentence_length_variance ./documents/

# Read from stdin
cat essay.md | inkcheck -m hedging_analysis

# List available metrics
inkcheck -m help

# Initialize config file
inkcheck config init

# Show current configuration
inkcheck config list

# Download semantic model
inkcheck model download
```

## Available Metrics

### Structure (6 metrics)

| Metric | Output | Description |
|--------|--------|-------------|
| `paragraph_variance` | CV (float) | Coefficient of variation of paragraph word counts |
| `sentence_length_variance` | CV (float) | CV of sentence word counts |
| `sentence_opener_diversity` | Structured | Unique opener ratio and Shannon entropy of opener distribution |
| `sentence_type_distribution` | Structured | Counts of declarative/interrogative/imperative/exclamatory sentences and entropy |
| `paragraph_position_analysis` | Structured | Opening/closing paragraph deviation from body mean |
| `punctuation_profile` | Structured | Distribution of 8 punctuation types |

### Rhetoric (16 metrics)

| Metric | Output | Description |
|--------|--------|-------------|
| `transition_word_density` | Structured | Variety and repetition of transition phrases |
| `vocabulary_sophistication` | Structured | TTR, MATTR, frequency bands, formal word count |
| `hedging_analysis` | Structured | Hedging density and category breakdown |
| `specificity_score` | Structured | Per-sentence concreteness rating (1-5 scale) |
| `voice_consistency` | Structured | Passive voice ratio and variance across paragraphs |
| `figurative_language` | Structured | Similes, rhetorical questions, alliteration |
| `rhetorical_diversity` | Structured | Shannon entropy of sentence type distribution |
| `claim_support_ratio` | Structured | Ratio of evidence sentences to claim sentences |
| `counterargument_engagement` | Structured | Density of counterargument phrases |
| `audience_awareness` | Structured | Second-person pronouns, questions, jargon density |
| `argument_structure` | Structured | Thesis-evidence-conclusion marker detection |
| `tension_and_resolution` | Structured | Narrative arc from tension to resolution markers |
| `stance_analysis` | Structured | Pronoun-based stance (reader-centric vs impersonal) |
| `contraction_rate` | Structured | Count and rate of contractions |
| `temporal_orientation` | Structured | Future modal, past tense, evidential, and aspiration densities |
| `economy_analysis` | Structured | Wordy phrase density, average sentence length, subordination index |

### Semantic (5 metrics)

| Metric | Output | Description |
|--------|--------|-------------|
| `topic_coherence` | Structured | Consecutive paragraph similarity via word embeddings |
| `semantic_progression` | Structured | Topic drift rate between paragraphs |
| `redundancy_detection` | Structured | Non-adjacent paragraph pairs above similarity threshold (default 0.90) |
| `information_novelty` | Structured | Per-paragraph novelty relative to prior paragraphs |
| `emotional_tone` | Structured | Valence and arousal scores via seed-word projection (Russell circumplex) |

Semantic metrics require a word2vec model (~310 MB). Download it with:

```bash
inkcheck model download
```

The model is stored in `~/.inkcheck/models/` by default. Set `INKCHECK_MODEL_PATH`
to use a custom location, or change `model_dir` in the config file.

## Interpreting Results

The metrics provide insights into different aspects of text characteristics:

- **Variance metrics** measure consistency vs. variety in paragraph and sentence lengths
- **Opener diversity** indicates how repetitive or varied sentence beginnings are
- **Punctuation variety** shows the range of punctuation types used
- **Paragraph structure** reveals whether opening/closing paragraphs differ from body paragraphs
- **Vocabulary metrics** analyze word choice patterns, including formal vocabulary usage
- **Specificity range** measures variation in abstraction levels across sentences
- **Voice consistency** tracks how passive/active voice usage varies throughout the text

Lower variance and higher uniformity typically indicate more formulaic writing,
while higher variance and variety often suggest more dynamic, varied writing styles.

## Credits

- Built using [Claude Code](https://claude.ai/claude-code)
- Licensed under [MIT License](LICENSE)
- See [CREDITS.md](CREDITS.md) for third-party attributions
