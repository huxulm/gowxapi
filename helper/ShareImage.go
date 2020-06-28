package helper

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"math"
	"os/exec"

	"github.com/jackdon/gowxapi/config"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/nfnt/resize"
)

// LOGO is a bytes holder of logo img
var LOGO []byte

func init() {
	l := config.C.StaticResource.LogoPath
	LOGO, _ = LoadLogo(l)
}

// LoadLogo loads logo into bytes from a given path
func LoadLogo(path string) ([]byte, error) {
	if b, err := ioutil.ReadFile(path); err == nil {
		return b, nil
	} else {
		return nil, err
	}
}

func GenHTMLImage(src []byte, options *map[string]string) ([]byte, error) {
	_, err := exec.LookPath("wkhtmltoimage")
	if err != nil {
		return nil, err
	}
	opts := make([]string, 0)
	if options != nil {
		for key, opt := range *options {
			opts = append(opts, fmt.Sprintf("--%s", key), opt)
		}
	}
	// set input and output use stdin & stdout
	opts = append(opts, "-", "-")
	cmd := exec.Command("wkhtmltoimage", opts...)
	cmd.Stdin = bytes.NewReader(src)

	stdout, err := cmd.StdoutPipe()
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)

	if _, err := io.Copy(buf, stdout); err == nil {
		// return buf.Bytes(), nil
	}

	if err := cmd.Wait(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// ShareImage is an interface
type ShareImage interface{}

// Generate composites html, avatar and logo image bytes and
// return a pointer of image.Image
func Generate(code, avatar, logo []byte, userInfo map[string]string) ([]byte, error) {
	bgImgR := bytes.NewReader(code)
	var rect image.Rectangle
	if bgConfig, err := png.DecodeConfig(bgImgR); err == nil {
		rect = image.Rect(0, 0, 380, 200+bgConfig.Height)
	} else {
		rect = image.Rect(0, 0, 380, 600)
	}
	dest := image.NewRGBA(rect)
	draw2d.SetFontFolder(config.C.StaticResource.FontsPath)
	gc := draw2dimg.NewGraphicContext(dest)
	gc.SetFontData(draw2d.FontData{"AiNiChuanYueRenHai", 4, draw2d.FontStyleNormal})
	for x := 0; x <= 380; x++ {
		for y := 0; y <= 580; y++ {
			dest.Set(x, y, color.White)
		}
	}
	if bgImg, err := png.Decode(bgImgR); err == nil {
		gc.Translate(0, 160)
		gc.DrawImage(bgImg)
		gc.Translate(0, -160)
	}

	logoImgR := bytes.NewReader(logo)
	if logImg, err := png.Decode(logoImgR); err == nil {
		logImg = resize.Resize(80, 80, logImg, resize.MitchellNetravali)
		gc.Translate(280, 500)
		gc.DrawImage(cropImage(logImg, 80))
		gc.Translate(-280, -500)
	}

	avatarImgR := bytes.NewReader(avatar)
	if avatarImg, err := jpeg.Decode(avatarImgR); err == nil {
		// resize
		avatarImg = resize.Resize(60, 60, avatarImg, resize.MitchellNetravali)
		gc.Translate(20, 20)
		gc.DrawImage(cropImage(avatarImg, 80))

		gc.Translate(-20, -20)

		gc.Translate(120, 40)
		gc.SetFontSize(16)
		gc.SetFillColor(color.Black)
		gc.FillStringAt(userInfo["nick"], 0, 25)
		gc.Translate(-120, -40)

		gc.Translate(120, 76)
		gc.SetFontSize(12)
		gc.SetFillColor(color.Black)
		gc.FillStringAt("邀你一起Go golang~", 0, 25)
		gc.Translate(-120, -76)
	}

	w := new(bytes.Buffer)
	if err := png.Encode(w, dest); err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

func cropImage(src image.Image, size int) *image.RGBA {
	width, height := src.Bounds().Dx(), src.Bounds().Dy()
	if size > width && size > height {
		size = int(math.Min(float64(width), float64(height)))
	}
	posX := (width - size) / 2
	posY := (height - size) / 2

	var clacDimension = func(x, y, cx, cy, r int) float64 {
		return math.Sqrt(float64((cx-x)*(cx-x) + (cy-y)*(cy-y)))
	}

	ci := image.NewRGBA(image.Rect(0, 0, size, size))
	for x := posX; x < (posX + size); x++ {
		for y := posY; y < (posY + size); y++ {
			if d := clacDimension(x, y, posX+size/2, posY+size/2, size/2); float64(size/2) >= d {
				ci.Set(x, y, src.At(x, y))
			}
		}
	}
	return ci
}
