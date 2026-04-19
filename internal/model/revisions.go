package model

import "time"

// Revision represents a revision event/marker.
type Revision struct {
	ID      string     `json:"id"`
	Parent  string     `json:"parent,omitempty"`
	Index   int        `json:"index"`
	Authors []string   `json:"authors"` // UUIDs
	Label   string     `json:"label"`
	Created *time.Time `json:"created"`
	Meta    Meta       `json:"meta,omitempty"`
}

// RevisionColor represents the industry-standard revision color system.
type RevisionColor string

const (
	RevisionWhite     RevisionColor = "white"
	RevisionBlue      RevisionColor = "blue"
	RevisionPink      RevisionColor = "pink"
	RevisionYellow    RevisionColor = "yellow"
	RevisionGreen     RevisionColor = "green"
	RevisionGoldenrod RevisionColor = "goldenrod"
	RevisionBuff      RevisionColor = "buff"
	RevisionSalmon    RevisionColor = "salmon"
	RevisionCherry    RevisionColor = "cherry"
)

// RevisionColorOrder defines the standard revision color sequence.
var RevisionColorOrder = []RevisionColor{
	RevisionWhite,
	RevisionBlue,
	RevisionPink,
	RevisionYellow,
	RevisionGreen,
	RevisionGoldenrod,
	RevisionBuff,
	RevisionSalmon,
	RevisionCherry,
}

// NewRevision creates a new revision.
func NewRevision(id string, index int, authors []string, label string, created time.Time) Revision {
	return Revision{
		ID:      id,
		Index:   index,
		Authors: authors,
		Label:   label,
		Created: &created,
	}
}
