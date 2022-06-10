package lib

import (
	"image"
	"image/png"
	"io/ioutil"

	"github.com/kbinani/screenshot"
	"gocv.io/x/gocv"
)

// ScreenshotBounds struct contains details about the screenshot
type ScreenshotBounds struct {
	x      int
	y      int
	width  int
	height int
}
type ImageUtils struct{}

var ImageTools ImageUtils

// IsImageWithin determines if a small image is withing a big image
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

// SystrayShot takes a screenshot of the systray
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

// SaveAsTempPNG saves an image as a temp png
func (i *ImageUtils) SaveAsTempPNG(img *image.RGBA) (string, error) {
	file, err := ioutil.TempFile("", "utiligo-")
	if err != nil {
		return "", err
	}
	defer file.Close()
	err = png.Encode(file, img)

	if err != nil {
		return "", err
	}

	return file.Name(), nil
}

// Screenshot takes a screenshot using ScreenshotBounds
func (i *ImageUtils) Screenshot(b *ScreenshotBounds) (*image.RGBA, error) {
	return screenshot.Capture(b.x, b.y, b.width, b.height)
}
