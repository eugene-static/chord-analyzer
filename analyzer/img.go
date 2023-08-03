package analyzer

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"strconv"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type drawInfo struct {
	name    string
	pattern string
	fret    int
	capo    bool
}

const (
	nameFontsize  = 24
	fretsFontSize = 18
	fontDPI       = 72
	zero          = 0
	cellWidth     = 75
	cellHeight    = 50
	grayColor     = 25
)
const (
	fretBoardPath = "analyzer/src/fretboard.png"
	symbolsPath   = "analyzer/src/symbols.png"
	verdanaPath   = "analyzer/src/verdana.ttf"
)

func newDrawInfo(name, pattern string, fret int, capo bool) *drawInfo {
	return &drawInfo{
		name:    name,
		pattern: pattern,
		fret:    fret,
		capo:    capo,
	}
}

func (info *drawInfo) buildPNG() ([]byte, error) {
	tab, err := toArray(info.pattern)
	if err != nil {
		return nil, err
	}
	fretboard, err := readPNG(fretBoardPath)
	if err != nil {
		return nil, err
	}
	symbols, err := readPNG(symbolsPath)
	if err != nil {
		return nil, err
	}
	fingerZP := image.Pt(zero, zero)
	openZP := image.Pt(cellWidth, zero)
	mutedZP := image.Pt(cellWidth*2, zero)
	capoZP := image.Pt(cellWidth*3, zero)
	cell := image.Rect(zero, zero, cellWidth, cellHeight)
	for i, str := range tab {
		height := i*cellHeight + cellHeight
		switch str {
		case -1:
			move(&cell, zero, height)
			draw.Draw(fretboard, cell, symbols, mutedZP, draw.Over)
		case 0:
			move(&cell, zero, height)
			draw.Draw(fretboard, cell, symbols, openZP, draw.Over)
		default:
			move(&cell, str*cellWidth, height)
			draw.Draw(fretboard, cell, symbols, fingerZP, draw.Over)
		}
	}
	if info.capo && info.fret != 0 {
		move(&cell, zero, cellHeight*7)
		draw.Draw(fretboard, cell, symbols, capoZP, draw.Over)
	}
	err = drawText(fretboard, info.name, info.fret)
	if err != nil {
		return nil, err
	}
	return writeToBytes(fretboard)
}
func writeToBytes(img *image.RGBA) ([]byte, error) {
	var b []byte
	buffer := bytes.NewBuffer(b)
	err := png.Encode(buffer, img)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func drawText(img *image.RGBA, name string, fret int) error {
	fontFile, err := os.ReadFile(verdanaPath)
	if err != nil {
		return err
	}
	bg := color.Gray{Y: grayColor}
	cell := image.Rect(zero, zero, cellWidth*5, cellHeight)
	canvas := image.NewRGBA(cell)
	fontFace, err := freetype.ParseFont(fontFile)
	draw.Draw(canvas, cell, &image.Uniform{C: bg}, image.Pt(zero, zero), draw.Src)
	faceOptions := &truetype.Options{
		Size:    fretsFontSize,
		DPI:     fontDPI,
		Hinting: font.HintingNone,
	}
	fontDrawer := &font.Drawer{
		Dst:  canvas,
		Src:  image.NewUniform(color.White),
		Face: truetype.NewFace(fontFace, faceOptions),
	}
	fontDrawer.Dot = fixed.Point26_6{
		Y: fixed.I(cellHeight / 2),
	}

	for i := 0; i < 5; i++ {
		txt := strconv.Itoa(fret + i + 1)
		fontDrawer.Dot.X = (fixed.I(cellWidth)-fontDrawer.MeasureString(txt))/2 + fixed.I(cellWidth*i)
		fontDrawer.DrawString(txt)
	}
	cell = image.Rect(cellWidth, cellHeight*7, cellWidth*6, cellHeight*8)
	draw.Draw(img, cell, canvas, image.Pt(zero, zero), draw.Over)
	faceOptions.Size = nameFontsize
	fontDrawer.Dst = img
	fontDrawer.Face = truetype.NewFace(fontFace, faceOptions)
	fontDrawer.Dot = fixed.Point26_6{
		X: fixed.I(cellWidth) + (fixed.I(img.Bounds().Max.X-cellWidth-cellWidth/3)-fontDrawer.MeasureString(name))/2,
		Y: fixed.I(cellHeight+nameFontsize) / 2,
	}
	fontDrawer.DrawString(name)
	return nil
}
func readPNG(name string) (*image.RGBA, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	grid, err := png.Decode(file)
	if err != nil {
		return nil, err
	}
	rect := image.Rect(grid.Bounds().Min.X, grid.Bounds().Min.Y, grid.Bounds().Max.X, grid.Bounds().Max.Y)
	img := image.NewRGBA(rect)
	draw.Draw(img, img.Bounds(), grid, image.Pt(0, 0), draw.Src)
	return img, nil
}

func toArray(pattern string) (result []int, err error) {
	for _, r := range pattern {
		if r == 'X' {
			result = append(result, -1)
		} else {
			trueFret := int(r - 48)
			if trueFret > 5 {
				err = fmt.Errorf("fret number '%d' must be < 6", trueFret)
				return nil, err
			}
			result = append(result, int(r-48))
		}
	}
	return
}

func move(cell *image.Rectangle, x, y int) {
	cell.Min.X = x
	cell.Min.Y = y
	cell.Max.X = x + cellWidth
	cell.Max.Y = y + cellHeight
}
