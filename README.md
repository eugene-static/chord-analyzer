# Chord Analyzer
***
***beta-test***

Package provides method for analyzing chord based on guitar fingering pattern.

The guitar fingering pattern consist of the string pattern like "34353X"
from highest to lowest string and offset fret number.
The method returns interval information about chord: root, quality, extension, alteration
and omission. It also returns other chords, that can be constructed basing on used notes. 

### Example:
```
pattern := "34353X"
fret := 0
chords, err := analyzer.GetNames(pattern, fret)
if err != nil {
    fmt.Println(err)
    return
}

fmt.Printf("Base chord:\nRoot: %s, Quality: %s, Extension: %s, Alteration: %s, Omission: %s\n",
    chords.Base.Root, chords.Base.Quality, chords.Base.Extended, chords.Base.Altered, chords.Base.Omitted)
for _, chord := range chords.Variations {
    fmt.Printf("\tChord:\n\tRoot: %s, Quality: %s, Extension: %s, Alteration: %s, Omission: %s\n",
        chord.Root, chord.Quality, chord.Extended, chord.Altered, chord.Omitted)
}
```
**Result:**
```
Base chord:
Root: C, Quality: m, Extension: 7, Alteration: , Omission:
        Chord:
        Root: G, Quality: m, Extension: b6/11, Alteration: , Omission:
        Chord:
        Root: Eb, Quality: , Extension: 6, Alteration: , Omission:
        Chord:
        Root: Bb, Quality: sus4, Extension: 6/9, Alteration: , Omission:
```
