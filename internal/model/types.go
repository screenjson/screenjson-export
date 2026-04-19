package model

// Text represents translatable text keyed by BCP 47 language tag.
// Example: {"en": "Hello", "fr": "Bonjour"}
type Text map[string]string

// Meta represents arbitrary key/value metadata.
type Meta map[string]string

// Get returns the text for the given language, or empty string if not found.
func (t Text) Get(lang string) string {
	if t == nil {
		return ""
	}
	return t[lang]
}

// GetOrDefault returns the text for the given language, or the first available if not found.
func (t Text) GetOrDefault(lang string) string {
	if t == nil {
		return ""
	}
	if v, ok := t[lang]; ok {
		return v
	}
	// Return first available
	for _, v := range t {
		return v
	}
	return ""
}

// Set sets the text for the given language.
func (t Text) Set(lang, value string) {
	if t == nil {
		return
	}
	t[lang] = value
}

// Languages returns all language keys.
func (t Text) Languages() []string {
	if t == nil {
		return nil
	}
	langs := make([]string, 0, len(t))
	for k := range t {
		langs = append(langs, k)
	}
	return langs
}

// NewText creates a new Text with the given language and value.
func NewText(lang, value string) Text {
	return Text{lang: value}
}

// NewEnglishText creates a new Text with English content.
func NewEnglishText(value string) Text {
	return Text{"en": value}
}
