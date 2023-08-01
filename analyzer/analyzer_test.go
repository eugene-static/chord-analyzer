package analyzer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetNames(t *testing.T) {
	testCase := []struct {
		pattern  string
		fret     int
		expected *ChordNames
		err      error
	}{
		// not predictable due to using maps. Slice order can differ
		//{
		//	pattern: "231231",
		//	fret:    1,
		//	expected: &ChordNames{
		//		Base: &chordName{
		//			Root:     "F",
		//			Quality:  "m",
		//			Extended: "maj13",
		//			Altered:  "b9",
		//			Omitted:  "",
		//		},
		//		Variations: []*chordName{
		//			{
		//				Root:     "C",
		//				Quality:  "aug",
		//				Extended: "",
		//				Altered:  "b5,add9,add11",
		//				Omitted:  "",
		//			},
		//			{
		//				Root:     "E",
		//				Quality:  "",
		//				Extended: "9",
		//				Altered:  "b9",
		//				Omitted:  "b13",
		//			},
		//			{
		//				Root:     "Ab",
		//				Quality:  "",
		//				Extended: "13",
		//				Altered:  "#5,#11",
		//				Omitted:  "",
		//			},
		//			{
		//				Root:     "D",
		//				Quality:  "",
		//				Extended: "9",
		//				Altered:  "#9",
		//				Omitted:  "#11",
		//			},
		//			{
		//				Root:     "F#",
		//				Quality:  "sus2",
		//				Extended: "7,maj7",
		//				Altered:  "#11,b13",
		//				Omitted:  "",
		//			},
		//		},
		//	},
		//	err: nil,
		//},
		{
			pattern: "XXX20X",
			fret:    4,
			expected: &ChordNames{
				Base: chordName{
					Root:     "C#",
					Quality:  "",
					Extended: "5",
					Altered:  "",
					Omitted:  "",
				},
				Variations: []chordName{
					{
						Root:     "G#",
						Quality:  "sus4",
						Extended: "",
						Altered:  "",
						Omitted:  "",
					},
				},
			},
			err: nil,
		},
		{
			pattern:  "X123452X",
			fret:     1,
			expected: nil,
			err:      lengthError,
		},
		{
			pattern:  "XXXXXX",
			fret:     1,
			expected: nil,
			err:      EmptyError,
		},
		{
			pattern:  "XY456X",
			fret:     1,
			expected: nil,
			err:      wrongSymbolsError,
		},
		{
			pattern:  "X15XXX",
			fret:     25,
			expected: nil,
			err:      fretNumberError,
		},
	}
	for _, r := range testCase {
		actual, err := GetNames(r.pattern, r.fret)
		assert.Equal(t, r.expected, actual)
		if r.err != nil {
			assert.EqualError(t, err, r.err.Error())
		}
	}
}
