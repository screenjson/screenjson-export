package model

import "time"

// ElementType defines the type of screenplay element.
type ElementType string

const (
	ElementAction        ElementType = "action"
	ElementCharacter     ElementType = "character"
	ElementDialogue      ElementType = "dialogue"
	ElementParenthetical ElementType = "parenthetical"
	ElementTransition    ElementType = "transition"
	ElementShot          ElementType = "shot"
	ElementGeneral       ElementType = "general"
)

// Element represents a screenplay body element.
// This is a unified type that covers all element types.
type Element struct {
	// Common fields
	ID           string       `json:"id"`
	Type         ElementType  `json:"type"`
	Scene        string       `json:"scene,omitempty"`
	Authors      []string     `json:"authors"` // UUIDs
	Contributors []string     `json:"contributors,omitempty"`
	Access       []string     `json:"access,omitempty"`
	Notes        []Note       `json:"notes,omitempty"`
	Charset      string       `json:"charset,omitempty"`
	Dir          string       `json:"dir,omitempty"`
	Class        string       `json:"class,omitempty"`
	DOM          string       `json:"dom,omitempty"`
	Encrypt      *Encrypt     `json:"encrypt,omitempty"`
	Locked       bool         `json:"locked,omitempty"`
	Omit         bool         `json:"omit,omitempty"`
	Revisions    []Revision   `json:"revisions,omitempty"`
	Styles       []string     `json:"styles,omitempty"`
	Meta         Meta         `json:"meta,omitempty"`

	// Text content (for action, dialogue, parenthetical, transition, shot, general)
	Text Text `json:"text,omitempty"`

	// Character-specific fields
	Character string `json:"character,omitempty"` // UUID for cue and dialogue
	Display   string `json:"display,omitempty"`   // Display override for cue

	// Dialogue-specific fields
	Origin string `json:"origin,omitempty"` // V.O., O.S., O.C., FILTER
	Dual   bool   `json:"dual,omitempty"`   // Dual dialogue flag

	// Shot-specific fields
	FOV         float64 `json:"fov,omitempty"`
	Perspective string  `json:"perspective,omitempty"` // 2D or 3D
}

// Note represents an ancillary note attached to an element.
type Note struct {
	ID          string     `json:"id"`
	Parent      string     `json:"parent,omitempty"`
	Highlight   [][2]int   `json:"highlight,omitempty"` // Character index ranges
	Contributor string     `json:"contributor,omitempty"`
	Created     *time.Time `json:"created"`
	Text        Text       `json:"text"`
	Color       string     `json:"color,omitempty"`
	Meta        Meta       `json:"meta,omitempty"`
}

// NewAction creates a new action element.
func NewAction(id string, authors []string, text Text) Element {
	return Element{
		ID:      id,
		Type:    ElementAction,
		Authors: authors,
		Text:    text,
	}
}

// NewCharacterCue creates a new character cue element.
func NewCharacterCue(id string, authors []string, characterID string, display string) Element {
	return Element{
		ID:        id,
		Type:      ElementCharacter,
		Authors:   authors,
		Character: characterID,
		Display:   display,
	}
}

// NewDialogue creates a new dialogue element.
func NewDialogue(id string, authors []string, characterID string, text Text, origin string, dual bool) Element {
	return Element{
		ID:        id,
		Type:      ElementDialogue,
		Authors:   authors,
		Character: characterID,
		Text:      text,
		Origin:    origin,
		Dual:      dual,
	}
}

// NewParenthetical creates a new parenthetical element.
func NewParenthetical(id string, authors []string, text Text) Element {
	return Element{
		ID:      id,
		Type:    ElementParenthetical,
		Authors: authors,
		Text:    text,
	}
}

// NewTransition creates a new transition element.
func NewTransition(id string, authors []string, text Text) Element {
	return Element{
		ID:      id,
		Type:    ElementTransition,
		Authors: authors,
		Text:    text,
	}
}

// NewShot creates a new shot element.
func NewShot(id string, authors []string, text Text) Element {
	return Element{
		ID:      id,
		Type:    ElementShot,
		Authors: authors,
		Text:    text,
		FOV:     40,
	}
}

// NewGeneral creates a new general element.
func NewGeneral(id string, authors []string, text Text) Element {
	return Element{
		ID:      id,
		Type:    ElementGeneral,
		Authors: authors,
		Text:    text,
	}
}
