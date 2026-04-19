package bridge

import (
	"strings"

	ftnmodel "screenjson/export/internal/formats/fountain/model"
	"screenjson/export/internal/model"
)

// FromScreenJSON converts a ScreenJSON document to a Fountain document.
func FromScreenJSON(doc *model.Document, lang string) *ftnmodel.Document {
	if lang == "" {
		lang = "en"
	}

	ftn := &ftnmodel.Document{}

	// Build title page
	ftn.TitlePage = buildTitlePage(doc, lang)

	// Build character lookup
	charMap := make(map[string]string) // UUID -> name
	for _, c := range doc.Characters {
		charMap[c.ID] = c.Name
	}

	// Convert content
	ftn.Elements = convertElements(doc, charMap, lang)

	return ftn
}

// buildTitlePage creates a Fountain title page from ScreenJSON metadata.
func buildTitlePage(doc *model.Document, lang string) *ftnmodel.TitlePage {
	tp := &ftnmodel.TitlePage{
		Custom: make(map[string]string),
	}

	// Title
	if title := doc.Title.GetOrDefault(lang); title != "" {
		tp.Title = title
	}

	// Author(s)
	if len(doc.Authors) > 0 {
		var authors []string
		for _, a := range doc.Authors {
			name := strings.TrimSpace(a.Given + " " + a.Family)
			if name != "" {
				authors = append(authors, name)
			}
		}
		if len(authors) == 1 {
			tp.Author = authors[0]
		} else if len(authors) > 1 {
			tp.Authors = strings.Join(authors, "\n       ")
		}
	}

	// Source
	if len(doc.Sources) > 0 {
		var sources []string
		for _, s := range doc.Sources {
			if title := s.Title.GetOrDefault(lang); title != "" {
				sources = append(sources, title)
			}
		}
		if len(sources) > 0 {
			tp.Source = "Based on " + strings.Join(sources, ", ")
		}
	}

	// Copyright
	if doc.License != nil {
		tp.Copyright = doc.License.ID
	}

	return tp
}

// convertElements converts ScreenJSON scenes to Fountain elements.
func convertElements(doc *model.Document, charMap map[string]string, lang string) []ftnmodel.Element {
	var elements []ftnmodel.Element

	if doc.Content == nil {
		return elements
	}

	for i, scene := range doc.Content.Scenes {
		// Add blank line between scenes (except first)
		if i > 0 {
			elements = append(elements, ftnmodel.Element{Type: ftnmodel.ElementBlank})
		}

		// Add scene heading
		if scene.Heading != nil {
			elements = append(elements, ftnmodel.Element{
				Type: ftnmodel.ElementSceneHeading,
				Text: formatSlugline(scene.Heading),
			})
		}

		// Add body elements
		for _, elem := range scene.Body {
			ftnElem := convertElement(elem, charMap, lang)
			if ftnElem != nil {
				elements = append(elements, *ftnElem)
			}
		}
	}

	return elements
}

// formatSlugline formats a ScreenJSON slugline for Fountain.
func formatSlugline(slug *model.Slugline) string {
	var parts []string

	context := slug.Context
	if context == "" {
		context = "INT"
	}
	parts = append(parts, context)

	if slug.Setting != "" {
		parts = append(parts, slug.Setting)
	}

	result := strings.Join(parts, ". ")

	if slug.Time != "" {
		result += " - " + slug.Time
	}

	return result
}

// convertElement converts a ScreenJSON element to a Fountain element.
func convertElement(elem model.Element, charMap map[string]string, lang string) *ftnmodel.Element {
	switch elem.Type {
	case model.ElementAction:
		return &ftnmodel.Element{
			Type: ftnmodel.ElementAction,
			Text: elem.Text.GetOrDefault(lang),
		}

	case model.ElementCharacter:
		display := elem.Display
		if display == "" {
			if name, ok := charMap[elem.Character]; ok {
				display = strings.ToUpper(name)
			}
		}
		return &ftnmodel.Element{
			Type: ftnmodel.ElementCharacter,
			Text: display,
			Dual: elem.Dual,
		}

	case model.ElementDialogue:
		return &ftnmodel.Element{
			Type: ftnmodel.ElementDialogue,
			Text: elem.Text.GetOrDefault(lang),
		}

	case model.ElementParenthetical:
		text := elem.Text.GetOrDefault(lang)
		// Ensure parentheses
		if !strings.HasPrefix(text, "(") {
			text = "(" + text
		}
		if !strings.HasSuffix(text, ")") {
			text = text + ")"
		}
		return &ftnmodel.Element{
			Type: ftnmodel.ElementParenthetical,
			Text: text,
		}

	case model.ElementTransition:
		return &ftnmodel.Element{
			Type: ftnmodel.ElementTransition,
			Text: elem.Text.GetOrDefault(lang),
		}

	case model.ElementShot:
		// Shots are typically ALL CAPS action in Fountain
		return &ftnmodel.Element{
			Type: ftnmodel.ElementAction,
			Text: strings.ToUpper(elem.Text.GetOrDefault(lang)),
		}

	case model.ElementGeneral:
		return &ftnmodel.Element{
			Type: ftnmodel.ElementAction,
			Text: elem.Text.GetOrDefault(lang),
		}
	}

	return nil
}
