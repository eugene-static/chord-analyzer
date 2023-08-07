package analyzer

import (
	"strconv"
	"strings"
)

type tabInfo struct {
	name    string
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

func newTabInfo(chordName, pattern string, fret int, capo bool) *tabInfo {
	return &tabInfo{
		name:    chordName,
		pattern: pattern,
		fret:    fret,
		capo:    capo,
	}
}

func (c *tabInfo) buildTab() string {
	chordTab := c.name + "\n"
	for _, fr := range c.pattern {
		switch fr {
		case 'X':
			chordTab += deadEnd + strings.Repeat(guitarString, 5) + "\n"
		case '0':
			chordTab += openString + strings.Repeat(guitarString, 5) + "\n"
		default:
			pos := int(fr - 48)
			chordTab += pushedString + strings.Repeat(guitarString, pos-1) +
				guitarString[:1] + finger + guitarString[2:] +
				strings.Repeat(guitarString, 5-pos) + "\n"
		}
	}
	if c.capo && c.fret != 0 {
		chordTab += capodastro
	} else {
		chordTab += space
	}
	for i := 1; i < 6; i++ {
		chordTab += doubleSpace + strconv.Itoa(c.fret+i) + space
	}
	return chordTab
}
