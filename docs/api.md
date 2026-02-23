# Library API

## inkcheck (top-level package)

```go
import "github.com/inkcheck"
```

The top-level package provides a unified API for running all metrics at once.

### Result Types

```go
type Result struct {
    Structure StructureResult
    Rhetoric  RhetoricResult
    Semantic  SemanticResult
}

type StructureResult struct {
    ParagraphVariance        float64
    ParagraphLengths         []int
    SentenceLengthVariance   float64
    SentenceOpenerDiversity  structure.SentenceOpenerDiversityResult
    SentenceTypeDistribution structure.SentenceTypeResult
    ParagraphPosition        structure.ParagraphPositionResult
    Punctuation              structure.PunctuationProfile
}

type RhetoricResult struct {
    TransitionWordDensity   rhetoric.TransitionResult
    VocabSophistication     rhetoric.VocabSophisticationResult
    Hedging                 rhetoric.HedgingResult
    Specificity             rhetoric.SpecificityResult
    VoiceConsistency        rhetoric.VoiceConsistencyResult
    FigurativeLanguage      rhetoric.FigurativeLanguageResult
    RhetoricalDiversity     rhetoric.RhetoricalDiversityResult
    ClaimSupport            rhetoric.ClaimSupportResult
    Counterargument         rhetoric.CounterargumentResult
    AudienceAwareness       rhetoric.AudienceAwarenessResult
    ArgumentStructure       rhetoric.ArgumentStructureResult
    TensionAndResolution    rhetoric.TensionResolutionResult
    Stance                  rhetoric.StanceResult
    Contraction             rhetoric.ContractionResult
    Temporal                rhetoric.TemporalResult
    Economy                 rhetoric.EconomyResult
}

type SemanticResult struct {
    TopicCoherence      semantic.TopicCoherenceResult
    SemanticProgression semantic.SemanticProgressionResult
    Redundancy          semantic.RedundancyResult
    InformationNovelty  semantic.InformationNoveltyCurveResult
    EmotionalTone       semantic.EmotionalToneResult
}
```

### Analyze

```go
func Analyze(text string) Result
```

Runs all structure and rhetoric metrics (22 total). Semantic metrics are left
empty. This is the simplest entry point — no model download required.

```go
result := inkcheck.Analyze(markdownText)
fmt.Println(result.Structure.ParagraphVariance)
fmt.Println(result.Rhetoric.Hedging.Density)
```

### AnalyzeAll

```go
func AnalyzeAll(text string, model *semantic.ModelManager) Result
```

Runs all 27 metrics including semantic analysis. Requires a loaded
`ModelManager`.

```go
model, err := semantic.LoadModel()
if err != nil {
    log.Fatal(err)
}
result := inkcheck.AnalyzeAll(text, model)
fmt.Println(result.Semantic.TopicCoherence.MeanSimilarity)
```

### AnalyzeStructure

```go
func AnalyzeStructure(text string) StructureResult
```

Runs all 6 structural metrics.

### AnalyzeRhetoric

```go
func AnalyzeRhetoric(text string) RhetoricResult
```

Runs all 16 rhetorical metrics.

### AnalyzeSemantic

```go
func AnalyzeSemantic(text string, model *semantic.ModelManager) SemanticResult
```

Runs all 5 semantic metrics.

## Sub-packages

For individual metric access, import sub-packages directly:

```go
import (
    "github.com/inkcheck/shared"     // text processing and statistics
    "github.com/inkcheck/structure"  // structural metrics
    "github.com/inkcheck/rhetoric"   // rhetorical metrics
    "github.com/inkcheck/semantic"   // semantic metrics (requires word2vec model)
)
```

## shared

Common text processing utilities used by all other packages.

### Text Functions

```go
func ExtractProse(markdown string) []string
```

Parses markdown and returns prose paragraph text, stripping headings, code
blocks, lists, block quotes, tables, and HTML blocks.

```go
func ExtractProseText(markdown string) string
```

Returns all prose paragraphs joined into a single string.

```go
func SplitParagraphs(markdown string) []string
```

Alias for `ExtractProse`.

```go
func SplitSentences(text string) []string
```

Splits text into sentences using prose's NLP sentence tokenizer.

```go
func ListWords(text string) []string
func CountWords(text string) int
```

Whitespace-based word splitting and counting.

```go
func Tokenize(text string) []Token
```

Returns POS-tagged tokens using prose's NLP pipeline.

### Token

```go
type Token struct {
    Text  string // The token text
    Tag   string // POS tag (e.g., "NN", "VB", "JJ")
    Label string // NER label (e.g., "PERSON", "GPE"), empty if not an entity
}
```

### Statistics Functions

```go
func Mean(values []float64) float64
func StdDev(values []float64) float64
func CoefficientOfVariation(values []float64) float64
func CosineSimilarity(a, b []float32) float64
func Entropy(distribution []float64) float64
```

## structure

Six structural metrics. All functions take markdown text as input.

### ParagraphVariance

```go
func ParagraphVariance(text string) float64
func ParagraphLengths(text string) []int
```

Returns the coefficient of variation of paragraph word counts. Low CV (below
0.25) suggests uniform paragraph lengths.

### SentenceLengthVariance

```go
func SentenceLengthVariance(text string) float64
```

Returns the CV of sentence word counts. Higher CV indicates more varied
sentence lengths (e.g., 4 to 45 words), while lower CV suggests consistent
sentence lengths (e.g., clustering around 15-25 words).

### SentenceOpenerDiversity

```go
func SentenceOpenerDiversity(cfg config.Config, text string) SentenceOpenerDiversityResult
```

```go
type SentenceOpenerDiversityResult struct {
    Ratio   float64 // distinct openers / total sentences (0–1)
    Entropy float64 // Shannon entropy of opener distribution (bits)
}
```

Returns both the unique-opener ratio and the Shannon entropy of the opener
frequency distribution.

### SentenceTypeDistribution

```go
func SentenceTypeDistribution(text string) SentenceTypeResult
```

```go
type SentenceTypeResult struct {
    Declarative   int
    Interrogative int
    Imperative    int
    Exclamatory   int
    Total         int
    Entropy       float64 // Shannon entropy (bits); max ≈ 2.0
}
```

Classifies each sentence as declarative, interrogative, imperative, or
exclamatory and computes the Shannon entropy of the distribution.

### ParagraphPositionAnalysis

```go
func ParagraphPositionAnalysis(text string) ParagraphPositionResult
```

```go
type ParagraphPositionResult struct {
    OpeningLength    int
    ClosingLength    int
    BodyMeanLength   float64
    OpeningDeviation float64
    ClosingDeviation float64
    Uniform          bool    // true if both within 20% of body mean
}
```

### PunctuationAnalysis

```go
func PunctuationAnalysis(text string) PunctuationProfile
```

```go
type PunctuationProfile struct {
    Periods, Commas, Semicolons, Colons     int
    Dashes, Parentheses, Questions, Exclamations int
}

func (p PunctuationProfile) Variety() int  // types used out of 8
func (p PunctuationProfile) Total() int
```

## rhetoric

Sixteen rhetorical metrics. All functions take markdown text as input.

### TransitionWordDensity

```go
func TransitionWordDensity(text string) TransitionResult
```

```go
type TransitionResult struct {
    Total    int
    Distinct int
    Variety  float64   // Distinct / Total
    Repeated []string  // phrases appearing 3+ times
}
```

### VocabSophisticationDistribution

```go
func VocabSophisticationDistribution(text string) VocabSophisticationResult
```

```go
type VocabSophisticationResult struct {
    TotalWords      int
    UniqueWords     int
    TypeTokenRatio  float64
    MATTR           float64         // Moving-Average TTR (window=50)
    BandCounts      [5]int          // frequency band counts
    BandRatios      [5]float64
    BandCV          float64
    FormalWordCount int
    FormalWordRatio float64
    FormalWords     map[string]int
}
```

### HedgingAnalysis

```go
func HedgingAnalysis(text string) HedgingResult
```

```go
type HedgingResult struct {
    Total      int
    Density    float64           // hedges per 100 words
    Distinct   int
    Variety    float64
    Categories HedgingCategories // modal, approximator, plausibility, attribution, frequency
    Hedges     []HedgeInstance
}
```

### SpecificityScore

```go
func SpecificityScore(text string) SpecificityResult
```

```go
type SpecificityResult struct {
    Mean   float64
    Range  float64
    CV     float64
    Scores []int   // per-sentence scores (1=abstract, 5=concrete)
}
```

### VoiceConsistency

```go
func VoiceConsistency(text string) VoiceConsistencyResult
```

```go
type VoiceConsistencyResult struct {
    PassiveRatio           float64
    ParagraphPassiveRatios []float64
    CV                     float64
}
```

### FigurativeLanguagePresence

```go
func FigurativeLanguagePresence(text string) FigurativeLanguageResult
```

```go
type FigurativeLanguageResult struct {
    SimileCount, RhetoricalQuestionCount, AlliterationCount int
    TotalInstances   int
    DensityPer100Words float64
}
```

### RhetoricalDiversity

```go
func RhetoricalDiversity(text string) RhetoricalDiversityResult
```

```go
type RhetoricalDiversityResult struct {
    Questions, Exclamations, Imperatives, Conditionals, Declaratives int
    Total   int
    Entropy float64  // Shannon entropy of sentence type distribution
}
```

### ClaimSupportRatio

```go
func ClaimSupportRatio(text string) ClaimSupportResult
```

```go
type ClaimSupportResult struct {
    ClaimCount, SupportCount, NeutralCount int
    Total int
    Ratio float64  // Support / Claim
}
```

### CounterargumentEngagement

```go
func CounterargumentEngagement(text string) CounterargumentResult
```

```go
type CounterargumentResult struct {
    Instances     int
    DensityPer100 float64
    Phrases       []string
}
```

### AudienceAwareness

```go
func AudienceAwareness(text string) AudienceAwarenessResult
```

```go
type AudienceAwarenessResult struct {
    SecondPersonCount, DirectQuestionCount, ParentheticalCount int
    JargonDensity   float64
    EngagementScore float64
}
```

### ArgumentStructureCoherence

```go
func ArgumentStructureCoherence(text string) ArgumentStructureResult
```

```go
type ArgumentStructureResult struct {
    HasThesisMarker, HasEvidenceMarkers, HasConclusionMarker bool
    ThesisPosition, ConclusionPosition float64  // 0.0=start, 1.0=end
    CoherenceScore float64                      // 0-1
}
```

### TensionAndResolution

```go
func TensionAndResolution(text string) TensionResolutionResult
```

```go
type TensionResolutionResult struct {
    TensionMarkers, ResolutionMarkers int
    HasArc    bool     // true if tension precedes resolution
    ArcScore  float64  // 0-1
}
```

### StanceAnalysis

```go
func StanceAnalysis(text string) StanceResult
```

```go
type StanceResult struct {
    SecondPerson     int
    FirstPlural      int
    FirstSingular    int
    ThirdImpersonal  int
    TotalPronouns    int
    ReaderCentricity float64 // weighted pronoun score / total words
}
```

Analyses pronoun-based stance. `ReaderCentricity` uses weighted counts
(you×1.0 + we×0.6 + I×0.4 + they/one×0.1) divided by total words.

### ContractionRate

```go
func ContractionRate(text string) ContractionResult
```

```go
type ContractionResult struct {
    Count int
    Rate  float64 // Count / total words (0.0–1.0)
}
```

Counts contractions (e.g. "don't", "we're", "it's") and returns the rate
per total words.

### TemporalOrientation

```go
func TemporalOrientation(text string) TemporalResult
```

```go
type TemporalResult struct {
    FutureModalCount   int
    PastTenseCount     int
    EvidentialCount    int
    AspirationCount    int
    FutureModalDensity float64 // per 100 words
    PastTenseDensity   float64 // per 100 words
    EvidentialDensity  float64 // per 100 words
    AspirationDensity  float64 // per 100 words
}
```

Analyses the epistemic-temporal mode of the text: future modals, past-tense
indicators, evidential phrases, and aspiration/intention phrases.

### EconomyAnalysis

```go
func EconomyAnalysis(text string) EconomyResult
```

```go
type EconomyResult struct {
    AvgSentenceLength  float64 // mean words per sentence
    WordyPhraseCount   int
    WordyPhraseDensity float64 // wordy phrases per 100 words
    WordsPerClause     float64 // mean words per clause (approximated)
    ClauseCount        int
    SubordinationIndex float64 // subordinating conjunctions per sentence
}
```

Measures writing efficiency: detects wordy/redundant phrases, approximates
clause count via punctuation and subordinating conjunctions, and computes
average sentence length.

## semantic

Five semantic metrics using word2vec embeddings. All functions take a
`*ModelManager` as the first parameter.

### ModelManager

```go
func LoadModel() (*ModelManager, error)
```

Loads the word2vec model from `~/.inkcheck/models/` (or downloads it on first
use). Respects the `INKCHECK_MODEL_PATH` environment variable.

```go
func (m *ModelManager) Embedding(word string) []float32
func (m *ModelManager) SentenceEmbedding(sentence string) []float32
func (m *ModelManager) ParagraphEmbedding(paragraph string) []float32
func (m *ModelManager) Similarity(text1, text2 string) float64
```

### TopicCoherence

```go
func TopicCoherence(m *ModelManager, text string) TopicCoherenceResult
```

```go
type TopicCoherenceResult struct {
    PairSimilarities []float64
    MeanSimilarity   float64
    CV               float64
}
```

### SemanticProgression

```go
func SemanticProgression(m *ModelManager, text string) SemanticProgressionResult
```

```go
type SemanticProgressionResult struct {
    DriftRates []float64
    MeanDrift  float64
    CV         float64
}
```

### RedundancyDetection

```go
func RedundancyDetection(m *ModelManager, text string) RedundancyResult
```

```go
type RedundancyResult struct {
    Pairs     []RedundancyPair  // IndexA, IndexB, Similarity
    PairCount int
}
```

### InformationNoveltyCurve

```go
func InformationNoveltyCurve(m *ModelManager, text string) InformationNoveltyCurveResult
```

```go
type InformationNoveltyCurveResult struct {
    NoveltyScores []float64  // first paragraph = 1.0
    MeanNovelty   float64
    CV            float64
}
```

### EmotionalTone

```go
func EmotionalTone(cfg config.Config, m *ModelManager, text string) EmotionalToneResult
```

```go
type EmotionalToneResult struct {
    Valence      float64 // -1 (very negative) to +1 (very positive)
    Arousal      float64 // -1 (very calm) to +1 (very excited/energetic)
    CoveredWords int     // content words that had embeddings in the model
}
```

Estimates valence and arousal using the Russell circumplex model. Projects
content-word embeddings onto axes defined by positive/negative seed-word
centroids. Returns zero result if the model is nil or no words matched.
