package analyzer

import (
	"errors"
	"fmt"
	"unicode"
)

// ChordInfo stores request information
//
// Pattern must look like "01220X" from the highest string to the lowest, have length = 6
// and consist of 'X' for muted strings and digit from 0 to 5.
// If Fret is 0 the chord will be Am. If Fret == 2 and Capo == false, the chord will be A6/9sus4.
// If Fret == 2 and Capo == true, the chord will be Bm
//
// Capo influences on opened strings. If it is false, open strings (0 in pattern) will be calculated as open ¯\_(ツ)_/¯
// If true, open strings will be calculated as there is capo on Fret.
//
// If you want to use frets over 5th, just increase Fret value. It supports frets up to 18 (+5).
type ChordInfo struct {
	Pattern string
	Fret    int
	Capo    bool
}

const (
	patternLength = 6
	maxFretNumber = 18
)

// ChordNames stores information about all variations of chord, built on notes in pattern.
//
// Field Base is for chord with the lowest fingered string.
//
// Field Variations is for other chords, which can be constructed using same notes.
type ChordNames struct {
	Base       ChordName
	Variations []ChordName
}

// ChordName stores information about chord construction.
// Example:
//
// Cm7(addb9)
//
// Root: C;
//
// Quality: m (also can be sus2, sus4, dim or aug);
//
// Extended: 7 (has various values, like 13, maj11, b6/9, etc...);
//
// Altered: addb9 (if more, they will be separated with comma);
//
// Omitted: this field will not be empty when the chord will not have some kind of 'third',
// i.e. the chord will not be minor, major or suspended.
type ChordName struct {
	Root     string
	Quality  string
	Extended string
	Altered  string
	Omitted  string
}

// EmptyError can be used for preventing calculating if pattern has got no notes; ex: "XXXXXX"
var EmptyError = errors.New("invalid request: pattern must contain at list one digit")

// NewChordInfo returns new storage for request information
func NewChordInfo(pattern string, fret int, capo bool) *ChordInfo {
	return &ChordInfo{
		Pattern: pattern,
		Fret:    fret,
		Capo:    capo,
	}
}

// GetNames calculates intervals from guitar chord pattern and returns their symbolic values according to music theory.
// Guitar pattern contains information about fingering of the chord
// Ex: "00023X", fret = 2, capo = true --> is equal to Dmaj7 chord.
// Method 'calculateNotes' deletes repeating notes, so it doesn't give information about note octave.
// After calculation "00023X" ("F# D A F# D X") becomes "DF#A" and transforms into interval array that looks like:
//
// "[true false false false true false false true false false false]"
//
// Further analyzing is quite complex series of "if else" statements.
// At the end it returns struct with information about base chord, which note is on lowest string, and array of chords,
// which can be constructed from used notes.
func (c *ChordInfo) GetNames() (*ChordNames, error) {
	chordPattern := newNameInfo(c.Pattern, c.Fret, c.Capo)
	err := validate(c.Pattern, c.Fret)
	if err != nil {
		return nil, err
	}
	var baseChordName ChordName
	var variations []ChordName
	notes, baseRoot, length := chordPattern.calculateNotes()
	for bass, intervals := range notes {
		root, quality, extended, altered, omitted := getNames(bass, intervals, length)
		if bass == baseRoot {
			baseChordName = ChordName{
				Root:     root,
				Quality:  quality,
				Extended: extended,
				Altered:  altered,
				Omitted:  omitted,
			}
		} else {
			variations = append(variations, ChordName{
				Root:     root,
				Quality:  quality,
				Extended: extended,
				Altered:  altered,
				Omitted:  omitted,
			})
		}
	}
	return &ChordNames{
		Base:       baseChordName,
		Variations: variations,
	}, nil
}

// BuildName returns string constructed from ChordName fields
func (c *ChordName) BuildName() string {
	var name string
	name = c.Root
	if c.Quality == "sus2" || c.Quality == "sus4" {
		name += c.Extended + c.Quality
	} else {
		name += c.Quality + c.Extended
	}
	if c.Altered != "" {
		name += fmt.Sprintf("(%s)", c.Altered)
	}
	name += c.Omitted
	return name
}

// BuildTab returns string containing chord fingering tab
func (c *ChordInfo) BuildTab(name string) (string, error) {
	if len(name) == 0 {
		return "", errors.New("chord name can't be empty")
	}
	if len(name) > 20 {
		return "", errors.New("chord name is too long")
	}
	info := newTabInfo(c.Pattern, c.Fret, c.Capo)
	return info.buildTab(name), nil
}

func (c *ChordInfo) BuildPNG(name string) ([]byte, error) {
	info := newPNGInfo(name, c.Pattern, c.Fret, c.Capo)
	return info.buildPNG()
}

func validate(pattern string, fret int) error {
	if len(pattern) != patternLength {
		return lengthError
	}
	countX := 0
	for _, r := range pattern {
		if !unicode.IsDigit(r) {
			if r == x {
				countX++
			} else {
				return wrongSymbolsError
			}
		} else if r-48 > 5 {
			return fretPatternError
		}
	}
	if countX == patternLength {
		return EmptyError
	}
	if fret < 0 || fret > maxFretNumber {
		return fretNumberError
	}
	return nil
}
