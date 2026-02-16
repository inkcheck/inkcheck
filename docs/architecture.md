# Architecture

## Package Structure

```
inkcheck/
├── cmd/
│   └── inkcheck/
│       └── main.go              # CLI entry point
├── docs/                        # Documentation
├── shared/
│   ├── text.go                  # Markdown extraction, sentence splitting, tokenization
│   ├── stats.go                 # Mean, StdDev, CV, CosineSimilarity, Entropy
│   └── token.go                 # Token type (Text, Tag, Label)
├── structure/
│   ├── paragraph_variance.go
│   ├── sentence_length_variance.go
│   ├── sentence_opener_diversity.go
│   ├── paragraph_position_analysis.go
│   └── punctuation_profile.go
├── rhetoric/
│   ├── wordlist.go              # 10,000-word frequency map (generated)
│   ├── transition_word_density.go
│   ├── vocabulary_sophistication_distribution.go
│   ├── hedging_analysis.go
│   ├── specificity_score.go
│   ├── voice_consistency.go
│   ├── figurative_language_presence.go
│   ├── rhetorical_diversity.go
│   ├── claim_support_ratio.go
│   ├── counterargument_engagement.go
│   ├── audience_awareness.go
│   ├── argument_structure_coherence.go
│   └── tension_and_resolution.go
├── semantic/
│   ├── model.go                 # Word2vec model download, cache, and loading
│   ├── topic_coherence.go
│   ├── semantic_progression.go
│   ├── redundancy_detection.go
│   └── information_novelty_curve.go
├── go.mod
└── go.sum
```

## Source Files

### shared/

| File | Purpose |
|------|---------|
| `text.go` | Markdown parsing via goldmark, sentence splitting via prose, POS tokenization |
| `stats.go` | Statistical functions: mean, standard deviation, CV, cosine similarity, entropy |
| `token.go` | `Token` struct with Text, Tag (POS), and Label (NER) fields |

### structure/

| File | Purpose |
|------|---------|
| `paragraph_variance.go` | CV of paragraph word counts |
| `sentence_length_variance.go` | CV of sentence word counts |
| `sentence_opener_diversity.go` | Unique opener ratio with first-two-word extraction |
| `paragraph_position_analysis.go` | Opening/closing vs body comparison |
| `punctuation_profile.go` | 8-type punctuation counting with Variety() and Total() methods |

### rhetoric/

| File | Purpose |
|------|---------|
| `wordlist.go` | Generated map of 10,000 words to frequency rank |
| `transition_word_density.go` | 32 transition phrases with word-boundary matching |
| `vocabulary_sophistication_distribution.go` | TTR, MATTR, frequency bands, formal word detection |
| `hedging_analysis.go` | 5-category hedging detection with 28 patterns |
| `specificity_score.go` | Per-sentence concreteness heuristic (1-5 scale) |
| `voice_consistency.go` | POS-based passive voice detection per paragraph |
| `figurative_language_presence.go` | Simile, rhetorical question, alliteration detection |
| `rhetorical_diversity.go` | Sentence type classification with Shannon entropy |
| `claim_support_ratio.go` | Keyword-based claim vs evidence classification |
| `counterargument_engagement.go` | 21 counterargument phrase patterns |
| `audience_awareness.go` | Second-person, questions, parentheticals, jargon density |
| `argument_structure_coherence.go` | Thesis/evidence/conclusion marker position analysis |
| `tension_and_resolution.go` | Tension/resolution marker tracking for narrative arc |

### semantic/

| File | Purpose |
|------|---------|
| `model.go` | Word2vec model management: download, cache, load, embedding functions |
| `topic_coherence.go` | Consecutive paragraph cosine similarity |
| `semantic_progression.go` | Drift rate (1 - similarity) between paragraphs |
| `redundancy_detection.go` | Non-adjacent paragraph pair similarity above threshold |
| `information_novelty_curve.go` | Per-paragraph novelty relative to all prior paragraphs |

## Dependencies

| Dependency | Purpose |
|------------|---------|
| `github.com/tsawler/prose/v3` | Sentence splitting, POS tagging, NER |
| `github.com/yuin/goldmark` | Markdown parsing for prose extraction |
| `github.com/danieldk/go2vec` | Word2vec binary model loading (semantic metrics only) |

The `go2vec` library requires CGO for BLAS operations. On macOS, link against
the Accelerate framework. On Linux, install OpenBLAS.

## Design Decisions

### Package Separation

Metrics are organized by category (structure, rhetoric, semantic) rather than
as a flat list. This allows importing only what you need — if you don't use
semantic metrics, you avoid the word2vec and CGO dependency entirely.

### Shared Text Processing

All packages use `shared.ExtractProse` for markdown-aware text extraction and
`shared.SplitSentences` for sentence boundary detection. This ensures
consistent text preprocessing across all metrics.

### Heuristic-First

The 17 non-semantic metrics use keyword matching, POS patterns, and statistical
variance rather than ML models. This makes them fast, deterministic, and easy
to understand. Metrics that would benefit from deeper analysis are marked with
`// TODO: LLM judge` comments.

### Lazy Model Loading

The word2vec model is only downloaded and loaded when a semantic metric is
requested. Non-semantic metrics work without any model file.

### Word2vec Model Management

The semantic package manages model lifecycle:
1. Check `INKCHECK_MODEL_PATH` environment variable
2. Fall back to `~/.inkcheck/models/`
3. Download if not cached (prints progress to stderr)
4. Load into memory via go2vec

## Data Flow

```
markdown text
  │
  ├─► shared.ExtractProse()          (goldmark AST)
  │     └─► []paragraphs
  │
  ├─► shared.SplitSentences()        (prose tokenizer)
  │     └─► []sentences
  │
  ├─► shared.Tokenize()              (prose POS tagger)
  │     └─► []Token{Text, Tag, Label}
  │
  ├─► structure.*()                  (statistical variance)
  │     └─► CV, ratios, counts
  │
  ├─► rhetoric.*()                   (keyword + POS heuristics)
  │     └─► densities, ratios, classifications
  │
  └─► semantic.*()                   (word2vec embeddings)
        └─► similarities, drift rates, novelty scores
```
