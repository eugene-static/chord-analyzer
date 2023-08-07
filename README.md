# Chord Analyzer
***beta-test***

Package provides method for analyzing chord based on guitar fingering pattern.
It also includes method for building chord tab.

The guitar fingering pattern consist of the string pattern like "34353X"
from highest to lowest string and offset fret number.
The method returns interval information about chord: root, quality, extension, alteration
and omission. It also returns other chords, that can be constructed based on used notes. 

### Example:
```
pattern := "00023X"
fret := 2
capo := true
chord := analyzer.NewChordInfo(pattern, fret, capo)
names, err := chord.GetNames()
if err != nil {...}
	
fmt.Printf("Base chord:\nRoot: %s, Quality: %s, Extension: %s, Alteration: %s, Omission: %s\n",
    names.Base.Root, names.Base.Quality, names.Base.Extended, names.Base.Altered, names.Base.Omitted)
for _, name := range names.Variations {
    fmt.Printf("\tChord:\n\tRoot: %s, Quality: %s, Extension: %s, Alteration: %s, Omission: %s\n",
	name.Root, name.Quality, name.Extended, name.Altered, name.Omitted)
}
```
**Result:**
```
Base chord:
Root: D, Quality: , Extension: maj7, Alteration: , Omission:
        Chord:
        Root: F#, Quality: m, Extension: b6, Alteration: , Omission:
        Chord:
        Root: C#, Quality: sus4, Extension: b6, Alteration: #5,b9, Omission:
        Chord:
        Root: A, Quality: , Extension: 6/11, Alteration: , Omission:

```

To get name from analyzed chord, you can use 'BuildName' method.

```
name := names.Base.BuildName() // Dmaj7
```

Use 'BuildTab' method to get chord tab.
You can use your own name, if you are not agree with analyzed name

```
chord.Name = name
tab, err := chord.BuildTab()
fmt.Println(tab)
```
**Result:**
```
Dmaj7
0|---|---|---|---|---|
0|---|---|---|---|---|
0|---|---|---|---|---|
-|---|-#-|---|---|---|
-|---|---|-#-|---|---|
X|---|---|---|---|---|
c  3   4   5   6   7

// "c" is for capo
```
