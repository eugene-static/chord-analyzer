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

const (
	nameFontsize = 32
	fontDPI      = 72
	zero         = 0
	cellWidth    = 100
	cellHeight   = 60
	fretMax      = 18
)
const (
	fretBoardPath = "analyzer/assets/fretboard.png"
	symbolsPath   = "analyzer/assets/symbols.png"
	verdanaPath   = "analyzer/assets/verdana.ttf"
)

type pngInfo struct {
	Name    string
	Pattern string
	Fret    int
	Capo    bool
}

func newPNGInfo(name, pattern string, fret int, capo bool) *pngInfo {
	return &pngInfo{
		Name:    name,
		Pattern: pattern,
		Fret:    fret,
		Capo:    capo,
	}
}

func (info *pngInfo) buildPNG() ([]byte, error) {
	path, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	tab, err := info.toArray()
	if err != nil {
		return nil, err
	}
	fretboard, err := readPNG(filepath.Join(path, fretBoardPath))
	if err != nil {
		return nil, err
	}
	sym, err := readPNG(filepath.Join(path, symbolsPath))
	if err != nil {
		return nil, err
	}
	fingerZP := image.Pt(zero, zero)
	openZP := image.Pt(cellWidth, zero)
	mutedZP := image.Pt(cellWidth*2, zero)
	capoZP := image.Pt(cellWidth*3, zero)
	canvas := image.NewRGBA(image.Rect(0, 0, cellWidth*6+cellWidth/2, cellHeight*7+cellHeight/2))
	draw.Draw(canvas, canvas.Bounds(), fretboard, image.Pt(info.Fret*cellWidth, zero), draw.Src)
	if info.Fret != 0 {
		draw.Draw(canvas, image.Rect(0, 0, cellWidth, canvas.Bounds().Max.Y),
			fretboard, image.Pt(zero, zero), draw.Src)
	}
	if info.Fret != fretMax {
		draw.Draw(canvas, image.Rect(cellWidth*6+2, 0, cellWidth*6+cellWidth/2, canvas.Bounds().Max.Y),
			fretboard, image.Pt(zero, zero), draw.Src)
	}
	cell := image.Rect(zero, zero, cellWidth, cellHeight)
	for i, str := range tab {
		height := i*cellHeight + cellHeight
		switch str {
		case -1:
			move(&cell, zero, height)
			draw.Draw(canvas, cell, sym, mutedZP, draw.Over)
		case 0:
			move(&cell, zero, height)
			draw.Draw(canvas, cell, sym, openZP, draw.Over)
		default:
			move(&cell, str*cellWidth, height)
			draw.Draw(canvas, cell, sym, fingerZP, draw.Over)
		}
	}
	if info.Capo && info.Fret != 0 {
		move(&cell, zero, cellHeight*6+cellHeight/2)
		draw.Draw(canvas, cell, sym, capoZP, draw.Over)
	}
	err = info.drawText(canvas, path)
	if err != nil {
		return nil, err
	}
	dir, err := info.write(canvas)
	return dir, err
}
func (info *pngInfo) write(img *image.RGBA) ([]byte, error) {
	var b bytes.Buffer
	err := png.Encode(&b, img)
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func (info *pngInfo) drawText(img *image.RGBA, path string) error {
	fontData, err := os.ReadFile(filepath.Join(path, verdanaPath))
	if err != nil {
		return err
	}
	fontFace, err := freetype.ParseFont(fontData)
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
		X: fixed.I(cellWidth) + (fixed.I(img.Bounds().Max.X-cellWidth-cellWidth/2)-fontDrawer.MeasureString(info.Name))/2,
		Y: fixed.I(cellHeight+nameFontsize) / 2,
	}
	fontDrawer.DrawString(info.Name)
	return nil
}

func (info *pngInfo) toArray() (result []int, err error) {
	for _, r := range info.Pattern {
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

func readPNG(path string) (*image.RGBA, error) {
	fileData, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	grid, err := png.Decode(fileData)
	if err != nil {
		return nil, err
	}
	rect := image.Rect(grid.Bounds().Min.X, grid.Bounds().Min.Y, grid.Bounds().Max.X, grid.Bounds().Max.Y)
	img := image.NewRGBA(rect)
	draw.Draw(img, img.Bounds(), grid, image.Pt(0, 0), draw.Src)
	return img, nil
}

func move(cell *image.Rectangle, x, y int) {
	cell.Min.X = x
	cell.Min.Y = y
	cell.Max.X = x + cellWidth
	cell.Max.Y = y + cellHeight
}
