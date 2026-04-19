// Package codec provides FDX encoding/decoding.
package codec

import (
	"context"
	"encoding/xml"
	"fmt"

	fdxmodel "screenjson/export/internal/formats/fdx/model"
)

// Decoder decodes FDX files.
type Decoder struct{}

// NewDecoder creates a new FDX decoder.
func NewDecoder() *Decoder {
	return &Decoder{}
}

// Decode parses FDX XML data into the FDX model.
func (d *Decoder) Decode(ctx context.Context, data []byte) (*fdxmodel.FinalDraft, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty FDX input")
	}
	var fdx fdxmodel.FinalDraft
	if err := xml.Unmarshal(data, &fdx); err != nil {
		return nil, fmt.Errorf("FDX XML parse error: %w", err)
	}
	return &fdx, nil
}
