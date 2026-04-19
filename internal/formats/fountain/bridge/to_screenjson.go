// Package bridge provides conversion between Fountain and ScreenJSON models.
package bridge

import (
	"strings"

	"github.com/google/uuid"

	ftnmodel "screenjson/export/internal/formats/fountain/model"
	"screenjson/export/internal/model"
)

// ToScreenJSON converts a Fountain document to a ScreenJSON document.
func ToScreenJSON(ftn *ftnmodel.Document, lang string) *model.Document {
	if lang == "" {
		lang = "en"
	}

	authorID := uuid.New().String()
	title := "Untitled"
	authorName := "Unknown"
	authorFamily := "Author"

	// Extract metadata from title page
	if ftn.TitlePage != nil {
		if ftn.TitlePage.Title != "" {
			title = ftn.TitlePage.Title
		}
		if ftn.TitlePage.Author != "" {
			parts := strings.Fields(ftn.TitlePage.Author)
			if len(parts) >= 2 {
				authorName = parts[0]
				authorFamily = strings.Join(parts[1:], " ")
			} else if len(parts) == 1 {
				authorName = parts[0]
				authorFamily = ""
			}
		}
	}

	doc := &model.Document{
		ID:      uuid.New().String(),
		Version: "1.0.0",
		Generator: &model.Generator{
			Name:    "screenjson-export",
			Version: "1.0.0",
		},
		Title:   model.Text{lang: title},
		Lang:    lang,
		Charset: "utf-8",
		Dir:     "ltr",
		Authors: []model.Author{
			{
				ID:     authorID,
				Given:  authorName,
				Family: authorFamily,
			},
		},
	}

	// Extract characters
	characters := extractCharacters(ftn)
	doc.Characters = characters

	// Build character lookup
	charMap := make(map[string]string) // name -> UUID
	for _, c := range characters {
		charMap[strings.ToUpper(c.Name)] = c.ID
	}

	// Convert content
	doc.Content = convertContent(ftn, authorID, charMap, lang)

	return doc
}

// extractCharacters extracts unique character names from the document.
func extractCharacters(ftn *ftnmodel.Document) []model.Character {
	seen := make(map[string]bool)
	var chars []model.Character

	for _, elem := range ftn.Elements {
		if elem.Type == ftnmodel.ElementCharacter {
			name := cleanCharacterName(elem.Text)
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
	// Remove parenthetical extensions
	if idx := strings.Index(name, "("); idx > 0 {
		name = name[:idx]
	}
	return strings.TrimSpace(name)
}

// convertContent converts Fountain content to ScreenJSON content.
func convertContent(ftn *ftnmodel.Document, authorID string, charMap map[string]string, lang string) *model.Content {
	title := "Untitled"
	if ftn.TitlePage != nil && ftn.TitlePage.Title != "" {
		title = ftn.TitlePage.Title
	}

	content := &model.Content{
		Cover: &model.Cover{
			Title:   model.Text{lang: title},
			Authors: []string{authorID},
		},
	}

	var scenes []model.Scene
	var currentScene *model.Scene
	var sceneNumber int

	for _, elem := range ftn.Elements {
		switch elem.Type {
		case ftnmodel.ElementSceneHeading:
			// Start new scene
			if currentScene != nil {
				scenes = append(scenes, *currentScene)
			}
			sceneNumber++
			currentScene = &model.Scene{
				ID:      uuid.New().String(),
				Authors: []string{authorID},
				Heading: parseSlugline(elem.Text, sceneNumber),
				Body:    []model.Element{},
			}

		case ftnmodel.ElementAction:
			if currentScene == nil {
				sceneNumber++
				currentScene = createDefaultScene(authorID, sceneNumber)
			}
			currentScene.Body = append(currentScene.Body, model.Element{
				ID:      uuid.New().String(),
				Type:    model.ElementAction,
				Authors: []string{authorID},
				Text:    model.Text{lang: elem.Text},
			})

		case ftnmodel.ElementCharacter:
			if currentScene == nil {
				sceneNumber++
				currentScene = createDefaultScene(authorID, sceneNumber)
			}
			name := cleanCharacterName(elem.Text)
			charID := charMap[strings.ToUpper(name)]
			currentScene.Body = append(currentScene.Body, model.Element{
				ID:        uuid.New().String(),
				Type:      model.ElementCharacter,
				Authors:   []string{authorID},
				Character: charID,
				Display:   elem.Text,
			})
			// Track cast
			if charID != "" {
				found := false
				for _, c := range currentScene.Cast {
					if c == charID {
						found = true
						break
					}
				}
				if !found {
					currentScene.Cast = append(currentScene.Cast, charID)
				}
			}

		case ftnmodel.ElementDialogue:
			if currentScene == nil {
				sceneNumber++
				currentScene = createDefaultScene(authorID, sceneNumber)
			}
			dual := false
			// Check if last element was a dual dialogue character
			if len(currentScene.Body) > 0 {
				last := currentScene.Body[len(currentScene.Body)-1]
				if last.Type == model.ElementCharacter && last.Dual {
					dual = true
				}
			}
			currentScene.Body = append(currentScene.Body, model.Element{
				ID:      uuid.New().String(),
				Type:    model.ElementDialogue,
				Authors: []string{authorID},
				Text:    model.Text{lang: elem.Text},
				Dual:    dual,
			})

		case ftnmodel.ElementParenthetical:
			if currentScene == nil {
				sceneNumber++
				currentScene = createDefaultScene(authorID, sceneNumber)
			}
			currentScene.Body = append(currentScene.Body, model.Element{
				ID:      uuid.New().String(),
				Type:    model.ElementParenthetical,
				Authors: []string{authorID},
				Text:    model.Text{lang: elem.Text},
			})

		case ftnmodel.ElementTransition:
			if currentScene == nil {
				sceneNumber++
				currentScene = createDefaultScene(authorID, sceneNumber)
			}
			currentScene.Body = append(currentScene.Body, model.Element{
				ID:      uuid.New().String(),
				Type:    model.ElementTransition,
				Authors: []string{authorID},
				Text:    model.Text{lang: elem.Text},
			})

		case ftnmodel.ElementCentered, ftnmodel.ElementLyrics:
			if currentScene == nil {
				sceneNumber++
				currentScene = createDefaultScene(authorID, sceneNumber)
			}
			currentScene.Body = append(currentScene.Body, model.Element{
				ID:      uuid.New().String(),
				Type:    model.ElementGeneral,
				Authors: []string{authorID},
				Text:    model.Text{lang: elem.Text},
			})

		// Skip blank, section, synopsis, note, boneyard, page break
		}
	}

	if currentScene != nil {
		scenes = append(scenes, *currentScene)
	}

	content.Scenes = scenes
	return content
}

// createDefaultScene creates a default scene.
func createDefaultScene(authorID string, sceneNumber int) *model.Scene {
	return &model.Scene{
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

// parseSlugline parses a scene heading into a Slugline.
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
	} else if strings.HasPrefix(upper, "EST.") || strings.HasPrefix(upper, "EST ") {
		slug.Context = "EXT"
		text = strings.TrimPrefix(upper, "EST.")
		text = strings.TrimPrefix(text, "EST ")
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
