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
	chordTab := name + "\n"
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
	var sp string
	for i := 1; i < 6; i++ {
		if c.fret+i < 10 {
			sp = space
		} else {
			sp = ""
		}
		chordTab += doubleSpace + strconv.Itoa(c.fret+i) + sp
	}
	return chordTab
}
