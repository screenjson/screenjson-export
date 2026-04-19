// Package model defines Fountain screenplay format structures.
package model

// Document represents a Fountain screenplay.
type Document struct {
	TitlePage  *TitlePage
	Elements   []Element
}

// TitlePage holds metadata from the Fountain title page.
type TitlePage struct {
	Title     string
	Credit    string
	Author    string
	Authors   string
	Source    string
	DraftDate string
	Contact   string
	Copyright string
	Notes     string
	Custom    map[string]string
}

// ElementType defines the type of Fountain element.
type ElementType int

const (
	ElementSceneHeading ElementType = iota
	ElementAction
	ElementCharacter
	ElementDialogue
	ElementParenthetical
	ElementTransition
	ElementLyrics
	ElementCentered
	ElementPageBreak
	ElementSection
	ElementSynopsis
	ElementNote
	ElementBoneyard
	ElementBlank
)

// Element represents a Fountain document element.
type Element struct {
	Type    ElementType
	Text    string
	Depth   int    // For sections
	SceneNo string // Optional scene number
	Dual    bool   // For dual dialogue
	Forced  bool   // Element type was forced with syntax
}

// Common Fountain syntax patterns
const (
	SceneHeadingPrefixes = "INT|EXT|EST|INT./EXT|INT/EXT|I/E"
	TransitionSuffix     = "TO:"
)
