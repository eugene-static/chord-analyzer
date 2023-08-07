package analyzer

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"

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
	fretBoardPath = "fretboard.png"
	symbolsPath   = "symbols.png"
	verdanaPath   = "verdana.ttf"
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
	storage, err := NewStorage("./data")
	if err != nil {
		return nil, err
	}
	fretboard, err := readPNG(storage, fretBoardPath)
	if err != nil {
		return nil, err
	}
	symbols, err := readPNG(storage, symbolsPath)
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
	err = drawText(storage, canvas, info.name)
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

func drawText(storage *Storage, img *image.RGBA, name string) error {
	fontData, err := storage.Get(verdanaPath)
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
		X: fixed.I(cellWidth) + (fixed.I(img.Bounds().Max.X-cellWidth-cellWidth/2)-fontDrawer.MeasureString(name))/2,
		Y: fixed.I(cellHeight+nameFontsize) / 2,
	}
	fontDrawer.DrawString(name)
	return nil
}
func readPNG(storage *Storage, name string) (*image.RGBA, error) {
	fileData, err := storage.Get(name)
	if err != nil {
		return nil, err
	}
	grid, err := png.Decode(bytes.NewReader(fileData))
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

//func load(storage *Storage, names ...string) error {
//	for _, name := range names {
//		data, err := os.ReadFile(name)
//		if err != nil {
//			return err
//		}
//		err = storage.Save(filepath.Base(name), data)
//	}
//	return nil
//}
