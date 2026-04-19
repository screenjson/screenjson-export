package model

// Character represents a character entity referenced throughout the screenplay.
type Character struct {
	ID      string   `json:"id"`
	Slug    string   `json:"slug,omitempty"`
	Name    string   `json:"name"`
	Aliases []string `json:"aliases,omitempty"`
	Desc    Text     `json:"desc,omitempty"`
	Traits  []string `json:"traits,omitempty"`
	Meta    Meta     `json:"meta,omitempty"`
}

// NewCharacter creates a new character with the given ID and name.
func NewCharacter(id, name string) Character {
	return Character{
		ID:   id,
		Name: name,
	}
}

// WithSlug sets the character's slug.
func (c Character) WithSlug(slug string) Character {
	c.Slug = slug
	return c
}

// WithAliases sets the character's aliases.
func (c Character) WithAliases(aliases ...string) Character {
	c.Aliases = aliases
	return c
}

// WithDesc sets the character's description.
func (c Character) WithDesc(desc Text) Character {
	c.Desc = desc
	return c
}

// WithTraits sets the character's traits.
func (c Character) WithTraits(traits ...string) Character {
	c.Traits = traits
	return c
}
