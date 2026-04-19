// Package codec provides Fountain encoding/decoding.
package codec

import (
	"bufio"
	"context"
	"fmt"
	"regexp"
	"strings"

	ftnmodel "screenjson/export/internal/formats/fountain/model"
)

var (
	sceneHeadingPattern = regexp.MustCompile(`^(?i)(INT|EXT|EST|INT\./EXT|INT/EXT|I/E)[\.\s]`)
	transitionPattern   = regexp.MustCompile(`^[A-Z\s]+TO:$`)
	characterPattern    = regexp.MustCompile(`^[A-Z][A-Z0-9\s\-'\.]+(\s*\([^)]+\))?$`)
	parentheticalPattern = regexp.MustCompile(`^\([^)]+\)$`)
	forcedCharPattern   = regexp.MustCompile(`^@(.+)$`)
	forcedScenePattern  = regexp.MustCompile(`^\.(.+)$`)
	forcedActionPattern = regexp.MustCompile(`^!(.+)$`)
	forcedTransPattern  = regexp.MustCompile(`^>(.+)$`)
	centeredPattern     = regexp.MustCompile(`^>(.+)<$`)
	sectionPattern      = regexp.MustCompile(`^(#{1,6})\s*(.+)$`)
	synopsisPattern     = regexp.MustCompile(`^=\s*(.+)$`)
	notePattern         = regexp.MustCompile(`\[\[([^\]]+)\]\]`)
	boneyardStart       = regexp.MustCompile(`/\*`)
	boneyardEnd         = regexp.MustCompile(`\*/`)
	pageBreakPattern    = regexp.MustCompile(`^===+$`)
	titleKeyPattern     = regexp.MustCompile(`^([A-Za-z\s]+):\s*(.*)$`)
	dualDialoguePattern = regexp.MustCompile(`\s*\^$`)
)

// Decoder decodes Fountain files.
type Decoder struct{}

// NewDecoder creates a new Fountain decoder.
func NewDecoder() *Decoder {
	return &Decoder{}
}

// Decode parses Fountain text into the Fountain model.
func (d *Decoder) Decode(ctx context.Context, data []byte) (*ftnmodel.Document, error) {
	doc := &ftnmodel.Document{}
	
	if len(data) == 0 {
		return nil, fmt.Errorf("empty Fountain input")
	}
	
	text := string(data)
	lines := strings.Split(text, "\n")
	
	// Parse title page if present
	titlePage, contentStart := parseTitlePage(lines)
	doc.TitlePage = titlePage
	
	// Parse content
	doc.Elements = parseContent(lines[contentStart:])
	
	return doc, nil
}

// parseTitlePage extracts title page metadata.
func parseTitlePage(lines []string) (*ftnmodel.TitlePage, int) {
	tp := &ftnmodel.TitlePage{
		Custom: make(map[string]string),
	}
	
	// Title page ends at first blank line followed by content
	// or after no title page markers are found
	
	i := 0
	inTitlePage := false
	var currentKey string
	var currentValue strings.Builder
	
	for i < len(lines) {
		line := lines[i]
		trimmed := strings.TrimSpace(line)
		
		// Check for title page key
		if match := titleKeyPattern.FindStringSubmatch(line); match != nil {
			if currentKey != "" {
				setTitlePageValue(tp, currentKey, strings.TrimSpace(currentValue.String()))
			}
			currentKey = strings.TrimSpace(match[1])
			currentValue.Reset()
			currentValue.WriteString(match[2])
			inTitlePage = true
		} else if inTitlePage && trimmed == "" {
			// End of title page
			if currentKey != "" {
				setTitlePageValue(tp, currentKey, strings.TrimSpace(currentValue.String()))
			}
			i++
			break
		} else if inTitlePage && (strings.HasPrefix(line, "   ") || strings.HasPrefix(line, "\t")) {
			// Continuation of multi-line value
			currentValue.WriteString("\n")
			currentValue.WriteString(strings.TrimSpace(line))
		} else if !inTitlePage && trimmed != "" {
			// No title page, content starts here
			break
		}
		
		i++
	}
	
	if currentKey != "" {
		setTitlePageValue(tp, currentKey, strings.TrimSpace(currentValue.String()))
	}
	
	if !inTitlePage {
		return nil, 0
	}
	
	return tp, i
}

// setTitlePageValue sets a title page field by key.
func setTitlePageValue(tp *ftnmodel.TitlePage, key, value string) {
	switch strings.ToLower(key) {
	case "title":
		tp.Title = value
	case "credit":
		tp.Credit = value
	case "author":
		tp.Author = value
	case "authors":
		tp.Authors = value
	case "source":
		tp.Source = value
	case "draft date":
		tp.DraftDate = value
	case "contact":
		tp.Contact = value
	case "copyright":
		tp.Copyright = value
	case "notes":
		tp.Notes = value
	default:
		tp.Custom[key] = value
	}
}

// parseContent parses the main screenplay content.
func parseContent(lines []string) []ftnmodel.Element {
	var elements []ftnmodel.Element
	
	scanner := &lineScanner{lines: lines, pos: 0}
	var lastCharacter bool
	inBoneyard := false
	
	for scanner.hasMore() {
		line := scanner.next()
		trimmed := strings.TrimSpace(line)
		
		// Handle boneyard (comments)
		if boneyardStart.MatchString(line) {
			inBoneyard = true
			continue
		}
		if inBoneyard {
			if boneyardEnd.MatchString(line) {
				inBoneyard = false
			}
			continue
		}
		
		// Blank line
		if trimmed == "" {
			lastCharacter = false
			elements = append(elements, ftnmodel.Element{Type: ftnmodel.ElementBlank})
			continue
		}
		
		// Page break
		if pageBreakPattern.MatchString(trimmed) {
			elements = append(elements, ftnmodel.Element{Type: ftnmodel.ElementPageBreak})
			continue
		}
		
		// Section heading
		if match := sectionPattern.FindStringSubmatch(trimmed); match != nil {
			elements = append(elements, ftnmodel.Element{
				Type:  ftnmodel.ElementSection,
				Text:  match[2],
				Depth: len(match[1]),
			})
			continue
		}
		
		// Synopsis
		if match := synopsisPattern.FindStringSubmatch(trimmed); match != nil {
			elements = append(elements, ftnmodel.Element{
				Type: ftnmodel.ElementSynopsis,
				Text: match[1],
			})
			continue
		}
		
		// Centered text
		if match := centeredPattern.FindStringSubmatch(trimmed); match != nil {
			elements = append(elements, ftnmodel.Element{
				Type:   ftnmodel.ElementCentered,
				Text:   strings.TrimSpace(match[1]),
				Forced: true,
			})
			continue
		}
		
		// Forced scene heading
		if match := forcedScenePattern.FindStringSubmatch(trimmed); match != nil {
			elements = append(elements, ftnmodel.Element{
				Type:   ftnmodel.ElementSceneHeading,
				Text:   match[1],
				Forced: true,
			})
			lastCharacter = false
			continue
		}
		
		// Scene heading
		if sceneHeadingPattern.MatchString(trimmed) {
			elements = append(elements, ftnmodel.Element{
				Type: ftnmodel.ElementSceneHeading,
				Text: trimmed,
			})
			lastCharacter = false
			continue
		}
		
		// Forced transition
		if match := forcedTransPattern.FindStringSubmatch(trimmed); match != nil {
			if !centeredPattern.MatchString(trimmed) {
				elements = append(elements, ftnmodel.Element{
					Type:   ftnmodel.ElementTransition,
					Text:   strings.TrimSpace(match[1]),
					Forced: true,
				})
				continue
			}
		}
		
		// Transition
		if transitionPattern.MatchString(trimmed) {
			elements = append(elements, ftnmodel.Element{
				Type: ftnmodel.ElementTransition,
				Text: trimmed,
			})
			continue
		}
		
		// Parenthetical (only after character or dialogue)
		if lastCharacter && parentheticalPattern.MatchString(trimmed) {
			elements = append(elements, ftnmodel.Element{
				Type: ftnmodel.ElementParenthetical,
				Text: trimmed,
			})
			continue
		}
		
		// Forced character
		if match := forcedCharPattern.FindStringSubmatch(trimmed); match != nil {
			dual := dualDialoguePattern.MatchString(match[1])
			name := dualDialoguePattern.ReplaceAllString(match[1], "")
			elements = append(elements, ftnmodel.Element{
				Type:   ftnmodel.ElementCharacter,
				Text:   name,
				Dual:   dual,
				Forced: true,
			})
			lastCharacter = true
			continue
		}
		
		// Character (ALL CAPS followed by dialogue)
		if !lastCharacter && characterPattern.MatchString(trimmed) && scanner.hasMore() {
			nextLine := strings.TrimSpace(scanner.peek())
			if nextLine != "" && !sceneHeadingPattern.MatchString(nextLine) {
				dual := dualDialoguePattern.MatchString(trimmed)
				name := dualDialoguePattern.ReplaceAllString(trimmed, "")
				elements = append(elements, ftnmodel.Element{
					Type: ftnmodel.ElementCharacter,
					Text: name,
					Dual: dual,
				})
				lastCharacter = true
				continue
			}
		}
		
		// Dialogue (after character)
		if lastCharacter {
			elements = append(elements, ftnmodel.Element{
				Type: ftnmodel.ElementDialogue,
				Text: trimmed,
			})
			continue
		}
		
		// Forced action
		if match := forcedActionPattern.FindStringSubmatch(trimmed); match != nil {
			elements = append(elements, ftnmodel.Element{
				Type:   ftnmodel.ElementAction,
				Text:   match[1],
				Forced: true,
			})
			continue
		}
		
		// Default: action
		elements = append(elements, ftnmodel.Element{
			Type: ftnmodel.ElementAction,
			Text: trimmed,
		})
		lastCharacter = false
	}
	
	return elements
}

// lineScanner provides a simple scanner for lines.
type lineScanner struct {
	lines []string
	pos   int
}

func (s *lineScanner) hasMore() bool {
	return s.pos < len(s.lines)
}

func (s *lineScanner) next() string {
	if s.pos >= len(s.lines) {
		return ""
	}
	line := s.lines[s.pos]
	s.pos++
	return line
}

func (s *lineScanner) peek() string {
	if s.pos >= len(s.lines) {
		return ""
	}
	return s.lines[s.pos]
}

// ReadFountain is a convenience function to parse Fountain from a bufio.Scanner.
func ReadFountain(scanner *bufio.Scanner) (*ftnmodel.Document, error) {
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	
	d := NewDecoder()
	return d.Decode(context.Background(), []byte(strings.Join(lines, "\n")))
}
