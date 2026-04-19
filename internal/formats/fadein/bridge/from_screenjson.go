package bridge

import (
	"fmt"
	"strings"

	"github.com/google/uuid"

	osfmodel "screenjson/export/internal/formats/fadein/model"
	"screenjson/export/internal/model"
)

// FromScreenJSON converts a ScreenJSON document to an OSF document.
func FromScreenJSON(doc *model.Document, lang string) *osfmodel.Document {
	if lang == "" {
		lang = "en"
	}

	osf := &osfmodel.Document{
		Type:    "Open Screenplay Format document",
		Version: "40",
	}

	// Set info
	osf.Info = &osfmodel.Info{
		UUID:      uuid.New().String(),
		DraftUUID: uuid.New().String(),
	}

	// Set default settings
	osf.Settings = &osfmodel.Settings{
		PageWidth:          "2159",
		PageHeight:         "2794",
		MarginTop:          "266",
		MarginBottom:       "228",
		MarginLeft:         "317",
		MarginRight:        "317",
		NormalLinesPerInch: "6.0",
		ElementSpacing:     "1.00",
		ContText:           "(cont'd)",
		MoreText:           "(MORE)",
		SceneTimeSeparator: " - ",
	}

	// Set default styles
	osf.Styles = &osfmodel.Styles{
		Styles: []osfmodel.Style{
			{Name: "Normal Text", Builtin: "1", BuiltinIndex: "0", Font: "Courier Prime", Size: "12"},
			{Name: "Scene Heading", Builtin: "1", BuiltinIndex: "1", BaseStyle: "Normal Text", Font: "Courier Prime", Size: "12", Bold: "1", SpaceBefore: "2.0", KeepWithNext: "1", AllCaps: "1"},
			{Name: "Action", Builtin: "1", BuiltinIndex: "2", BaseStyle: "Normal Text", Font: "Courier Prime", Size: "12", SpaceBefore: "1.0"},
			{Name: "Character", Builtin: "1", BuiltinIndex: "3", BaseStyle: "Normal Text", Font: "Courier Prime", Size: "12", SpaceBefore: "1.0", KeepWithNext: "1", LeftIndent: "635", AllCaps: "1"},
			{Name: "Parenthetical", Builtin: "1", BuiltinIndex: "4", BaseStyle: "Normal Text", Font: "Courier Prime", Size: "12", KeepWithNext: "1", LeftIndent: "508", RightIndent: "508"},
			{Name: "Dialogue", Builtin: "1", BuiltinIndex: "5", BaseStyle: "Normal Text", Font: "Courier Prime", Size: "12", LeftIndent: "330", RightIndent: "254"},
			{Name: "Transition", Builtin: "1", BuiltinIndex: "6", BaseStyle: "Normal Text", Font: "Courier Prime", Size: "12", SpaceBefore: "1.0", Align: "right", LeftIndent: "1016", RightIndent: "127", AllCaps: "1"},
			{Name: "Shot", Builtin: "1", BuiltinIndex: "7", BaseStyle: "Normal Text", Font: "Courier Prime", Size: "12", SpaceBefore: "1.0", KeepWithNext: "1", AllCaps: "1"},
		},
		HeaderStyle: &osfmodel.Style{BaseStyle: "Normal Text"},
		FooterStyle: &osfmodel.Style{BaseStyle: "Normal Text"},
	}

	// Build character lookup
	charMap := make(map[string]string) // UUID -> name
	for _, c := range doc.Characters {
		charMap[c.ID] = c.Name
	}

	// Convert content
	var paras []osfmodel.Para

	if doc.Content != nil {
		for _, scene := range doc.Content.Scenes {
			// Add scene heading
			if scene.Heading != nil {
				headingText := formatSlugline(scene.Heading)
				paras = append(paras, osfmodel.Para{
					Style: &osfmodel.ParaStyle{BaseStyle: osfmodel.StyleSceneHeading},
					Text:  headingText,
				})
			}

			// Add body elements
			for _, elem := range scene.Body {
				p := convertElementToPara(elem, charMap, lang)
				if p != nil {
					paras = append(paras, *p)
				}
			}
		}
	}

	osf.Paragraphs = &osfmodel.Paragraphs{
		Paragraphs: paras,
	}

	// Update page count estimate
	pageCount := len(paras) / 55 // Rough estimate
	if pageCount < 1 {
		pageCount = 1
	}
	osf.Info.PageCount = fmt.Sprintf("%d", pageCount)

	return osf
}

// formatSlugline formats a slugline for OSF output.
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

	result := strings.Join(parts, " ")

	if slug.Time != "" {
		result += " - " + slug.Time
	}

	return result
}

// convertElementToPara converts a ScreenJSON element to an OSF paragraph.
func convertElementToPara(elem model.Element, charMap map[string]string, lang string) *osfmodel.Para {
	p := &osfmodel.Para{}

	switch elem.Type {
	case model.ElementAction:
		p.Style = &osfmodel.ParaStyle{BaseStyle: osfmodel.StyleAction}
		p.Text = elem.Text.GetOrDefault(lang)

	case model.ElementCharacter:
		p.Style = &osfmodel.ParaStyle{BaseStyle: osfmodel.StyleCharacter}
		display := elem.Display
		if display == "" {
			if name, ok := charMap[elem.Character]; ok {
				display = strings.ToUpper(name)
			}
		}
		p.Text = display

	case model.ElementDialogue:
		p.Style = &osfmodel.ParaStyle{BaseStyle: osfmodel.StyleDialogue}
		p.Text = elem.Text.GetOrDefault(lang)

	case model.ElementParenthetical:
		p.Style = &osfmodel.ParaStyle{BaseStyle: osfmodel.StyleParenthetical}
		text := elem.Text.GetOrDefault(lang)
		if !strings.HasPrefix(text, "(") {
			text = "(" + text
		}
		if !strings.HasSuffix(text, ")") {
			text = text + ")"
		}
		p.Text = text

	case model.ElementTransition:
		p.Style = &osfmodel.ParaStyle{BaseStyle: osfmodel.StyleTransition}
		p.Text = elem.Text.GetOrDefault(lang)

	case model.ElementShot:
		p.Style = &osfmodel.ParaStyle{BaseStyle: osfmodel.StyleShot}
		p.Text = elem.Text.GetOrDefault(lang)

	case model.ElementGeneral:
		p.Style = &osfmodel.ParaStyle{BaseStyle: osfmodel.StyleAction}
		p.Text = elem.Text.GetOrDefault(lang)

	default:
		return nil
	}

	return p
}
