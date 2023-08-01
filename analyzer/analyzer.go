package analyzer

import (
	"errors"
)

// ChordNames stores information about chord construction.
// It represents Root, Quality, Extension, Alterations and Omitted fields
//
// Ex: Cm7(addb9)
//
// Root: C;
//
// Quality: m (also can be sus2, sus4, dim or aug);
//
// Extension: 7 (has various values, like 13, maj11, b6/9, etc...);
//
// Alterations: addb9 (if more, they will be separated with comma);
//
// Omitted: this field will not be empty when the chord will not have some kind of 'third',
// i.e. the chord will not be minor, major or suspended.
//
// Field Base is for chord with the lowest fingered string.
//
// Field Variations is for other chords, which can be constructed using same notes.
type ChordNames struct {
	Base       chordName
	Variations []chordName
}

type chordName struct {
	Root     string
	Quality  string
	Extended string
	Altered  string
	Omitted  string
}

// EmptyError can be used for preventing calculating if pattern has got no notes; ex: "XXXXXX"
var EmptyError = errors.New("invalid request: pattern must contain at list one digit")

// GetNames calculates intervals from guitar chord pattern and returns their symbolic values according to music theory.
// Guitar pattern contains information about fingering of the chord
// Ex: "01023X", Fret = 0 --> is equal to C major chord.
// Method 'calculateNotes' deletes repeating notes, so it doesn't give information about note octave.
// After calculation "01023X" ("ECGECX") becomes "CEG" and transforms in some interval array that looks like:
//
//	"[true false false false true false false true false false false]"
//
// Further analyzing is quite complex series of "if else" statements.
// At the end it returns struct with information about base chord, which note is on lowest string, and array of chords,
// which can be constructed from used notes.
// Pattern must have length = 6
// If you want to use frets over 9th, just increase "fret" value. It supports frets below 23.
func GetNames(pattern string, fret int) (*ChordNames, error) {
	chordPattern := newChordPattern(pattern, fret)
	err := chordPattern.validate()
	if err != nil {
		return nil, err
	}
	var baseChordName chordName
	var variations []chordName
	notes, baseRoot, length := chordPattern.calculateNotes()
	for bass, intervals := range notes {
		root, quality, extended, altered, omitted := getIntervalsNames(bass, intervals, length)
		if bass == baseRoot {
			baseChordName = chordName{
				Root:     root,
				Quality:  quality,
				Extended: extended,
				Altered:  altered,
				Omitted:  omitted,
			}
		} else {
			variations = append(variations, chordName{
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
