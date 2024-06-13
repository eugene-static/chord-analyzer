package analyzer

import (
	"strconv"
	"strings"
)

type tabInfo struct {
	pattern string
	fret    int
	capo    bool
}

const (
	guitarString = "---|"
	deadEnd      = "X|"
	openString   = "0|"
	pushedString = "-|"
	space        = "\u00A0"
	doubleSpace  = space + space
	capodastro   = "c"
	finger       = "#"
)

func newTabInfo(pattern string, fret int, capo bool) *tabInfo {
	return &tabInfo{
		pattern: pattern,
		fret:    fret,
		capo:    capo,
	}
}

func (c *tabInfo) buildTab(name string) string {
	chordTab := strings.Builder{}
	chordTab.WriteString(name)
	chordTab.WriteRune('\n')
	for _, fr := range c.pattern {
		switch fr {
		case 'X':
			chordTab.WriteString(deadEnd)
			chordTab.WriteString(strings.Repeat(guitarString, 5))
		case '0':
			chordTab.WriteString(openString)
			chordTab.WriteString(strings.Repeat(guitarString, 5))
		default:
			pos := int(fr - 48)
			chordTab.WriteString(pushedString)
			chordTab.WriteString(strings.Repeat(guitarString, pos-1))
			chordTab.WriteString(guitarString[:1])
			chordTab.WriteString(finger)
			chordTab.WriteString(guitarString[2:])
			chordTab.WriteString(strings.Repeat(guitarString, 5-pos))
		}
		chordTab.WriteRune('\n')
	}
	if c.capo && c.fret != 0 {
		chordTab.WriteString(capodastro)
	} else {
		chordTab.WriteString(space)
	}
	var sp string
	for i := 1; i < 6; i++ {
		if c.fret+i < 10 {
			sp = space
		} else {
			sp = ""
		}
		chordTab.WriteString(doubleSpace)
		chordTab.WriteString(strconv.Itoa(c.fret + i))
		chordTab.WriteString(sp)
	}
	return chordTab.String()
}
