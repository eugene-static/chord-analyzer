package analyzer

import (
	"errors"
	"unicode"
)

type chordPattern struct {
	Pattern string
	Fret    int
}

func newChordPattern(pattern string, fret int) *chordPattern {
	return &chordPattern{
		Pattern: pattern,
		Fret:    fret,
	}
}

const (
	x             = 'X'
	patternLength = 6
	maxFretNumber = 23
)

var (
	lengthError       = errors.New("invalid request: pattern must consist of six symbols")
	wrongSymbolsError = errors.New("invalid request: pattern must contain only digits and 'X'")
	fretNumberError   = errors.New("invalid request: offset fret number must be positive and less than '23'")
)

func (c *chordPattern) validate() error {
	if len(c.Pattern) != patternLength {
		return lengthError
	}
	countX := 0
	for _, r := range c.Pattern {
		if !unicode.IsDigit(r) {
			if r == x {
				countX++
			} else {
				return wrongSymbolsError
			}
		}
	}
	if countX == patternLength {
		return EmptyError
	}
	if c.Fret < 0 || c.Fret > maxFretNumber {
		return fretNumberError
	}
	return nil
}
func (c *chordPattern) calculateNotes() (map[int][]bool, int, int) {
	var root, length int
	var intervals []bool
	res := make(map[int][]bool)
	for i, n := range c.Pattern {
		if n != x {
			note := findNote(i, int(n), c.Fret)
			intervals, length = c.getIntervals(note)
			if _, ok := res[note]; !ok {
				res[note] = intervals
				root = note
			}
		}
	}
	return res, root, length
}

func (c *chordPattern) getIntervals(noteIndex int) ([]bool, int) {
	iArr := make([]bool, 12)
	length := 0
	for i, n := range c.Pattern {
		if n != x {
			note := findNote(i, int(n), c.Fret)
			ivl := (12 - (noteIndex - note)) % 12
			if !iArr[ivl] {
				length++
				iArr[ivl] = true
			}
		}
	}
	return iArr, length
}

func findNote(str, pos, fret int) (res int) {
	switch str {
	case 0, 5:
		res = (fret + pos - 48) % 12
	case 1:
		res = (fret + pos - 41) % 12
	case 2:
		res = (fret + pos - 45) % 12
	case 3:
		res = (fret + pos - 38) % 12
	case 4:
		res = (fret + pos - 43) % 12
	}
	return
}
