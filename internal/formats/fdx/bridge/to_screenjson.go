// Package bridge provides conversion between FDX and ScreenJSON models.
package bridge

import (
	"strings"

	"github.com/google/uuid"

	fdxmodel "screenjson/export/internal/formats/fdx/model"
	"screenjson/export/internal/model"
)

// ToScreenJSON converts an FDX document to a ScreenJSON document.
func ToScreenJSON(fdx *fdxmodel.FinalDraft, lang string) *model.Document {
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
		Title:   model.Text{lang: extractTitle(fdx)},
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

	// Extract characters from the document
	characters := extractCharacters(fdx)
	doc.Characters = characters

	// Build character lookup map
	charMap := make(map[string]string) // name -> UUID
	for _, c := range characters {
		charMap[strings.ToUpper(c.Name)] = c.ID
	}

	// Convert content
	doc.Content = convertContent(fdx, authorID, charMap, lang)

	return doc
}

// extractTitle attempts to extract the title from title page or first scene.
func extractTitle(fdx *fdxmodel.FinalDraft) string {
	if fdx.TitlePage != nil {
		for _, content := range fdx.TitlePage.Content {
			for _, p := range content.Paragraphs {
				if strings.ToLower(p.Type) == "title" {
					return getTextContent(p.Text)
				}
			}
		}
	}
	return "Untitled"
}

// extractCharacters extracts unique character names from the document.
func extractCharacters(fdx *fdxmodel.FinalDraft) []model.Character {
	seen := make(map[string]bool)
	var chars []model.Character

	if fdx.Content == nil {
		return chars
	}

	for _, p := range fdx.Content.Paragraphs {
		if strings.ToLower(p.Type) == "character" {
			name := strings.TrimSpace(getTextContent(p.Text))
			// Remove extensions like (V.O.), (O.S.), (CONT'D)
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
	// Remove common extensions
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

// convertContent converts FDX content to ScreenJSON content.
func convertContent(fdx *fdxmodel.FinalDraft, authorID string, charMap map[string]string, lang string) *model.Content {
	content := &model.Content{
		Cover: &model.Cover{
			Title:   model.Text{lang: extractTitle(fdx)},
			Authors: []string{authorID},
		},
	}

	if fdx.Content == nil {
		return content
	}

	var scenes []model.Scene
	var currentScene *model.Scene
	var sceneNumber int

	for _, p := range fdx.Content.Paragraphs {
		ptype := strings.ToLower(p.Type)

		switch ptype {
		case "scene heading":
			// Start a new scene
			if currentScene != nil {
				scenes = append(scenes, *currentScene)
			}
			sceneNumber++
			currentScene = createScene(p, authorID, sceneNumber, lang)

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

			element := convertParagraphToElement(p, authorID, charMap, lang)
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

// createScene creates a new scene from a scene heading paragraph.
func createScene(p fdxmodel.Paragraph, authorID string, sceneNumber int, lang string) *model.Scene {
	text := strings.TrimSpace(getTextContent(p.Text))
	
	scene := &model.Scene{
		ID:      uuid.New().String(),
		Authors: []string{authorID},
		Body:    []model.Element{},
	}

	// Parse scene heading
	scene.Heading = parseSlugline(text, sceneNumber)

	return scene
}

// parseSlugline parses a scene heading string into a Slugline.
func parseSlugline(text string, sceneNumber int) *model.Slugline {
	slug := &model.Slugline{
		No:   sceneNumber,
		Time: "DAY", // Default
	}

	text = strings.TrimSpace(text)
	upper := strings.ToUpper(text)

	// Detect context (INT/EXT)
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
		slug.Context = "INT" // Default
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
		parts = strings.Split(text, " – ") // em-dash
		if len(parts) >= 2 {
			slug.Setting = strings.TrimSpace(parts[0])
			slug.Time = strings.TrimSpace(parts[len(parts)-1])
		} else {
			slug.Setting = text
		}
	}

	return slug
}

// convertParagraphToElement converts an FDX paragraph to a ScreenJSON element.
func convertParagraphToElement(p fdxmodel.Paragraph, authorID string, charMap map[string]string, lang string) *model.Element {
	ptype := strings.ToLower(p.Type)
	text := getTextContent(p.Text)

	switch ptype {
	case "action":
		return &model.Element{
			ID:      uuid.New().String(),
			Type:    model.ElementAction,
			Authors: []string{authorID},
			Text:    model.Text{lang: text},
		}

	case "character":
		name := cleanCharacterName(strings.TrimSpace(text))
		charID := charMap[strings.ToUpper(name)]
		return &model.Element{
			ID:        uuid.New().String(),
			Type:      model.ElementCharacter,
			Authors:   []string{authorID},
			Character: charID,
			Display:   text,
		}

	case "dialogue":
		return &model.Element{
			ID:      uuid.New().String(),
			Type:    model.ElementDialogue,
			Authors: []string{authorID},
			Text:    model.Text{lang: text},
		}

	case "parenthetical":
		return &model.Element{
			ID:      uuid.New().String(),
			Type:    model.ElementParenthetical,
			Authors: []string{authorID},
			Text:    model.Text{lang: text},
		}

	case "transition":
		return &model.Element{
			ID:      uuid.New().String(),
			Type:    model.ElementTransition,
			Authors: []string{authorID},
			Text:    model.Text{lang: text},
		}

	case "shot":
		return &model.Element{
			ID:      uuid.New().String(),
			Type:    model.ElementShot,
			Authors: []string{authorID},
			Text:    model.Text{lang: text},
		}

	case "general":
		return &model.Element{
			ID:      uuid.New().String(),
			Type:    model.ElementGeneral,
			Authors: []string{authorID},
			Text:    model.Text{lang: text},
		}

	default:
		// Unknown type, treat as general
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

// getTextContent extracts text content from FDX Text elements.
func getTextContent(texts []fdxmodel.Text) string {
	var parts []string
	for _, t := range texts {
		if t.Content != "" {
			parts = append(parts, t.Content)
		}
	}
	return strings.Join(parts, "")
}
