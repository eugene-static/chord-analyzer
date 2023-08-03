package analyzer

import (
	"errors"
)

type nameInfo struct {
	pattern string
	fret    int
	capo    bool
}

func newNameInfo(pattern string, fret int, capo bool) *nameInfo {
	return &nameInfo{
		pattern: pattern,
		fret:    fret,
		capo:    capo,
	}
}

const (
	x = 'X'

	eString = 0
	bString = 7
	gString = 3
	dString = 10
	aString = 5
)

var (
	lengthError       = errors.New("invalid request: pattern must consist of six symbols")
	wrongSymbolsError = errors.New("invalid request: pattern must contain only digits and 'X'")
	fretNumberError   = errors.New("invalid request: offset fret number must be positive and less or equal '18'")
	fretPatternError  = errors.New("invalid request: fret number must be less or equal '5'")
)

func (c *nameInfo) calculateNotes() (map[int][]bool, int, int) {
	var root, length int
	var intervals []bool
	res := make(map[int][]bool)
	for i, n := range c.pattern {
		if n != x {
			note := findNote(i, int(n), c.fret, c.capo)
			intervals, length = c.getIntervals(note)
			if _, ok := res[note]; !ok {
				res[note] = intervals
			}
			root = note
		}
	}
	return res, root, length
}

func (c *nameInfo) getIntervals(noteIndex int) ([]bool, int) {
	iArr := make([]bool, 12)
	length := 0
	for i, n := range c.pattern {
		if n != x {
			note := findNote(i, int(n), c.fret, c.capo)
			ivl := (12 - (noteIndex - note)) % 12
			if !iArr[ivl] {
				length++
				iArr[ivl] = true
			}
		}
	}
	return iArr, length
}

func findNote(str, pos, fret int, capo bool) (res int) {
	if pos == 48 && fret != 0 && !capo {
		fret = 0
	}
	switch str {
	case 0, 5:
		res = (eString + fret + pos - 48) % 12
	case 1:
		res = (bString + fret + pos - 48) % 12
	case 2:
		res = (gString + fret + pos - 48) % 12
	case 3:
		res = (dString + fret + pos - 48) % 12
	case 4:
		res = (aString + fret + pos - 48) % 12
	}

	return
}
