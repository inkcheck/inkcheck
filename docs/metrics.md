# Metrics Reference

## Structure Metrics

These metrics measure the physical shape and rhythm of writing, including
variance in paragraph and sentence lengths, punctuation usage, and structural patterns.

### paragraph_variance

Coefficient of variation (standard deviation / mean) of paragraph word counts.

- **High variance (0.4-0.8+):** paragraphs vary significantly in length
- **Low variance (below 0.25):** paragraphs are roughly uniform in length

### sentence_length_variance

Coefficient of variation of sentence word counts.

- **High variance (0.4-0.7+):** sentences vary widely (4 to 45+ words)
- **Low variance (below 0.3):** sentences cluster around a consistent length (15-25 words)

### sentence_opener_diversity

Ratio of unique sentence openers (first two words, lowercased) to total
sentences.

- **High diversity (0.8-1.0):** varied sentence openings
- **Low diversity (below 0.7):** repetitive openings (e.g., "This", "The", "However", "It")

### paragraph_position_analysis

Compares the first and last paragraphs against the body paragraphs. Typically,
introductions and conclusions differ structurally from body paragraphs, while
more uniform writing has less variation between opening/closing and body sections.

Returns:
- Opening and closing paragraph word counts
- Body mean word count
- Deviation ratios (how far opening/closing deviate from body mean)
- Uniform flag (true if both within 20% of body mean)

### punctuation_profile

Counts 8 punctuation types: periods, commas, semicolons, colons, dashes,
parentheses, question marks, exclamation marks.

- **High variety (6-8 types):** diverse punctuation including semicolons, dashes, parenthetical asides
- **Low variety (3-4 types):** primarily periods and commas

## Rhetoric Metrics

These metrics analyze writing style, argument structure, and rhetorical
techniques.

### transition_word_density

Counts 32 transition phrases (e.g. "however", "furthermore", "for example") and
measures variety (distinct / total).

- **Low variety with high count:** over-reliance on a few transitions like "however",
  "additionally", "furthermore", "moreover"
- **Repeated:** lists transitions appearing 3 or more times

### vocabulary_sophistication

Analyzes word frequency distribution across 5 bands (most common to rarest)
using the Google 10,000 English words list.

- **TTR:** Type-Token Ratio (unique words / total words)
- **MATTR:** Moving-Average TTR (window=50), more stable across text lengths
- **Band CV:** coefficient of variation across 5 frequency bands
- **Formal words:** count of words from a curated list of 50+ overly formal/corporate
  vocabulary ("utilize", "leverage", "nuanced", "delve", "tapestry")

### hedging_analysis

Detects hedging language across 5 categories:

| Category | Examples |
|----------|----------|
| Modal | may, might, could |
| Approximator | approximately, roughly, somewhat |
| Plausibility | perhaps, probably, possibly, likely |
| Attribution | it seems, appears to, tends to |
| Frequency | sometimes, often, generally, usually |

Returns density (hedges per 100 words), distinct forms, and variety ratio.

### specificity_score

Rates each sentence on a 1-5 scale from abstract to concrete using heuristics:
numbers, percentages, proper nouns, quotes, signal phrases.

- **High variance:** wide range of specificity (1-5), varied abstraction levels
- **Low variance:** narrow range, sentences cluster around similar abstraction levels (2-3)

### voice_consistency

Detects passive voice using POS tag patterns (be-verb + past participle) and
measures the coefficient of variation of passive ratios across paragraphs.

- **High CV:** voice varies naturally across paragraphs
- **Low CV:** uniform passive/active ratio throughout

### figurative_language

Detects three types of figurative language:

- **Similes:** "like a/an", "as ... as" patterns
- **Rhetorical questions:** sentences ending with "?"
- **Alliteration:** 3+ consecutive words starting with the same letter

Returns instance counts and density per 100 words.

### rhetorical_diversity

Classifies sentences into 5 types: questions, exclamations, imperatives,
conditionals, and declaratives. Computes Shannon entropy of the distribution.

- **High entropy:** diverse sentence types
- **Low entropy:** primarily declaratives (formal writing style)

### claim_support_ratio

Classifies sentences as claims (containing "should", "must", "clearly"),
support (containing "for example", "according to", "research shows"), or
neutral. Returns the support-to-claim ratio.

Well-argued writing has a ratio above 1.0 (more evidence than claims).

### counterargument_engagement

Counts phrases that signal engagement with opposing viewpoints: "on the other
hand", "critics argue", "while it is true", "admittedly", etc.

Returns instance count, density per 100 words, and which phrases were found.

### audience_awareness

Measures audience engagement through:

- Second-person pronoun count ("you", "your")
- Direct question count
- Parenthetical aside count
- Jargon density (words outside top 6,000 frequency)
- Composite engagement score

### argument_structure

Checks for thesis-evidence-conclusion structure by scanning for marker phrases
at expected positions in the text.

- Thesis markers: "I argue", "the purpose", "in this article"
- Evidence markers: "for example", "according to", "research shows"
- Conclusion markers: "in conclusion", "to summarize", "therefore"

Coherence score rewards thesis appearing early and conclusion appearing late.

### tension_and_resolution

Tracks tension markers ("the problem", "however", "remains unclear") and
resolution markers ("the solution", "this demonstrates", "in conclusion")
to detect narrative arc structure.

- **Has arc:** tension markers appear before resolution markers
- **Arc score:** 0-1 based on presence and ordering of markers

## Semantic Metrics

These metrics use word2vec embeddings to measure meaning-level patterns.
They require a pre-trained model (~300 MB, downloaded on first use).

### topic_coherence

Computes cosine similarity between consecutive paragraph embeddings (average
of word vectors). Returns mean similarity and CV.

- **High mean + low CV:** topic flows very smoothly with consistent transitions
- **Variable similarities:** topic progression with varying degrees of connection

### semantic_progression

Measures drift rate (1 - similarity) between consecutive paragraphs.

- **High mean drift:** topics change significantly between paragraphs
- **CV of drifts:** how consistently the topic shifts

### redundancy_detection

Finds all non-adjacent paragraph pairs with cosine similarity above 0.85.
These represent paragraphs that say essentially the same thing in different
parts of the text.

### information_novelty

Computes per-paragraph novelty as 1 - max(similarity to all prior paragraphs).
The first paragraph always gets novelty 1.0.

- **Declining curve:** natural as topics build on prior content
- **Flat low novelty:** text is repetitive
- **Flat high novelty:** text lacks coherence
