package model

import "time"

// Analysis contains derived machine-readable representations.
// This is optional and discardable.
type Analysis struct {
	Embeddings map[string][]Embedding `json:"embeddings,omitempty"`
	Passages   []Passage              `json:"passages,omitempty"`
	Summaries  []Summary              `json:"summaries,omitempty"`
	Settings   *AnalysisSettings      `json:"settings,omitempty"`
	Meta       Meta                   `json:"meta,omitempty"`
}

// Embedding represents a derived numeric representation for semantic search.
type Embedding struct {
	ID         string     `json:"id"`
	Model      string     `json:"model"`
	Dimensions int        `json:"dimensions"`
	Values     []float64  `json:"values"`
	Source     string     `json:"source"` // name, text, desc, heading, composite
	Lang       string     `json:"lang,omitempty"`
	Tokens     int        `json:"tokens,omitempty"`
	Created    *time.Time `json:"created"`
	Meta       Meta       `json:"meta,omitempty"`
}

// EmbeddingSource defines what was embedded.
type EmbeddingSource string

const (
	EmbeddingSourceName      EmbeddingSource = "name"
	EmbeddingSourceText      EmbeddingSource = "text"
	EmbeddingSourceDesc      EmbeddingSource = "desc"
	EmbeddingSourceHeading   EmbeddingSource = "heading"
	EmbeddingSourceComposite EmbeddingSource = "composite"
)

// Passage represents a derived passage for retrieval.
type Passage struct {
	ID       string   `json:"id"`
	Scene    string   `json:"scene"`
	Elements []string `json:"elements"` // Element UUIDs
	Text     Text     `json:"text"`
	Tokens   int      `json:"tokens"`
	Overlap  int      `json:"overlap,omitempty"`
	Meta     Meta     `json:"meta,omitempty"`
}

// Summary represents a derived summary.
type Summary struct {
	ID        string     `json:"id"`
	Scope     string     `json:"scope"` // document or scene
	Target    string     `json:"target,omitempty"` // Scene UUID if scope is scene
	Text      Text       `json:"text"`
	Generated bool       `json:"generated,omitempty"`
	Model     string     `json:"model,omitempty"`
	Created   *time.Time `json:"created"`
	Meta      Meta       `json:"meta,omitempty"`
}

// SummaryScope defines the scope of a summary.
type SummaryScope string

const (
	SummaryScopeDocument SummaryScope = "document"
	SummaryScopeScene    SummaryScope = "scene"
)

// AnalysisSettings records how derived data was produced.
type AnalysisSettings struct {
	Model     string `json:"model,omitempty"`
	Size      int    `json:"size,omitempty"`      // Passage size target (tokens)
	Overlap   int    `json:"overlap,omitempty"`   // Passage overlap (tokens)
	Tokeniser string `json:"tokeniser,omitempty"`
	Meta      Meta   `json:"meta,omitempty"`
}
