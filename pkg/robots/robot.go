package robots

type ImageURL = string

type Sentence struct {
	Text     string
	Keywords []string
	Images   []ImageURL
}

// State shared between robots
type State struct {
	SearchTerm             string
	Prefix                 string
	SourceContentOriginal  string
	SourceContentSanitized string
	Sentences              []*Sentence
}

// Robot is the base for all robots
type Robot interface {
	Start(state *State)
}
