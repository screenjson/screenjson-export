// Package model defines Final Draft XML structures.
package model

import "encoding/xml"

// FinalDraft is the root element of an FDX document.
type FinalDraft struct {
	XMLName      xml.Name  `xml:"FinalDraft"`
	DocumentType string    `xml:"DocumentType,attr"`
	Template     string    `xml:"Template,attr"`
	Version      string    `xml:"Version,attr"`
	Content      *Content  `xml:"Content"`
	TitlePage    *TitlePage `xml:"TitlePage"`
	ElementSettings []ElementSetting `xml:"ElementSettings>ElementSetting"`
	HeaderAndFooter *HeaderAndFooter `xml:"HeaderAndFooter"`
	SpellCheckIgnoreLists *SpellCheckIgnoreLists `xml:"SpellCheckIgnoreLists"`
	PageLayout   *PageLayout `xml:"PageLayout"`
	WindowState  *WindowState `xml:"WindowState"`
	TextState    *TextState `xml:"TextState"`
	Macros       *Macros `xml:"Macros"`
	Actors       *Actors `xml:"Actors"`
	Cast         *Cast `xml:"Cast"`
	SceneNumberOptions *SceneNumberOptions `xml:"SceneNumberOptions"`
	Revisions    *Revisions `xml:"Revisions"`
	SplitState   *SplitState `xml:"SplitState"`
	Watermarking *Watermarking `xml:"Watermarking"`
	MoresAndContinueds *MoresAndContinueds `xml:"MoresAndContinueds"`
	LockedPages  *LockedPages `xml:"LockedPages"`
	SmartType    *SmartType `xml:"SmartType"`
	TagData      *TagData `xml:"TagData"`
}

// Content holds the screenplay paragraphs.
type Content struct {
	Paragraphs []Paragraph `xml:"Paragraph"`
}

// Paragraph represents a screenplay element.
type Paragraph struct {
	Type            string          `xml:"Type,attr"`
	Number          string          `xml:"Number,attr,omitempty"`
	SceneProperties *SceneProperties `xml:"SceneProperties,omitempty"`
	Text            []Text          `xml:"Text"`
	DualDialogue    *DualDialogue   `xml:"DualDialogue,omitempty"`
}

// SceneProperties holds scene metadata.
type SceneProperties struct {
	Length string `xml:"Length,attr,omitempty"`
	Page   string `xml:"Page,attr,omitempty"`
	Title  string `xml:"Title,attr,omitempty"`
}

// Text represents formatted text content.
type Text struct {
	Style    string `xml:"Style,attr,omitempty"`
	Content  string `xml:",chardata"`
}

// DualDialogue marks dual dialogue.
type DualDialogue struct {
	Position string `xml:"Position,attr,omitempty"`
}

// TitlePage holds title page content.
type TitlePage struct {
	Content []TitlePageContent `xml:"Content"`
}

// TitlePageContent holds a title page section.
type TitlePageContent struct {
	Paragraphs []Paragraph `xml:"Paragraph"`
}

// ElementSetting defines formatting for an element type.
type ElementSetting struct {
	Type      string     `xml:"Type,attr"`
	FontSpec  *FontSpec  `xml:"FontSpec"`
	ParagraphSpec *ParagraphSpec `xml:"ParagraphSpec"`
	Behavior  *Behavior  `xml:"Behavior"`
}

// FontSpec defines font settings.
type FontSpec struct {
	AdornmentStyle string `xml:"AdornmentStyle,attr,omitempty"`
	Background     string `xml:"Background,attr,omitempty"`
	Color          string `xml:"Color,attr,omitempty"`
	Font           string `xml:"Font,attr,omitempty"`
	RevisionID     string `xml:"RevisionID,attr,omitempty"`
	Size           string `xml:"Size,attr,omitempty"`
	Style          string `xml:"Style,attr,omitempty"`
}

// ParagraphSpec defines paragraph settings.
type ParagraphSpec struct {
	Alignment     string `xml:"Alignment,attr,omitempty"`
	FirstIndent   string `xml:"FirstIndent,attr,omitempty"`
	Leading       string `xml:"Leading,attr,omitempty"`
	LeftIndent    string `xml:"LeftIndent,attr,omitempty"`
	RightIndent   string `xml:"RightIndent,attr,omitempty"`
	SpaceBefore   string `xml:"SpaceBefore,attr,omitempty"`
	SpaceAfter    string `xml:"SpaceAfter,attr,omitempty"`
	StartsNewPage string `xml:"StartsNewPage,attr,omitempty"`
}

// Behavior defines element behavior.
type Behavior struct {
	PaginateAs string `xml:"PaginateAs,attr,omitempty"`
	ReturnKey  string `xml:"ReturnKey,attr,omitempty"`
	Shortcut   string `xml:"Shortcut,attr,omitempty"`
}

// HeaderAndFooter holds header/footer settings.
type HeaderAndFooter struct {
	Header    *HeaderFooter `xml:"Header"`
	Footer    *HeaderFooter `xml:"Footer"`
	StartingPage string `xml:"StartingPage,attr,omitempty"`
}

// HeaderFooter holds header or footer content.
type HeaderFooter struct {
	Paragraphs []Paragraph `xml:"Paragraph"`
}

// PageLayout defines page layout settings.
type PageLayout struct {
	PageWidth     string `xml:"PageWidth,attr,omitempty"`
	PageHeight    string `xml:"PageHeight,attr,omitempty"`
	TopMargin     string `xml:"TopMargin,attr,omitempty"`
	BottomMargin  string `xml:"BottomMargin,attr,omitempty"`
	LeftMargin    string `xml:"LeftMargin,attr,omitempty"`
	RightMargin   string `xml:"RightMargin,attr,omitempty"`
}

// WindowState holds window state.
type WindowState struct {
	Mode       string `xml:"Mode,attr,omitempty"`
	Top        string `xml:"Top,attr,omitempty"`
	Left       string `xml:"Left,attr,omitempty"`
	Right      string `xml:"Right,attr,omitempty"`
	Bottom     string `xml:"Bottom,attr,omitempty"`
}

// TextState holds text state.
type TextState struct {
	Zoom       string `xml:"Zoom,attr,omitempty"`
	Selection  string `xml:"Selection,attr,omitempty"`
}

// Macros holds macro definitions.
type Macros struct {
	// Implementation can be extended
}

// Actors holds actor definitions.
type Actors struct {
	// Implementation can be extended
}

// Cast holds cast list.
type Cast struct {
	Members []CastMember `xml:"Member"`
}

// CastMember represents a cast member.
type CastMember struct {
	Name string `xml:",chardata"`
}

// SceneNumberOptions holds scene numbering settings.
type SceneNumberOptions struct {
	LeftLocation  string `xml:"LeftLocation,attr,omitempty"`
	RightLocation string `xml:"RightLocation,attr,omitempty"`
	ShowNumbers   string `xml:"ShowNumbers,attr,omitempty"`
}

// Revisions holds revision information.
type Revisions struct {
	ActiveSet string     `xml:"ActiveSet,attr,omitempty"`
	Revision  []Revision `xml:"Revision"`
}

// Revision represents a single revision.
type Revision struct {
	ID    string `xml:"ID,attr,omitempty"`
	Color string `xml:"Color,attr,omitempty"`
	Name  string `xml:"Name,attr,omitempty"`
	Mark  string `xml:"Mark,attr,omitempty"`
	Style string `xml:"Style,attr,omitempty"`
}

// SplitState holds split view state.
type SplitState struct {
	// Implementation can be extended
}

// Watermarking holds watermark settings.
type Watermarking struct {
	// Implementation can be extended
}

// MoresAndContinueds holds continuation settings.
type MoresAndContinueds struct {
	Mores      *Mores      `xml:"Mores"`
	Continueds *Continueds `xml:"Continueds"`
}

// Mores holds "more" settings.
type Mores struct {
	Show      string `xml:"Show,attr,omitempty"`
	Text      string `xml:"Text,attr,omitempty"`
}

// Continueds holds "continued" settings.
type Continueds struct {
	Show      string `xml:"Show,attr,omitempty"`
	Top       string `xml:"Top,attr,omitempty"`
	Bottom    string `xml:"Bottom,attr,omitempty"`
	DialogueTop string `xml:"DialogueTop,attr,omitempty"`
}

// LockedPages holds locked page information.
type LockedPages struct {
	// Implementation can be extended
}

// SmartType holds smart type settings.
type SmartType struct {
	// Implementation can be extended
}

// TagData holds tag information.
type TagData struct {
	// Implementation can be extended
}

// SpellCheckIgnoreLists holds spell check ignore lists.
type SpellCheckIgnoreLists struct {
	// Implementation can be extended
}
