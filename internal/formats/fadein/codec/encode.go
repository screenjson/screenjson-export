package codec

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/xml"
	"fmt"

	osfmodel "screenjson/export/internal/formats/fadein/model"
)

// Encoder encodes to FadeIn files.
type Encoder struct{}

// NewEncoder creates a new FadeIn encoder.
func NewEncoder() *Encoder {
	return &Encoder{}
}

// Encode serializes the OSF model to a FadeIn ZIP file.
func (e *Encoder) Encode(ctx context.Context, doc *osfmodel.Document) ([]byte, error) {
	if doc == nil {
		return nil, fmt.Errorf("nil OSF document")
	}
	// Marshal to XML
	xmlData, err := xml.MarshalIndent(doc, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("OSF XML encode error: %w", err)
	}

	// Add XML declaration
	xmlData = append([]byte(xml.Header), xmlData...)

	// Create ZIP archive
	var buf bytes.Buffer
	w := zip.NewWriter(&buf)

	// Add document.xml
	f, err := w.Create("document.xml")
	if err != nil {
		return nil, fmt.Errorf("failed to create document.xml in ZIP: %w", err)
	}

	if _, err := f.Write(xmlData); err != nil {
		return nil, fmt.Errorf("failed to write document.xml to ZIP: %w", err)
	}

	if err := w.Close(); err != nil {
		return nil, fmt.Errorf("failed to close ZIP: %w", err)
	}

	return buf.Bytes(), nil
}

// EncodeXML serializes just the OSF XML (without ZIP wrapper).
func (e *Encoder) EncodeXML(ctx context.Context, doc *osfmodel.Document) ([]byte, error) {
	xmlData, err := xml.MarshalIndent(doc, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("OSF XML encode error: %w", err)
	}
	return append([]byte(xml.Header), xmlData...), nil
}
