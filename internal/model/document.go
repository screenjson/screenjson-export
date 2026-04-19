// Package model defines the internal ScreenJSON document model.
// This package imports nothing and is the canonical validation/serialization boundary.
package model

import "time"

// Document represents a complete ScreenJSON screenplay document.
type Document struct {
	ID           string         `json:"id"`
	Version      string         `json:"version"`
	Generator    *Generator     `json:"generator,omitempty"`
	Title        Text           `json:"title"`
	Lang         string         `json:"lang"`
	Locale       string         `json:"locale,omitempty"`
	Charset      string         `json:"charset"`
	Dir          string         `json:"dir"`
	Authors      []Author       `json:"authors"`
	Contributors []Contributor  `json:"contributors,omitempty"`
	Characters   []Character    `json:"characters,omitempty"`
	Colors       []Color        `json:"colors,omitempty"`
	Sources      []Source       `json:"sources,omitempty"`
	Registrations []Registration `json:"registrations,omitempty"`
	Revisions    []Revision     `json:"revisions,omitempty"`
	Encrypt      *Encrypt       `json:"encrypt,omitempty"`
	License      *License       `json:"license,omitempty"`
	Taggable     []string       `json:"taggable,omitempty"`
	Genre        []string       `json:"genre,omitempty"`
	Themes       []string       `json:"themes,omitempty"`
	Logline      Text           `json:"logline,omitempty"`
	Content      *Content       `json:"document"`
	Analysis     *Analysis      `json:"analysis,omitempty"`
}

// Generator describes the tool that created this document.
type Generator struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Meta    Meta   `json:"meta,omitempty"`
}

// Author represents an original writer.
type Author struct {
	ID     string `json:"id"`
	Given  string `json:"given"`
	Family string `json:"family"`
	Meta   Meta   `json:"meta,omitempty"`
}

// Contributor represents any third party who contributed artistically.
type Contributor struct {
	ID     string   `json:"id"`
	Given  string   `json:"given"`
	Family string   `json:"family"`
	Roles  []string `json:"roles,omitempty"`
	Meta   Meta     `json:"meta,omitempty"`
}

// Source represents a source work this screenplay is based on.
type Source struct {
	ID    string `json:"id"`
	Type  string `json:"type"`
	Title Text   `json:"title"`
	Meta  Meta   `json:"meta,omitempty"`
}

// Registration represents WGA or other registration information.
type Registration struct {
	Authority string     `json:"authority"`
	ID        string     `json:"id"`
	Created   *time.Time `json:"created,omitempty"`
	Modified  *time.Time `json:"modified,omitempty"`
	Meta      Meta       `json:"meta,omitempty"`
}

// Encrypt describes encryption parameters.
type Encrypt struct {
	Cipher   string `json:"cipher"`
	Hash     string `json:"hash"`
	Encoding string `json:"encoding"`
	Meta     Meta   `json:"meta,omitempty"`
}

// License describes license or copyright information.
type License struct {
	ID   string `json:"id"`
	Ref  string `json:"ref,omitempty"`
	Meta Meta   `json:"meta,omitempty"`
}

// Color represents a reusable UI color.
type Color struct {
	ID    string `json:"id"`
	Title Text   `json:"title,omitempty"`
	RGB   [3]int `json:"rgb"`
	Hex   string `json:"hex"`
	Meta  Meta   `json:"meta,omitempty"`
}

// Content is the screenplay document container.
type Content struct {
	Cover     *Cover     `json:"cover"`
	Layout    *Layout    `json:"layout,omitempty"`
	Bookmarks []Bookmark `json:"bookmarks,omitempty"`
	Scenes    []Scene    `json:"scenes"`
	Meta      Meta       `json:"meta,omitempty"`
}

// Cover represents title page metadata.
type Cover struct {
	Title   Text     `json:"title"`
	Authors []string `json:"authors"` // UUIDs
	Sources []string `json:"sources,omitempty"` // UUIDs
	Extra   Text     `json:"extra,omitempty"`
	Meta    Meta     `json:"meta,omitempty"`
}

// Layout contains optional rendering rules.
type Layout struct {
	Header    *Ribbon    `json:"header,omitempty"`
	Footer    *Ribbon    `json:"footer,omitempty"`
	Status    *Status    `json:"status,omitempty"`
	Styles    []Style    `json:"styles,omitempty"`
	Templates []Template `json:"templates,omitempty"`
	Guides    []Format   `json:"guides,omitempty"`
	Meta      Meta       `json:"meta,omitempty"`
}

// Ribbon represents a header or footer ribbon.
type Ribbon struct {
	Cover bool   `json:"cover"`
	Show  bool   `json:"show"`
	Start int    `json:"start,omitempty"`
	Omit  []int  `json:"omit,omitempty"`
	Text  Text   `json:"text"`
	Meta  Meta   `json:"meta,omitempty"`
}

// Status represents script revision status.
type Status struct {
	Color   string     `json:"color"`
	Round   int        `json:"round"`
	Updated *time.Time `json:"updated"`
	Meta    Meta       `json:"meta,omitempty"`
}

// Style represents a reusable presentation style rule.
type Style struct {
	ID      string `json:"id"`
	Default bool   `json:"default"`
	Text    string `json:"text"`
	Meta    Meta   `json:"meta,omitempty"`
}

// Template represents a reusable presentation template.
type Template struct {
	ID      string `json:"id"`
	Default bool   `json:"default"`
	Text    string `json:"text"`
	Meta    Meta   `json:"meta,omitempty"`
}

// Format references an external screenplay presentation standard.
type Format struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	URL   string `json:"url"`
	Note  string `json:"note,omitempty"`
	Meta  Meta   `json:"meta,omitempty"`
}

// Bookmark represents a shortcut to a specific element.
type Bookmark struct {
	ID      string `json:"id"`
	Scene   string `json:"scene"`
	Element string `json:"element"`
	Title   Text   `json:"title"`
	Desc    Text   `json:"desc,omitempty"`
	Meta    Meta   `json:"meta,omitempty"`
}
