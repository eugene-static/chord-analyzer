package analyzer

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"

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
	nameFontsize = 32
	fontDPI      = 72
	zero         = 0
	cellWidth    = 100
	cellHeight   = 60
	fretMax      = 18
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
	canvas := image.NewRGBA(image.Rect(0, 0, cellWidth*6+cellWidth/2, cellHeight*7+cellHeight/2))
	draw.Draw(canvas, canvas.Bounds(), fretboard, image.Pt(info.fret*cellWidth, zero), draw.Src)
	if info.fret != 0 {
		draw.Draw(canvas, image.Rect(0, 0, cellWidth, canvas.Bounds().Max.Y),
			fretboard, image.Pt(zero, zero), draw.Src)
	}
	if info.fret != fretMax {
		draw.Draw(canvas, image.Rect(cellWidth*6+2, 0, cellWidth*6+cellWidth/2, canvas.Bounds().Max.Y),
			fretboard, image.Pt(zero, zero), draw.Src)
	}
	cell := image.Rect(zero, zero, cellWidth, cellHeight)
	for i, str := range tab {
		height := i*cellHeight + cellHeight
		switch str {
		case -1:
			move(&cell, zero, height)
			draw.Draw(canvas, cell, symbols, mutedZP, draw.Over)
		case 0:
			move(&cell, zero, height)
			draw.Draw(canvas, cell, symbols, openZP, draw.Over)
		default:
			move(&cell, str*cellWidth, height)
			draw.Draw(canvas, cell, symbols, fingerZP, draw.Over)
		}
	}
	if info.capo && info.fret != 0 {
		move(&cell, zero, cellHeight*6+cellHeight/2)
		draw.Draw(canvas, cell, symbols, capoZP, draw.Over)
	}
	err = drawText(canvas, info.name, info.fret)
	if err != nil {
		return nil, err
	}
	return writeToBytes(canvas)
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
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	fontFile, err := os.ReadFile(filepath.Join(dir, verdanaPath))
	if err != nil {
		return err
	}
	fontFace, err := freetype.ParseFont(fontFile)
	faceOptions := &truetype.Options{
		Size:    nameFontsize,
		DPI:     fontDPI,
		Hinting: font.HintingNone,
	}
	fontDrawer := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color.White),
		Face: truetype.NewFace(fontFace, faceOptions),
	}
	fontDrawer.Dot = fixed.Point26_6{
		X: fixed.I(cellWidth) + (fixed.I(img.Bounds().Max.X-cellWidth-cellWidth/2)-fontDrawer.MeasureString(name))/2,
		Y: fixed.I(cellHeight+nameFontsize) / 2,
	}
	fontDrawer.DrawString(name)
	return nil
}
func readPNG(name string) (*image.RGBA, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	file, err := os.Open(filepath.Join(dir, name))
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
