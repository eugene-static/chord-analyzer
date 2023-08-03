package analyzer

import (
	"strings"
)

const (
	minSecond     = 1 + iota //m2
	majSecond                //M2
	minThird                 //m3
	majThird                 //M3
	perfectFourth            //P4
	flatFifth                //d5
	perfectFifth             //P5
	flatSixth                //m6
	sixth                    //M6
	minSeventh               //m7
	majSeventh               //M7
)
const (
	sharpFifth     = flatSixth
	flatNinth      = minSecond
	ninth          = majSecond
	sharpNinth     = minThird
	eleventh       = perfectFourth
	sharpEleventh  = flatFifth
	flatThirteenth = flatSixth
	thirteenth     = sixth
)
const (
	qsus2 = iota
	qmin
	qdur
	qsus4
	qdim
	qaug
	e5
	eb6
	e6
	e7
	e9
	e11
	e13
	emaj7
	emaj9
	emaj11
	emaj13
	ab5
	as5
	ab9
	as9
	as11
	ab13
	aaddb9
	aadd9
	aadds9
	aadd11
	aadds11
	ono3
	ono5
	ono6
	ono7
)

type trapped struct {
	hooks    []bool
	symbols  []string
	notes    []string
	majorMap map[string]string
	minorMap map[string]string
	sharpMap map[string]string
}

func prepare() *trapped {
	hooks := make([]bool, 32)
	for i := range hooks {
		hooks[i] = false
	}
	return &trapped{hooks: hooks,
		symbols: []string{"sus2", "m", "dur", "sus4", "dim", "aug",
			"5", "b6", "6", "7", "9", "11", "13", "maj7", "maj9", "maj11", "maj13",
			"b5", "#5", "b9", "#9", "#11", "b13", "addb9", "add9", "add#9", "add11", "add#11",
			"no3", "no5", "no6", "no7"},
		notes: []string{"E", "F", "fg", "G", "ga", "A", "ab", "B", "C", "cd", "D", "de"},
		majorMap: map[string]string{"E": "E", "F": "F", "fg": "Gb", "G": "G", "ga": "Ab", "A": "A", "ab": "Bb", "B": "B",
			"C": "C", "cd": "Db", "D": "D", "de": "Eb"},
		minorMap: map[string]string{"E": "E", "F": "F", "fg": "F#", "G": "G", "ga": "G#", "A": "A", "ab": "Bb", "B": "B",
			"C": "C", "cd": "C#", "D": "D", "de": "Eb"},
		sharpMap: map[string]string{"E": "E", "F": "F", "fg": "F#", "G": "G", "ga": "G#", "A": "A", "ab": "A#", "B": "B",
			"C": "C", "cd": "C#", "D": "D", "de": "D#"},
	}
}
func (trap *trapped) hide() {
	trap.hooks = nil
	trap.symbols = nil
	trap.notes = nil
	trap.majorMap = nil
	trap.minorMap = nil
	trap.sharpMap = nil
}

func (trap *trapped) release(i ...int) {
	for _, r := range i {
		trap.hooks[r] = false
	}
}
func (trap *trapped) catch(i ...int) {
	for _, r := range i {
		trap.hooks[r] = true
	}
}
func (trap *trapped) cop(i int) {
	trap.hooks[i] = true
}
func (trap *trapped) rid(i int) {
	trap.hooks[i] = false
}

func (trap *trapped) install(c []bool) {
	trap.catch(ono3, ono5, ono6, ono7)
	if c[majSecond] && !(c[minThird] || c[majThird] || c[perfectFourth]) {
		trap.cop(qsus2)
		trap.rid(ono3)
	}
	if c[minThird] && !c[majThird] {
		trap.cop(qmin)
		trap.rid(ono3)
	}
	if c[majThird] {
		trap.cop(qdur)
		trap.rid(ono3)
	}
	if c[perfectFourth] && !(c[minThird] || c[majThird]) {
		trap.cop(qsus4)
		trap.rid(ono3)
	}
	if c[flatFifth] && !c[perfectFifth] {
		trap.cop(ab5)
		trap.rid(ono5)
		if c[minThird] && !(c[minSeventh] || c[majSeventh]) {
			trap.cop(qdim)
			trap.release(qmin, ab5)
		}
	}
	if c[perfectFifth] {
		trap.cop(e5)
		trap.rid(ono5)
	}
	if c[sharpFifth] && !(c[majSecond] || c[minThird] || /*c[perfectFourth] ||*/ c[perfectFifth]) {
		trap.cop(as5)
		trap.rid(ono5)
		if c[majThird] && !(c[minSeventh] || c[majSeventh]) {
			trap.cop(qaug)
			trap.release(qdur, ab5, as5)
		}
	}
	if c[flatSixth] && !(c[majThird] || c[sixth] || c[minSeventh] || c[majSeventh]) {
		if !trap.hooks[qaug] {
			trap.cop(eb6)
			trap.rid(ono6)
			//	trap.rid(eb6)
			//	trap.cop(ono6)
		}
	}
	if c[sixth] && !(c[minSeventh] || c[majSeventh]) {
		trap.cop(e6)
		trap.rid(ono6)
		if trap.hooks[qdim] {
			trap.release(e6, ono7)
			trap.catch(e7, ono6)
		}
	}
	if c[minSeventh] {
		if !trap.hooks[qdim] {
			trap.cop(e7)
			trap.rid(ono7)
		}
	}
	if c[majSeventh] {
		if !trap.hooks[qdim] {
			trap.cop(emaj7)
			trap.rid(ono7)
		}
	}
	if c[flatNinth] {
		if trap.hooks[ono7] && trap.hooks[ono6] {
			trap.cop(aaddb9)
		} else {
			trap.cop(ab9)
		}
	}
	if c[ninth] {
		if !trap.hooks[qsus2] {
			if trap.hooks[ono7] && trap.hooks[ono6] {
				trap.cop(aadd9)
			} else {
				if trap.hooks[ono7] {
					trap.cop(e9)
					trap.rid(eb6)
				} else {
					if trap.hooks[e7] {
						trap.cop(e9)
						trap.rid(e7)
					}
					if trap.hooks[emaj7] {
						trap.cop(emaj9)
						trap.rid(emaj7)
					}
				}
			}
		}

	}
	if c[sharpNinth] {
		if !(trap.hooks[qmin] || trap.hooks[qdim]) {
			if trap.hooks[ono7] && trap.hooks[ono6] {
				trap.cop(aadds9)
			} else {
				trap.cop(as9)
			}
		}
	}
	if c[eleventh] {
		if !trap.hooks[qsus4] {
			if trap.hooks[ono7] && trap.hooks[ono6] {
				trap.cop(aadd11)
			} else {
				if trap.hooks[ono7] {
					trap.cop(e11)
				} else {
					if trap.hooks[e7] || trap.hooks[e9] {
						trap.cop(e11)
						trap.release(e7, e9)
					}
					if trap.hooks[emaj7] || trap.hooks[emaj9] {
						trap.cop(emaj11)
						trap.release(emaj7, emaj9)
					}
				}
			}
		}
	}
	if c[sharpEleventh] {
		if !(trap.hooks[qmin] || trap.hooks[qdim] || trap.hooks[ab5]) {
			if trap.hooks[ono6] && trap.hooks[ono7] {
				trap.cop(aadds11)
			} else {
				trap.cop(as11)
			}
		}
	}
	if c[flatThirteenth] {
		//if !trap.hooks[as5] && !trap.hooks[qaug] && (!trap.hooks[ono7] || !trap.hooks[ono6]) {
		//	trap.cop(ab13)
		//	trap.rid(eb6)
		//}
		if !trap.hooks[qaug] && !trap.hooks[ono7] {
			trap.cop(ab13)
			trap.release(eb6, as5)
		}
	}
	if c[thirteenth] {
		if !trap.hooks[qdim] && trap.hooks[ono6] {
			if trap.hooks[e7] || trap.hooks[e9] || trap.hooks[e11] {
				trap.cop(e13)
				trap.release(e7, e9, e11)
			}
			if trap.hooks[emaj7] || trap.hooks[emaj9] || trap.hooks[emaj11] {
				trap.cop(emaj13)
				trap.release(emaj7, emaj9, emaj11)
			}
		}
	}
}

// getIntervalNames uses 'install' method to get matches between requested intervals and their names in trap.symbols array
// according to music theory
func getIntervalsNames(rootIndex int, intervals []bool, length int) (root, quality, extended, altered, omitted string) {
	trap := prepare()
	defer trap.hide()
	root = trap.notes[rootIndex]
	if length == 1 {
		root = trap.minorMap[root]
		return
	}
	trap.install(intervals)
	var ext, alt []string
	switch {
	case trap.hooks[qdur]:
		root = trap.majorMap[root]
	case trap.hooks[qdim], trap.hooks[qaug]:
		root = trap.sharpMap[root]
	default:
		root = trap.minorMap[root]
	}
	for i, r := range trap.hooks {
		if r {
			switch {
			case i < e5:
				if i != qdur {
					quality = trap.symbols[i]
				}
			case i < ab5:
				if i == e5 {
					if length == 2 {
						extended = trap.symbols[i]
						return
					}
				} else {
					ext = append(ext, trap.symbols[i])
				}
			case i < ono3:
				alt = append(alt, trap.symbols[i])
			case i == ono3:
				omitted = trap.symbols[i]
			}
		}
	}
	extended = strings.Join(ext, "/")
	altered = strings.Join(alt, ",")
	return
}
