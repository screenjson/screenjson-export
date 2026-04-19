package model

import "fmt"

// ParseError represents a parsing error that occurred during format conversion.
// Instead of aborting, errors are represented in the document to preserve structure.
type ParseError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Source  string `json:"source,omitempty"`  // Source format (fdx, fadein, etc.)
	Line    int    `json:"line,omitempty"`    // Line number in source
	Column  int    `json:"column,omitempty"`  // Column number in source
	Context string `json:"context,omitempty"` // Surrounding text for debugging
}

// Error implements the error interface.
func (e *ParseError) Error() string {
	if e.Line > 0 {
		return fmt.Sprintf("[%s] %s at line %d", e.Code, e.Message, e.Line)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Common error codes
const (
	ErrCodeInvalidXML        = "INVALID_XML"
	ErrCodeInvalidJSON       = "INVALID_JSON"
	ErrCodeInvalidYAML       = "INVALID_YAML"
	ErrCodeInvalidFormat     = "INVALID_FORMAT"
	ErrCodeMissingRequired   = "MISSING_REQUIRED"
	ErrCodeInvalidUUID       = "INVALID_UUID"
	ErrCodeInvalidElement    = "INVALID_ELEMENT"
	ErrCodeInvalidScene      = "INVALID_SCENE"
	ErrCodeInvalidCharacter  = "INVALID_CHARACTER"
	ErrCodeInvalidEncryption = "INVALID_ENCRYPTION"
	ErrCodeZipError          = "ZIP_ERROR"
	ErrCodeIOError           = "IO_ERROR"
	ErrCodePDFError          = "PDF_ERROR"
	ErrCodeRDFError          = "RDF_ERROR"
)

// NewParseError creates a new parse error.
func NewParseError(code, message string) *ParseError {
	return &ParseError{
		Code:    code,
		Message: message,
	}
}

// WithSource adds source format information to the error.
func (e *ParseError) WithSource(source string) *ParseError {
	e.Source = source
	return e
}

// WithLocation adds line/column information to the error.
func (e *ParseError) WithLocation(line, column int) *ParseError {
	e.Line = line
	e.Column = column
	return e
}

// WithContext adds surrounding context to the error.
func (e *ParseError) WithContext(context string) *ParseError {
	e.Context = context
	return e
}

// ValidationError represents a schema validation error.
type ValidationError struct {
	Path    string `json:"path"`
	Message string `json:"message"`
	Value   any    `json:"value,omitempty"`
}

// Error implements the error interface.
func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Path, e.Message)
}

// ValidationResult contains the result of schema validation.
type ValidationResult struct {
	Valid  bool               `json:"valid"`
	Errors []*ValidationError `json:"errors,omitempty"`
}

// NewValidationResult creates a new validation result.
func NewValidationResult(valid bool) *ValidationResult {
	return &ValidationResult{Valid: valid}
}

// AddError adds a validation error.
func (r *ValidationResult) AddError(path, message string, value any) {
	r.Valid = false
	r.Errors = append(r.Errors, &ValidationError{
		Path:    path,
		Message: message,
		Value:   value,
	})
}
