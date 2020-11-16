package robots

// State shared between robots
type State struct {
	SearchTerm             string
	Prefix                 string
	SourceContentOriginal  string
	SourceContentSanitized string
	Sentences              []string
}

// Robot is the base for all robots
type Robot interface {
	Start(state *State) error
}
