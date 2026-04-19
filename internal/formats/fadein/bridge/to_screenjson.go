// Package bridge provides conversion between FadeIn/OSF and ScreenJSON models.
package bridge

import (
	"strings"

	"github.com/google/uuid"

	osfmodel "screenjson/export/internal/formats/fadein/model"
	"screenjson/export/internal/model"
)

// ToScreenJSON converts an OSF document to a ScreenJSON document.
func ToScreenJSON(osf *osfmodel.Document, lang string) *model.Document {
	if lang == "" {
		lang = "en"
	}

	authorID := uuid.New().String()

	doc := &model.Document{
		ID:      uuid.New().String(),
		Version: "1.0.0",
		Generator: &model.Generator{
			Name:    "screenjson-export",
			Version: "1.0.0",
		},
		Title:   model.Text{lang: "Untitled"},
		Lang:    lang,
		Charset: "utf-8",
		Dir:     "ltr",
		Authors: []model.Author{
			{
				ID:     authorID,
				Given:  "Unknown",
				Family: "Author",
			},
		},
	}

	// Extract title from header if available
	if osf.Settings != nil && osf.Settings.PageHeader != "" {
		// Header often contains title
		header := osf.Settings.PageHeader
		// Try to extract title from header format like "#. TITLE by AUTHOR"
		if idx := strings.Index(header, " by "); idx > 0 {
			title := strings.TrimPrefix(header[:idx], "#. ")
			doc.Title = model.Text{lang: title}
		}
	}

	// Extract characters
	characters := extractCharacters(osf)
	doc.Characters = characters

	// Build character lookup map
	charMap := make(map[string]string) // name -> UUID
	for _, c := range characters {
		charMap[strings.ToUpper(c.Name)] = c.ID
	}

	// Convert content
	doc.Content = convertContent(osf, authorID, charMap, lang)

	return doc
}

// extractCharacters extracts unique character names from the document.
func extractCharacters(osf *osfmodel.Document) []model.Character {
	seen := make(map[string]bool)
	var chars []model.Character

	if osf.Paragraphs == nil {
		return chars
	}

	for _, p := range osf.Paragraphs.Paragraphs {
		if p.Style != nil && p.Style.BaseStyle == osfmodel.StyleCharacter {
			name := strings.TrimSpace(p.Text)
			name = cleanCharacterName(name)
			upperName := strings.ToUpper(name)

			if name != "" && !seen[upperName] {
				seen[upperName] = true
				chars = append(chars, model.Character{
					ID:   uuid.New().String(),
					Name: name,
				})
			}
		}
	}

	return chars
}

// cleanCharacterName removes extensions from character names.
func cleanCharacterName(name string) string {
	name = strings.TrimSuffix(name, "(V.O.)")
	name = strings.TrimSuffix(name, "(V.O)")
	name = strings.TrimSuffix(name, "(O.S.)")
	name = strings.TrimSuffix(name, "(O.S)")
	name = strings.TrimSuffix(name, "(O.C.)")
	name = strings.TrimSuffix(name, "(O.C)")
	name = strings.TrimSuffix(name, "(CONT'D)")
	name = strings.TrimSuffix(name, "(CONT)")
	name = strings.TrimSuffix(name, "(cont'd)")
	name = strings.TrimSpace(name)
	return name
}

// convertContent converts OSF content to ScreenJSON content.
func convertContent(osf *osfmodel.Document, authorID string, charMap map[string]string, lang string) *model.Content {
	content := &model.Content{
		Cover: &model.Cover{
			Title:   model.Text{lang: "Untitled"},
			Authors: []string{authorID},
		},
	}

	if osf.Paragraphs == nil {
		return content
	}

	var scenes []model.Scene
	var currentScene *model.Scene
	var sceneNumber int

	for _, p := range osf.Paragraphs.Paragraphs {
		if p.Style == nil {
			continue
		}

		style := p.Style.BaseStyle

		switch style {
		case osfmodel.StyleSceneHeading:
			// Start a new scene
			if currentScene != nil {
				scenes = append(scenes, *currentScene)
			}
			sceneNumber++
			currentScene = createScene(p.Text, authorID, sceneNumber, lang)

		default:
			// Add element to current scene
			if currentScene == nil {
				// Create a default scene if none exists
				sceneNumber++
				currentScene = &model.Scene{
					ID:      uuid.New().String(),
					Authors: []string{authorID},
					Heading: &model.Slugline{
						No:      sceneNumber,
						Context: "INT",
						Setting: "UNKNOWN",
						Time:    "DAY",
					},
					Body: []model.Element{},
				}
			}

			element := convertParaToElement(p, authorID, charMap, lang)
			if element != nil {
				currentScene.Body = append(currentScene.Body, *element)

				// Track characters in scene
				if element.Character != "" {
					found := false
					for _, c := range currentScene.Cast {
						if c == element.Character {
							found = true
							break
						}
					}
					if !found {
						currentScene.Cast = append(currentScene.Cast, element.Character)
					}
				}
			}
		}
	}

	// Don't forget the last scene
	if currentScene != nil {
		scenes = append(scenes, *currentScene)
	}

	content.Scenes = scenes
	return content
}

// createScene creates a new scene from a scene heading.
func createScene(text string, authorID string, sceneNumber int, lang string) *model.Scene {
	scene := &model.Scene{
		ID:      uuid.New().String(),
		Authors: []string{authorID},
		Body:    []model.Element{},
	}

	scene.Heading = parseSlugline(text, sceneNumber)
	return scene
}

// parseSlugline parses a scene heading string into a Slugline.
func parseSlugline(text string, sceneNumber int) *model.Slugline {
	slug := &model.Slugline{
		No:   sceneNumber,
		Time: "DAY",
	}

	text = strings.TrimSpace(text)
	upper := strings.ToUpper(text)

	// Detect context
	if strings.HasPrefix(upper, "INT/EXT") || strings.HasPrefix(upper, "I/E") {
		slug.Context = "INT/EXT"
		text = strings.TrimPrefix(upper, "INT/EXT")
		text = strings.TrimPrefix(text, "I/E")
	} else if strings.HasPrefix(upper, "EXT/INT") {
		slug.Context = "EXT/INT"
		text = strings.TrimPrefix(upper, "EXT/INT")
	} else if strings.HasPrefix(upper, "INT.") || strings.HasPrefix(upper, "INT ") {
		slug.Context = "INT"
		text = strings.TrimPrefix(upper, "INT.")
		text = strings.TrimPrefix(text, "INT ")
	} else if strings.HasPrefix(upper, "EXT.") || strings.HasPrefix(upper, "EXT ") {
		slug.Context = "EXT"
		text = strings.TrimPrefix(upper, "EXT.")
		text = strings.TrimPrefix(text, "EXT ")
	} else {
		slug.Context = "INT"
	}

	text = strings.TrimSpace(text)
	text = strings.TrimPrefix(text, ".")
	text = strings.TrimPrefix(text, "-")
	text = strings.TrimSpace(text)

	// Split setting and time
	parts := strings.Split(text, " - ")
	if len(parts) >= 2 {
		slug.Setting = strings.TrimSpace(parts[0])
		slug.Time = strings.TrimSpace(parts[len(parts)-1])
	} else {
		parts = strings.Split(text, " – ")
		if len(parts) >= 2 {
			slug.Setting = strings.TrimSpace(parts[0])
			slug.Time = strings.TrimSpace(parts[len(parts)-1])
		} else {
			slug.Setting = text
		}
	}

	return slug
}

// convertParaToElement converts an OSF paragraph to a ScreenJSON element.
func convertParaToElement(p osfmodel.Para, authorID string, charMap map[string]string, lang string) *model.Element {
	if p.Style == nil {
		return nil
	}

	style := p.Style.BaseStyle
	text := p.Text

	switch style {
	case osfmodel.StyleAction:
		return &model.Element{
			ID:      uuid.New().String(),
			Type:    model.ElementAction,
			Authors: []string{authorID},
			Text:    model.Text{lang: text},
		}

	case osfmodel.StyleCharacter:
		name := cleanCharacterName(strings.TrimSpace(text))
		charID := charMap[strings.ToUpper(name)]
		return &model.Element{
			ID:        uuid.New().String(),
			Type:      model.ElementCharacter,
			Authors:   []string{authorID},
			Character: charID,
			Display:   text,
		}

	case osfmodel.StyleDialogue:
		return &model.Element{
			ID:      uuid.New().String(),
			Type:    model.ElementDialogue,
			Authors: []string{authorID},
			Text:    model.Text{lang: text},
		}

	case osfmodel.StyleParenthetical:
		return &model.Element{
			ID:      uuid.New().String(),
			Type:    model.ElementParenthetical,
			Authors: []string{authorID},
			Text:    model.Text{lang: text},
		}

	case osfmodel.StyleTransition:
		return &model.Element{
			ID:      uuid.New().String(),
			Type:    model.ElementTransition,
			Authors: []string{authorID},
			Text:    model.Text{lang: text},
		}

	case osfmodel.StyleShot:
		return &model.Element{
			ID:      uuid.New().String(),
			Type:    model.ElementShot,
			Authors: []string{authorID},
			Text:    model.Text{lang: text},
		}

	default:
		if text != "" {
			return &model.Element{
				ID:      uuid.New().String(),
				Type:    model.ElementGeneral,
				Authors: []string{authorID},
				Text:    model.Text{lang: text},
			}
		}
	}

	return nil
}
