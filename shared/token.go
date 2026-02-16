package shared

// Token represents a POS-tagged token from prose analysis.
type Token struct {
	Text  string // The token text.
	Tag   string // POS tag (e.g., "NN", "VB", "JJ").
	Label string // NER label (e.g., "PERSON", "GPE"), empty if not an entity.
}
