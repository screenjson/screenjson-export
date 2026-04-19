package bridge

import (
	"strings"

	fdxmodel "screenjson/export/internal/formats/fdx/model"
	"screenjson/export/internal/model"
)

// FromScreenJSON converts a ScreenJSON document to an FDX document.
func FromScreenJSON(doc *model.Document, lang string) *fdxmodel.FinalDraft {
	if lang == "" {
		lang = "en"
	}

	fdx := &fdxmodel.FinalDraft{
		DocumentType: "Script",
		Template:     "No",
		Version:      "2",
	}

	// Build character lookup map
	charMap := make(map[string]string) // UUID -> name
	for _, c := range doc.Characters {
		charMap[c.ID] = c.Name
	}

	// Convert content
	var paragraphs []fdxmodel.Paragraph

	if doc.Content != nil {
		for _, scene := range doc.Content.Scenes {
			// Add scene heading
			if scene.Heading != nil {
				headingText := formatSlugline(scene.Heading)
				paragraphs = append(paragraphs, fdxmodel.Paragraph{
					Type: "Scene Heading",
					Text: []fdxmodel.Text{{Content: headingText}},
				})
			}

			// Add body elements
			for _, elem := range scene.Body {
				p := convertElementToParagraph(elem, charMap, lang)
				if p != nil {
					paragraphs = append(paragraphs, *p)
				}
			}
		}
	}

	fdx.Content = &fdxmodel.Content{
		Paragraphs: paragraphs,
	}

	// Add title page
	title := doc.Title.GetOrDefault(lang)
	if title != "" {
		fdx.TitlePage = &fdxmodel.TitlePage{
			Content: []fdxmodel.TitlePageContent{
				{
					Paragraphs: []fdxmodel.Paragraph{
						{
							Type: "Title",
							Text: []fdxmodel.Text{{Content: title}},
						},
					},
				},
			},
		}
	}

	return fdx
}

// formatSlugline formats a slugline for FDX output.
func formatSlugline(slug *model.Slugline) string {
	var parts []string

	// Context
	context := slug.Context
	if context == "" {
		context = "INT"
	}
	parts = append(parts, context+".")

	// Setting
	if slug.Setting != "" {
		parts = append(parts, slug.Setting)
	}

	// Time
	if slug.Time != "" {
		parts = append(parts, "-", slug.Time)
	}

	return strings.Join(parts, " ")
}

// convertElementToParagraph converts a ScreenJSON element to an FDX paragraph.
func convertElementToParagraph(elem model.Element, charMap map[string]string, lang string) *fdxmodel.Paragraph {
	p := &fdxmodel.Paragraph{}

	switch elem.Type {
	case model.ElementAction:
		p.Type = "Action"
		p.Text = []fdxmodel.Text{{Content: elem.Text.GetOrDefault(lang)}}

	case model.ElementCharacter:
		p.Type = "Character"
		display := elem.Display
		if display == "" {
			if name, ok := charMap[elem.Character]; ok {
				display = strings.ToUpper(name)
			}
		}
		p.Text = []fdxmodel.Text{{Content: display}}

	case model.ElementDialogue:
		p.Type = "Dialogue"
		p.Text = []fdxmodel.Text{{Content: elem.Text.GetOrDefault(lang)}}

	case model.ElementParenthetical:
		p.Type = "Parenthetical"
		text := elem.Text.GetOrDefault(lang)
		// Ensure parentheses
		if !strings.HasPrefix(text, "(") {
			text = "(" + text
		}
		if !strings.HasSuffix(text, ")") {
			text = text + ")"
		}
		p.Text = []fdxmodel.Text{{Content: text}}

	case model.ElementTransition:
		p.Type = "Transition"
		p.Text = []fdxmodel.Text{{Content: elem.Text.GetOrDefault(lang)}}

	case model.ElementShot:
		p.Type = "Shot"
		p.Text = []fdxmodel.Text{{Content: elem.Text.GetOrDefault(lang)}}

	case model.ElementGeneral:
		p.Type = "General"
		p.Text = []fdxmodel.Text{{Content: elem.Text.GetOrDefault(lang)}}

	default:
		return nil
	}

	return p
}
