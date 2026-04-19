// Package codec provides FadeIn encoding/decoding.
package codec

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"

	osfmodel "screenjson/export/internal/formats/fadein/model"
)

// Decoder decodes FadeIn files.
type Decoder struct{}

const maxFadeInXMLSize = 50 << 20 // 50MB safety limit

// NewDecoder creates a new FadeIn decoder.
func NewDecoder() *Decoder {
	return &Decoder{}
}

// Decode parses a FadeIn ZIP file into the OSF model.
func (d *Decoder) Decode(ctx context.Context, data []byte) (*osfmodel.Document, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty FadeIn input")
	}
	if len(data) < 4 || data[0] != 'P' || data[1] != 'K' {
		return nil, fmt.Errorf("FadeIn input is not a ZIP archive")
	}
	// FadeIn files are ZIP archives containing document.xml
	r, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, fmt.Errorf("FadeIn ZIP open error: %w", err)
	}

	// Find document.xml
	var docFile *zip.File
	for _, f := range r.File {
		if f.Name == "document.xml" {
			docFile = f
			break
		}
	}

	if docFile == nil {
		return nil, fmt.Errorf("FadeIn file does not contain document.xml")
	}

	// Read document.xml
	rc, err := docFile.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open document.xml: %w", err)
	}
	defer rc.Close()

	limited := io.LimitReader(rc, maxFadeInXMLSize+1)
	xmlData, err := io.ReadAll(limited)
	if err != nil {
		return nil, fmt.Errorf("failed to read document.xml: %w", err)
	}
	if len(xmlData) > maxFadeInXMLSize {
		return nil, fmt.Errorf("document.xml exceeds size limit (%d bytes)", maxFadeInXMLSize)
	}

	// Parse XML
	var doc osfmodel.Document
	if err := xml.Unmarshal(xmlData, &doc); err != nil {
		return nil, fmt.Errorf("OSF XML parse error: %w", err)
	}

	return &doc, nil
}

// DecodeXML parses raw OSF XML (for testing or extracted files).
func (d *Decoder) DecodeXML(ctx context.Context, data []byte) (*osfmodel.Document, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty OSF XML input")
	}
	if len(data) > maxFadeInXMLSize {
		return nil, fmt.Errorf("OSF XML exceeds size limit (%d bytes)", maxFadeInXMLSize)
	}
	var doc osfmodel.Document
	if err := xml.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("OSF XML parse error: %w", err)
	}
	return &doc, nil
}
