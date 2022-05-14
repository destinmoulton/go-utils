package lib

import (
	"image"
	"image/png"
	"os"

	"github.com/kbinani/screenshot"
	"gocv.io/x/gocv"
)

// Bounds of a screenshot
type ScreenshotBounds struct {
	x      int
	y      int
	width  int
	height int
}
type ImageUtils struct{}

var ImageTools ImageUtils

// Is small pic within big pic?
func (i *ImageUtils) IsImageWithin(small string, big string) bool {
	img := gocv.IMRead(big, gocv.IMReadAnyColor)
	template := gocv.IMRead(small, gocv.IMReadAnyColor)
	matResult := gocv.NewMat()
	mask := gocv.NewMat()
	gocv.MatchTemplate(img, template, &matResult, gocv.TmCcorrNormed, mask)
	defer mask.Close()

	//minConfidence, maxConfidence, minLoc, maxLoc := gocv.MinMaxLoc(matResult)
	_, maxConfidence, _, _ := gocv.MinMaxLoc(matResult)
	defer matResult.Close()
	return maxConfidence > 0.95
}

// Take a screenshot of the systray
func (i *ImageUtils) SystrayShot(height int) (*image.RGBA, error) {
	xres := 1920
	yres := 1080
	b := &ScreenshotBounds{
		x:      0,
		y:      yres - height,
		width:  xres * 2,
		height: height,
	}
	return i.Screenshot(b)
}

// Save an image as a png
func (i *ImageUtils) SavePNG(img *image.RGBA, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return png.Encode(file, img)
}

// Take a screenshot using ScreenshotBounds
func (i *ImageUtils) Screenshot(b *ScreenshotBounds) (*image.RGBA, error) {
	return screenshot.Capture(b.x, b.y, b.width, b.height)
}
