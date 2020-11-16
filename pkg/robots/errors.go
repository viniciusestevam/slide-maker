package robots

import "errors"

var (
	ErrReadline              = errors.New("✖ error reading user's input")
	ErrAlgorithmInstatiation = errors.New("✖ error trying to instatiate WikipediaParser algorithm")
	ErrAlgorithm             = errors.New("✖ something occurred when trying to execute WikipediaParser algorithm")
)
