package electricity

import (
	"image"
	"image/color"
	"slices"

	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/llgcode/draw2d/draw2dkit"
)

var COLOR_BLACK = color.RGBA{0x00, 0x00, 0x00, 0xff}
var COLOR_WHITE = color.RGBA{0xff, 0xff, 0xff, 0xff}

func drawBars(prices []float64, width int, height int) image.Image {
	// Initialize the graphic context on an RGBA image
	canvas := image.Rect(0, 0, width, height)
	dest := image.NewRGBA(canvas)
	gc := draw2dimg.NewGraphicContext(dest)

	// Convenience
	w := float64(width)
	h := float64(height)

	// Draw background
	gc.Save()
	draw2dkit.Rectangle(gc, 0, 0, w, h)
	gc.SetFillColor(COLOR_WHITE)
	gc.FillStroke()
	gc.Restore()

	// Draw bars
	n := len(prices)

	xSpacing := 4.0
	ySpacing := 10.0

	maxBarHeight := h - ySpacing*2
	barWidth := (w - float64(n+1)*xSpacing) / float64(n)

	maxPrice := slices.Max(prices)

	for i, price := range prices {
		// TODO: Use delta instead?
		barHeight := (price / maxPrice) * maxBarHeight

		x0 := float64(i)*xSpacing + xSpacing + float64(i)*barWidth
		x1 := x0 + barWidth
		y0 := h - ySpacing
		y1 := y0 - barHeight

		gc.Save()
		draw2dkit.Rectangle(gc, x0, y0, x1, y1)
		gc.SetFillColor(COLOR_BLACK)
		gc.FillStroke()
		gc.Restore()

	}

	// Save to file
	draw2dimg.SaveToPngFile("hello.png", dest)
	return dest
}
