package codec

import (
	"context"
	"encoding/xml"
	"fmt"

	fdxmodel "screenjson/export/internal/formats/fdx/model"
)

// Encoder encodes to FDX files.
type Encoder struct{}

// NewEncoder creates a new FDX encoder.
func NewEncoder() *Encoder {
	return &Encoder{}
}

// Encode serializes the FDX model to XML.
func (e *Encoder) Encode(ctx context.Context, fdx *fdxmodel.FinalDraft) ([]byte, error) {
	if fdx == nil {
		return nil, fmt.Errorf("nil FDX document")
	}
	data, err := xml.MarshalIndent(fdx, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("FDX XML encode error: %w", err)
	}

	// Add XML declaration
	result := append([]byte(xml.Header), data...)
	return result, nil
}
