package test

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"testing"

	"github.com/llgcode/draw2d/draw2dimg"
)

func TestComposeImage(t *testing.T) {
	// Initialize the graphic context on an RGBA image
	dest := image.NewRGBA(image.Rect(0, 0, 374, 572))
	for x := 0; x < 375; x++ {
		for y := 0; y < 573; y++ {
			dest.Set(x, y, color.Black)
		}
	}
	gc := draw2dimg.NewGraphicContext(dest)

	// Set some properties
	gc.SetFillColor(color.RGBA{0x44, 0xff, 0x44, 0xff})
	gc.SetStrokeColor(color.RGBA{0x44, 0x44, 0x44, 0xff})
	gc.SetLineWidth(5)

	// Draw a closed shape
	gc.MoveTo(10, 10) // should always be called first for a new path
	gc.LineTo(100, 50)
	gc.QuadCurveTo(100, 10, 10, 10)
	gc.Close()
	gc.FillStroke()

	r, err := os.Open("code.png")
	if err != nil {
		t.Fatal(err)
	}
	codeImg, err := png.Decode(r)
	if err != nil {
		t.Fatal(err)
	}

	gc.Translate(0, 0)
	gc.DrawImage(codeImg)

	// Save to file
	draw2dimg.SaveToPngFile("hello.png", dest)
}
