package analyzer

import (
	"strings"
)

const (
	ROOT           = iota
	MIN_SECOND     //m2
	MAJ_SECOND     //M2
	MIN_THIRD      //m3
	MAJ_THIRD      //M3
	PERFECT_FOURTH //P4
	FLAT_FIFTH     //d5
	PERFECT_FIFTH  //P5
	FLAT_SIXTH     //m6
	SIXTH          //M6
	MIN_SEVENTH    //m7
	MAJ_SEVENTH    //M7
)
const (
	SHARP_FIFTH     = FLAT_SIXTH
	DIM_SEVEN       = SIXTH
	FLAT_NINTH      = MIN_SECOND
	NINTH           = MAJ_SECOND
	SHARP_NINTH     = MIN_THIRD
	ELEVENTH        = PERFECT_FOURTH
	SHARP_ELEVENTH  = FLAT_FIFTH
	FLAT_THIRTEENTH = FLAT_SIXTH
	THIRTEENTH      = SIXTH
)
const (
	Q_SUS2 = iota
	Q_MIN
	Q_DUR
	Q_SUS4
	Q_DIM
	Q_AUG
	E_5
	E_B6
	E_B69
	E_B611
	E_B6911
	E_6
	E_69
	E_611
	E_6911
	E_7
	E_9
	E_11
	E_13
	E_MAJ7
	E_MAJ9
	E_MAJ11
	E_MAJ13
	A_B5
	A_S5
	A_B9
	A_S9
	A_S11
	A_B13
	A_ADDB9
	A_ADD9
	A_ADDS9
	A_ADD11
	A_ADDS11
	O_NO3
	O_NO5
	O_NO6
	O_NO7
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
	hooks := make([]bool, 38)
	for i := range hooks {
		hooks[i] = false
	}
	return &trapped{hooks: hooks,
		symbols: []string{"sus2", "m", "dur", "sus4", "dim", "aug",
			"5", "b6", "b6/9", "b6/11", "b6/9/11", "6", "6/9", "6/11", "6/9/11", "7",
			"9", "11", "13", "maj7", "maj9", "maj11", "maj13",
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
	trap.catch(O_NO3, O_NO5, O_NO6, O_NO7)
	if c[MAJ_SECOND] && !(c[MIN_THIRD] || c[MAJ_THIRD] || c[PERFECT_FOURTH]) {
		trap.cop(Q_SUS2)
		trap.rid(O_NO3)
	}
	if c[MIN_THIRD] && !c[MAJ_THIRD] {
		trap.cop(Q_MIN)
		trap.rid(O_NO3)
	}
	if c[MAJ_THIRD] {
		trap.cop(Q_DUR)
		trap.rid(O_NO3)
	}
	if c[PERFECT_FOURTH] && !(c[MIN_THIRD] || c[MAJ_THIRD]) {
		trap.cop(Q_SUS4)
		trap.rid(O_NO3)
	}
	if c[FLAT_FIFTH] && !c[PERFECT_FIFTH] {
		trap.cop(A_B5)
		trap.rid(O_NO5)
		if c[MIN_THIRD] && !(c[MIN_SEVENTH] || c[MAJ_SEVENTH]) {
			trap.cop(Q_DIM)
			trap.release(Q_MIN, A_B5)
		}
	}
	if c[PERFECT_FIFTH] {
		trap.cop(E_5)
		trap.rid(O_NO5)
	}
	if c[SHARP_FIFTH] && !(c[MAJ_SECOND] || c[MIN_THIRD] || /*c[PERFECT_FOURTH] ||*/ c[PERFECT_FIFTH]) {
		trap.cop(A_S5)
		trap.rid(O_NO5)
		if c[MAJ_THIRD] && !(c[MIN_SEVENTH] || c[MAJ_SEVENTH]) {
			trap.cop(Q_AUG)
			trap.release(Q_DUR, A_B5, A_S5)
		}
	}
	if c[FLAT_SIXTH] && !(c[MAJ_THIRD] || c[SIXTH] || c[MIN_SEVENTH] || c[MAJ_SEVENTH]) {
		if !trap.hooks[Q_AUG] {
			trap.cop(E_B6)
			trap.rid(O_NO6)
			//	trap.rid(E_B6)
			//	trap.cop(O_NO6)
		}
	}
	if c[SIXTH] && !(c[MIN_SEVENTH] || c[MAJ_SEVENTH]) {
		trap.cop(E_6)
		trap.rid(O_NO6)
		if trap.hooks[Q_DIM] {
			trap.release(E_6, O_NO7)
			trap.catch(E_7, O_NO6)
		}
	}
	if c[MIN_SEVENTH] {
		if !trap.hooks[Q_DIM] {
			trap.cop(E_7)
			trap.rid(O_NO7)
		}
	}
	if c[MAJ_SEVENTH] {
		if !trap.hooks[Q_DIM] {
			trap.cop(E_MAJ7)
			trap.rid(O_NO7)
		}
	}
	if c[FLAT_NINTH] {
		if trap.hooks[O_NO7] && trap.hooks[O_NO6] {
			trap.cop(A_ADDB9)
		} else {
			trap.cop(A_B9)
		}
	}
	if c[NINTH] {
		if !trap.hooks[Q_SUS2] {
			if trap.hooks[O_NO7] && trap.hooks[O_NO6] {
				trap.cop(A_ADD9)
			} else {
				if trap.hooks[O_NO7] {
					trap.cop(E_9)
					trap.rid(E_B6)
				} else {
					if trap.hooks[E_7] {
						trap.cop(E_9)
						trap.rid(E_7)
					}
					if trap.hooks[E_MAJ7] {
						trap.cop(E_MAJ9)
						trap.rid(E_MAJ7)
					}
				}
			}
		}

	}
	if c[SHARP_NINTH] {
		if !(trap.hooks[Q_MIN] || trap.hooks[Q_DIM]) {
			if trap.hooks[O_NO7] && trap.hooks[O_NO6] {
				trap.cop(A_ADDS9)
			} else {
				trap.cop(A_S9)
			}
		}
	}
	if c[ELEVENTH] {
		if !trap.hooks[Q_SUS4] {
			if trap.hooks[O_NO7] && trap.hooks[O_NO6] {
				trap.cop(A_ADD11)
			} else {
				if trap.hooks[O_NO7] {
					trap.cop(E_11)
				} else {
					if trap.hooks[E_7] || trap.hooks[E_9] {
						trap.cop(E_11)
						trap.release(E_7, E_9)
					}
					if trap.hooks[E_MAJ7] || trap.hooks[E_MAJ9] {
						trap.cop(E_MAJ11)
						trap.release(E_MAJ7, E_MAJ9)
					}
				}
			}
		}
	}
	if c[SHARP_ELEVENTH] {
		if !(trap.hooks[Q_MIN] || trap.hooks[Q_DIM] || trap.hooks[A_B5]) {
			if trap.hooks[O_NO6] && trap.hooks[O_NO7] {
				trap.cop(A_ADDS11)
			} else {
				trap.cop(A_S11)
			}
		}
	}
	if c[FLAT_THIRTEENTH] {
		//if !trap.hooks[A_S5] && !trap.hooks[Q_AUG] && (!trap.hooks[O_NO7] || !trap.hooks[O_NO6]) {
		//	trap.cop(A_B13)
		//	trap.rid(E_B6)
		//}
		if !trap.hooks[Q_AUG] && !trap.hooks[O_NO7] {
			trap.cop(A_B13)
			trap.release(E_B6, A_S5)
		}
	}
	if c[THIRTEENTH] {
		if !trap.hooks[Q_DIM] && trap.hooks[O_NO6] {
			if trap.hooks[E_7] || trap.hooks[E_9] || trap.hooks[E_11] {
				trap.cop(E_13)
				trap.release(E_7, E_9, E_11)
			}
			if trap.hooks[E_MAJ7] || trap.hooks[E_MAJ9] || trap.hooks[E_MAJ11] {
				trap.cop(E_MAJ13)
				trap.release(E_MAJ7, E_MAJ9, E_MAJ11)
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
	case trap.hooks[Q_DUR]:
		root = trap.majorMap[root]
	case trap.hooks[Q_DIM], trap.hooks[Q_AUG]:
		root = trap.sharpMap[root]
	default:
		root = trap.minorMap[root]
	}
	for i, r := range trap.hooks {
		if r {
			switch {
			case i < E_5:
				if i != Q_DUR {
					quality = trap.symbols[i]
				}
			case i < A_B5:
				if i == E_5 {
					if length == 2 {
						extended = trap.symbols[i]
						return
					}
				} else {
					ext = append(ext, trap.symbols[i])
				}
			case i < O_NO3:
				alt = append(alt, trap.symbols[i])
			case i == O_NO3:
				omitted = trap.symbols[i]
			}
		}
	}
	extended = strings.Join(ext, "/")
	altered = strings.Join(alt, ",")
	return
}
