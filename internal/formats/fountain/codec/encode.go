package codec

import (
	"context"
	"fmt"
	"strings"

	ftnmodel "screenjson/export/internal/formats/fountain/model"
)

// Encoder encodes to Fountain files.
type Encoder struct{}

// NewEncoder creates a new Fountain encoder.
func NewEncoder() *Encoder {
	return &Encoder{}
}

// Encode serializes the Fountain model to Fountain text.
func (e *Encoder) Encode(ctx context.Context, doc *ftnmodel.Document) ([]byte, error) {
	if doc == nil {
		return nil, fmt.Errorf("nil Fountain document")
	}
	var sb strings.Builder
	
	// Write title page
	if doc.TitlePage != nil {
		writeTitlePage(&sb, doc.TitlePage)
		sb.WriteString("\n")
	}
	
	// Write elements
	for i, elem := range doc.Elements {
		writeElement(&sb, elem, i > 0)
	}
	
	return []byte(sb.String()), nil
}

// writeTitlePage writes the title page metadata.
func writeTitlePage(sb *strings.Builder, tp *ftnmodel.TitlePage) {
	if tp.Title != "" {
		sb.WriteString("Title: ")
		sb.WriteString(tp.Title)
		sb.WriteString("\n")
	}
	if tp.Credit != "" {
		sb.WriteString("Credit: ")
		sb.WriteString(tp.Credit)
		sb.WriteString("\n")
	}
	if tp.Author != "" {
		sb.WriteString("Author: ")
		sb.WriteString(tp.Author)
		sb.WriteString("\n")
	}
	if tp.Authors != "" {
		sb.WriteString("Authors: ")
		sb.WriteString(tp.Authors)
		sb.WriteString("\n")
	}
	if tp.Source != "" {
		sb.WriteString("Source: ")
		sb.WriteString(tp.Source)
		sb.WriteString("\n")
	}
	if tp.DraftDate != "" {
		sb.WriteString("Draft date: ")
		sb.WriteString(tp.DraftDate)
		sb.WriteString("\n")
	}
	if tp.Contact != "" {
		sb.WriteString("Contact: ")
		sb.WriteString(tp.Contact)
		sb.WriteString("\n")
	}
	if tp.Copyright != "" {
		sb.WriteString("Copyright: ")
		sb.WriteString(tp.Copyright)
		sb.WriteString("\n")
	}
	if tp.Notes != "" {
		sb.WriteString("Notes: ")
		sb.WriteString(tp.Notes)
		sb.WriteString("\n")
	}
	for key, value := range tp.Custom {
		sb.WriteString(key)
		sb.WriteString(": ")
		sb.WriteString(value)
		sb.WriteString("\n")
	}
}

// writeElement writes a single Fountain element.
func writeElement(sb *strings.Builder, elem ftnmodel.Element, needsBlank bool) {
	switch elem.Type {
	case ftnmodel.ElementBlank:
		sb.WriteString("\n")
		
	case ftnmodel.ElementSceneHeading:
		if needsBlank {
			sb.WriteString("\n")
		}
		if elem.Forced && !sceneHeadingPattern.MatchString(elem.Text) {
			sb.WriteString(".")
		}
		sb.WriteString(elem.Text)
		sb.WriteString("\n")
		
	case ftnmodel.ElementAction:
		if needsBlank {
			sb.WriteString("\n")
		}
		if elem.Forced {
			sb.WriteString("!")
		}
		sb.WriteString(elem.Text)
		sb.WriteString("\n")
		
	case ftnmodel.ElementCharacter:
		if needsBlank {
			sb.WriteString("\n")
		}
		if elem.Forced {
			sb.WriteString("@")
		}
		sb.WriteString(elem.Text)
		if elem.Dual {
			sb.WriteString(" ^")
		}
		sb.WriteString("\n")
		
	case ftnmodel.ElementDialogue:
		sb.WriteString(elem.Text)
		sb.WriteString("\n")
		
	case ftnmodel.ElementParenthetical:
		text := elem.Text
		if !strings.HasPrefix(text, "(") {
			text = "(" + text
		}
		if !strings.HasSuffix(text, ")") {
			text = text + ")"
		}
		sb.WriteString(text)
		sb.WriteString("\n")
		
	case ftnmodel.ElementTransition:
		if needsBlank {
			sb.WriteString("\n")
		}
		if elem.Forced && !transitionPattern.MatchString(elem.Text) {
			sb.WriteString("> ")
		}
		sb.WriteString(elem.Text)
		sb.WriteString("\n")
		
	case ftnmodel.ElementCentered:
		if needsBlank {
			sb.WriteString("\n")
		}
		sb.WriteString(">")
		sb.WriteString(elem.Text)
		sb.WriteString("<\n")
		
	case ftnmodel.ElementSection:
		if needsBlank {
			sb.WriteString("\n")
		}
		for i := 0; i < elem.Depth; i++ {
			sb.WriteString("#")
		}
		sb.WriteString(" ")
		sb.WriteString(elem.Text)
		sb.WriteString("\n")
		
	case ftnmodel.ElementSynopsis:
		sb.WriteString("= ")
		sb.WriteString(elem.Text)
		sb.WriteString("\n")
		
	case ftnmodel.ElementPageBreak:
		sb.WriteString("\n===\n")
		
	case ftnmodel.ElementLyrics:
		sb.WriteString("~")
		sb.WriteString(elem.Text)
		sb.WriteString("\n")
		
	case ftnmodel.ElementNote:
		sb.WriteString("[[")
		sb.WriteString(elem.Text)
		sb.WriteString("]]\n")
		
	case ftnmodel.ElementBoneyard:
		sb.WriteString("/*")
		sb.WriteString(elem.Text)
		sb.WriteString("*/\n")
	}
}
