package model

// Scene represents a screenplay scene.
type Scene struct {
	ID           string    `json:"id"`
	Authors      []string  `json:"authors"` // UUIDs
	Contributors []string  `json:"contributors,omitempty"`
	Heading      *Slugline `json:"heading"`
	Body         []Element `json:"body"`
	Cast         []string  `json:"cast,omitempty"` // Character UUIDs
	Animals      []string  `json:"animals,omitempty"`
	Extra        []string  `json:"extra,omitempty"`
	Locations    []string  `json:"locations,omitempty"`
	Moods        []string  `json:"moods,omitempty"`
	Props        []string  `json:"props,omitempty"`
	SFX          []string  `json:"sfx,omitempty"`
	Sounds       []string  `json:"sounds,omitempty"`
	Tags         []string  `json:"tags,omitempty"`
	VFX          []string  `json:"vfx,omitempty"`
	Wardrobe     []string  `json:"wardrobe,omitempty"`
	Meta         Meta      `json:"meta,omitempty"`
}

// Slugline represents a structured scene heading.
type Slugline struct {
	No      int      `json:"no,omitempty"`
	Context string   `json:"context"` // INT, EXT, INT/EXT, etc.
	Setting string   `json:"setting"`
	Time    string   `json:"time"` // DAY, NIGHT, etc.
	Mods    []string `json:"mods,omitempty"` // FLASHBACK, INTERCUT, etc.
	Desc    Text     `json:"desc,omitempty"`
	Meta    Meta     `json:"meta,omitempty"`
}

// ValidContexts are the allowed values for Slugline.Context.
var ValidContexts = []string{
	"I/E", "INT/EXT", "EXT/INT", "INT", "EXT", "POV",
}

// ValidTimes are common time-of-day values (custom values also allowed).
var ValidTimes = []string{
	"DAY", "NIGHT", "DAWN", "DUSK", "LATER", "MOMENTS LATER",
	"CONTINUOUS", "MORNING", "AFTERNOON", "EVENING", "THE NEXT DAY",
}
