// Package model defines Open Screenplay Format (OSF) structures for FadeIn.
package model

import "encoding/xml"

// Document is the root element of an OSF document.
type Document struct {
	XMLName xml.Name `xml:"document"`
	Type    string   `xml:"type,attr"`
	Version string   `xml:"version,attr"`
	Info    *Info    `xml:"info"`
	Settings *Settings `xml:"settings"`
	FadeInSettings *FadeInSettings `xml:"fadein_settings"`
	Styles  *Styles  `xml:"styles"`
	Paragraphs *Paragraphs `xml:"paragraphs"`
}

// Info contains document metadata.
type Info struct {
	UUID      string `xml:"uuid,attr"`
	DraftUUID string `xml:"draft_uuid,attr"`
	PageCount string `xml:"pagecount,attr"`
}

// Settings contains document settings.
type Settings struct {
	PageWidth       string `xml:"page_width,attr"`
	PageHeight      string `xml:"page_height,attr"`
	MarginTop       string `xml:"margin_top,attr"`
	MarginBottom    string `xml:"margin_bottom,attr"`
	MarginLeft      string `xml:"margin_left,attr"`
	MarginRight     string `xml:"margin_right,attr"`
	NormalLinesPerInch string `xml:"normal_linesperinch,attr"`
	ElementSpacing  string `xml:"element_spacing,attr"`
	BreakOnSentences string `xml:"break_on_sentences,attr"`
	DialogueContinues string `xml:"dialogue_continues,attr"`
	DialoguePagebreaks string `xml:"dialogue_pagebreaks,attr"`
	ContText        string `xml:"cont_text,attr"`
	MoreText        string `xml:"more_text,attr"`
	ScenesContinue  string `xml:"scenes_continue,attr"`
	ContinuedText   string `xml:"continued_text,attr"`
	NumberContinued string `xml:"number_continued,attr"`
	SceneTimeSeparator string `xml:"scene_time_separator,attr"`
	PageHeader      string `xml:"page_header,attr"`
	PageFooter      string `xml:"page_footer,attr"`
	HeaderAlignment string `xml:"header_alignment,attr"`
	FooterAlignment string `xml:"footer_alignment,attr"`
	HeaderFirstPage string `xml:"header_first_page,attr"`
	FooterFirstPage string `xml:"footer_first_page,attr"`
	PagesLocked     string `xml:"pages_locked,attr"`
	PagenumberStart string `xml:"pagenumber_start,attr"`
	PagenumberMode  string `xml:"pagenumber_mode,attr"`
	Revision        string `xml:"revision,attr"`
	DocumentRevision string `xml:"document_revision,attr"`
	RevisionMode    string `xml:"revision_mode,attr"`
	ShowRevisions   string `xml:"show_revisions,attr"`
	SelectedRevisions string `xml:"selected_revisions,attr"`
}

// FadeInSettings contains FadeIn-specific settings.
type FadeInSettings struct {
	SavedWith           string `xml:"saved_with,attr"`
	LastPosition        string `xml:"last_position,attr"`
	IndexCardsShow      string `xml:"index_cards_show,attr"`
	IndexCardsTextSize  string `xml:"index_cards_text_size,attr"`
	IndexCardsUseFolders string `xml:"index_cards_use_folders,attr"`
	IndexCardsUseColors string `xml:"index_cards_use_colors,attr"`
	NavigatorShow       string `xml:"navigator_show,attr"`
	NavigatorPagecount  string `xml:"navigator_pagecount,attr"`
	Narrator            string `xml:"narrator,attr"`
	NarratorDefault     string `xml:"narrator_default,attr"`
}

// Styles contains style definitions.
type Styles struct {
	Styles      []Style `xml:"style"`
	HeaderStyle *Style  `xml:"header_style"`
	FooterStyle *Style  `xml:"footer_style"`
}

// Style defines formatting for an element type.
type Style struct {
	Name         string `xml:"name,attr"`
	Builtin      string `xml:"builtin,attr,omitempty"`
	BuiltinIndex string `xml:"builtin_index,attr,omitempty"`
	BaseStyle    string `xml:"basestyle,attr,omitempty"`
	StyleEnter   string `xml:"style_enter,attr,omitempty"`
	StyleTabBefore string `xml:"style_tab_before,attr,omitempty"`
	StyleTabAfter string `xml:"style_tab_after,attr,omitempty"`
	Font         string `xml:"font,attr,omitempty"`
	Size         string `xml:"size,attr,omitempty"`
	Bold         string `xml:"bold,attr,omitempty"`
	Italic       string `xml:"italic,attr,omitempty"`
	Underline    string `xml:"underline,attr,omitempty"`
	SpaceBefore  string `xml:"spacebefore,attr,omitempty"`
	SpaceAfter   string `xml:"spaceafter,attr,omitempty"`
	KeepWithNext string `xml:"keepwithnext,attr,omitempty"`
	AllCaps      string `xml:"allcaps,attr,omitempty"`
	LeftIndent   string `xml:"leftindent,attr,omitempty"`
	RightIndent  string `xml:"rightindent,attr,omitempty"`
	Align        string `xml:"align,attr,omitempty"`
}

// Paragraphs contains all paragraph elements.
type Paragraphs struct {
	Paragraphs []Para `xml:"para"`
}

// Para represents a paragraph element.
type Para struct {
	EditedBy string     `xml:"edited_by,attr,omitempty"`
	Style    *ParaStyle `xml:"style"`
	Text     string     `xml:"text"`
}

// ParaStyle references a base style.
type ParaStyle struct {
	BaseStyle string `xml:"basestyle,attr"`
}

// StyleName constants for OSF element types.
const (
	StyleSceneHeading   = "Scene Heading"
	StyleAction         = "Action"
	StyleCharacter      = "Character"
	StyleDialogue       = "Dialogue"
	StyleParenthetical  = "Parenthetical"
	StyleTransition     = "Transition"
	StyleShot           = "Shot"
)
