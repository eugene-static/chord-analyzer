package analyzer

import "strings"

type filters struct {
	q []bool
	e []bool
	a []bool
}
type symbols struct {
	q     []string
	e     []string
	a     []string
	notes []string
	major map[string]string
	minor map[string]string
	sharp map[string]string
}

func install() *filters {
	return &filters{
		q: []bool{false, false, false, false, false, false, false},
		e: []bool{false, false, false, false, false, false, false, false, false, false, false, false, false},
		a: []bool{false, false, false, false, false, false, false, false, false, false, false, false},
	}
}
func initSymbols() *symbols {
	return &symbols{
		q:     []string{"sus2", "sus4", "m", "dim", "", "aug", "no3"},
		e:     []string{"5", "b6", "b13", "6", "7", "9", "11", "13", "maj7", "maj9", "maj11", "maj13", "onoE"},
		a:     []string{"b5", "#5", "b6", "b9", "#9", "#11", "b13", "addb9", "add9", "add#9", "add11", "add#11"},
		notes: []string{"E", "F", "fg", "G", "ga", "A", "ab", "B", "C", "cd", "D", "de"},
		major: map[string]string{"E": "E", "F": "F", "fg": "Gb", "G": "G", "ga": "Ab", "A": "A", "ab": "Bb", "B": "B",
			"C": "C", "cd": "Db", "D": "D", "de": "Eb"},
		minor: map[string]string{"E": "E", "F": "F", "fg": "F#", "G": "G", "ga": "G#", "A": "A", "ab": "Bb", "B": "B",
			"C": "C", "cd": "C#", "D": "D", "de": "Eb"},
		sharp: map[string]string{"E": "E", "F": "F", "fg": "F#", "G": "G", "ga": "G#", "A": "A", "ab": "A#", "B": "B",
			"C": "C", "cd": "C#", "D": "D", "de": "D#"},
	}
}

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
	qsus4
	qmin
	qdim
	qdur
	qaug
	ono3
)
const (
	e5 = iota
	eb6
	eb13
	e6
	e7
	e9
	e11
	e13
	emaj7
	emaj9
	emaj11
	emaj13
	onoe
)
const (
	ab5 = iota
	as5
	ab6
	ab9
	as9
	as11
	ab13
	addb9
	add9
	adds9
	add11
	adds11
)
const (
	comma = ","
	slash = "/"
)

func (f *filters) do(c []bool) {
	f.qualityFilter(c)
	f.extensionFilter(c)
	f.alterFilter(c)
}

func getNames(rootIndex int, intervals []bool, length int) (root, quality, extended, altered, omitted string) {
	filter := install()
	sym := initSymbols()
	defer func() {
		filter = nil
		sym = nil
	}()
	root = sym.notes[rootIndex]
	if length == 1 {
		root = sym.minor[root]
		return
	}
	switch {
	case filter.q[qdur]:
		root = sym.major[root]
	case filter.q[qdim], filter.q[qaug]:
		root = sym.sharp[root]
	default:
		root = sym.minor[root]
	}
	if filter.powerFilter(intervals, length) {
		extended = sym.e[e5]
		return
	}
	filter.do(intervals)
	var ext, alt []string
	if filter.q[ono3] {
		omitted = sym.q[ono3]
	} else {
		for i := 0; i < ono3; i++ {
			if filter.q[i] {
				quality = sym.q[i]
				break
			}
		}
	}
	if !filter.e[onoe] {
		for i := 1; i < onoe; i++ {
			if filter.e[i] {
				ext = append(ext, sym.e[i])
			}
		}
	}
	for i, a := range filter.a {
		if a {
			alt = append(alt, sym.a[i])
		}
	}
	extended = strings.Join(ext, filter.separator())
	if strings.ContainsRune(root, 'b') && quality == "" && filter.e[eb6] {
		altered = extended
		extended = ""
	}
	if altered != "" && alt != nil {
		altered += comma
	}
	altered += strings.Join(alt, comma)
	return
}

func (f *filters) powerFilter(c []bool, length int) bool {
	f.e[e5] = c[perfectFifth] && length == 2
	return f.e[e5]
}

func (f *filters) qualityFilter(c []bool) {
	var q int
	switch {
	case c[majThird]:
		q = qdur
	case c[minThird]:
		q = qmin
	case c[perfectFourth]:
		q = qsus4
	case c[majSecond]:
		q = qsus2
	default:
		q = ono3
	}
	if !c[perfectFifth] && !c[minSeventh] && !c[majSeventh] {
		if q == qdur && c[sharpFifth] && !c[flatFifth] { //!c[flatFifth]
			q = qaug
		}
		if q == qmin && c[flatFifth] && !c[sharpFifth] { //!c[sharpFifth]
			q = qdim
		}
	}
	f.qSelect(q)
}

func (f *filters) extensionFilter(c []bool) {
	switch {
	case f.q[qdim]:
		if c[sixth] {
			switch {
			case c[flatThirteenth]:
				f.eSelect(eb13)
			case c[eleventh]:
				f.eSelect(e11)
			case c[ninth]:
				f.eSelect(e9)
			default:
				f.eSelect(e7)
			}
		} else if c[flatSixth] {
			f.eSelect(eb6)
			f.e[e9] = c[ninth]
			f.e[e11] = c[eleventh]
		} else {
			f.eSelect(onoe)
		}
	case f.q[qaug]:
		if c[sixth] {
			f.eSelect(e6)
			f.e[e9] = c[ninth]
			f.e[e11] = c[eleventh]
		} else {
			f.eSelect(onoe)
		}

	case f.q[qsus2], f.q[qsus4], f.q[qmin], f.q[qdur], f.q[ono3]:
		if c[minSeventh] || c[majSeventh] {
			rise := 0
			switch {
			case c[thirteenth]:
				rise = e13
			case c[eleventh] && !f.q[qsus4]:
				rise = e11
			case c[ninth] && !f.q[qsus2]:
				rise = e9
			default:
				rise = e7
			}
			if c[minSeventh] {
				f.eSelect(rise)
			}
			if c[majSeventh] {
				f.eSelect(rise + 4)
			}
		} else if c[flatSixth] || c[sixth] {
			if c[sixth] {
				f.eSelect(e6)
			} else {
				f.eSelect(eb6)
			}
			f.e[e9] = c[ninth] && !f.q[qsus2]
			f.e[e11] = c[eleventh] && !f.q[qsus4]
		} else {
			f.eSelect(onoe)
		}
	}
}
func (f *filters) alterFilter(c []bool) {
	switch {
	case f.q[qdur], f.q[qsus4], f.q[qsus2]:
		if f.e[onoe] {
			f.a[addb9] = c[flatNinth]
			f.a[add9] = c[ninth] && !f.q[qsus2]
			f.a[adds9] = c[sharpNinth] && f.q[qdur]
			f.a[add11] = c[eleventh] && !f.q[qsus4]
			if c[perfectFifth] {
				f.a[adds11] = c[sharpEleventh]
				f.a[ab6] = c[flatSixth]
			} else {
				f.a[ab5] = c[flatFifth]
				f.a[as5] = c[sharpFifth]
			}
		} else {
			f.a[ab9] = c[flatNinth]
			f.a[as9] = c[sharpNinth]
			if c[perfectFifth] {
				f.a[as11] = c[sharpEleventh]
				if c[flatSixth] && !f.e[eb6] {
					if (f.e[e13] || f.e[emaj13]) && !f.e[e6] {
						f.aSelect(ab6)
					} else {
						f.aSelect(ab13)
					}
				}
			} else {
				f.a[ab5] = c[flatFifth]
				f.a[as5] = c[sharpFifth] && !f.e[eb6]
			}
		}
	case f.q[qmin]:
		if f.e[onoe] {
			f.a[addb9] = c[flatNinth]
			f.a[add9] = c[ninth]
			f.a[add11] = c[eleventh]
			if c[perfectFifth] {
				f.a[adds11] = c[sharpEleventh]
				f.a[ab6] = c[flatSixth]
			} else {
				f.a[ab5] = c[flatFifth]
				f.a[as5] = c[sharpFifth]
			}
		} else {
			f.a[ab9] = c[flatNinth]
			if c[flatFifth] {
				if c[perfectFifth] {
					f.aSelect(as11)
				} else {
					f.aSelect(ab5)
					f.a[as5] = c[sharpFifth] //&& !f.e[EB6]
				}
			} else if c[flatSixth] && !f.e[eb6] {
				if (f.e[e13] || f.e[emaj13]) && !f.e[e6] {
					f.aSelect(ab6)
				} else {
					f.aSelect(ab13)
				}
			}
		}
	case f.q[qdim]:
		if f.e[onoe] {
			f.a[addb9] = c[flatNinth]
			f.a[add9] = c[ninth]
			f.a[add11] = c[eleventh]
		} else {
			f.a[ab9] = c[flatNinth]
		}
	case f.q[qaug]:
		if f.e[onoe] {
			f.a[addb9] = c[flatNinth]
			f.a[add9] = c[ninth]
			f.a[adds9] = c[sharpNinth]
			f.a[add11] = c[eleventh]
			f.a[adds11] = c[sharpEleventh]
		} else {
			f.a[ab9] = c[flatNinth]
			f.a[as9] = c[sharpNinth]
		}
	case f.q[ono3]:
		if f.e[onoe] {
			f.a[addb9] = c[flatNinth]
			if c[perfectFifth] {
				f.a[adds11] = c[sharpEleventh]
				f.a[ab6] = c[flatSixth]
			} else {
				f.a[ab5] = c[flatFifth]
				f.a[as5] = c[sharpFifth]
			}
		} else {
			f.a[ab9] = c[flatNinth]
			f.a[as11] = c[sharpEleventh]
			if c[flatSixth] && !f.e[eb6] {
				if (f.e[e13] || f.e[emaj13]) && !f.e[e6] {
					f.aSelect(ab6)
				} else {
					f.aSelect(ab13)
				}
			}
		}
	}
}

func (f *filters) qSelect(i int) {
	f.q[i] = true
}
func (f *filters) eSelect(i int) {
	f.e[i] = true
}
func (f *filters) aSelect(i int) {
	f.a[i] = true
}
func (f *filters) separator() string {
	if f.e[e6] || f.e[eb6] {
		return slash
	}
	return comma
}
