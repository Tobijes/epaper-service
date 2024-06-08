package electricity

import (
	"image"
	"image/color"
	"strconv"

	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/llgcode/draw2d/draw2dkit"
)

var COLOR_BLACK = color.RGBA{0x00, 0x00, 0x00, 0xff}
var COLOR_WHITE = color.RGBA{0xff, 0xff, 0xff, 0xff}

func drawBars(prices []PriceRecord, width int, height int) image.Image {
	// Initialize the graphic context on an RGBA image
	draw2d.SetFontFolder("fonts")
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

	xSpacing := 8.0
	barWidth := (w - float64(n+1)*xSpacing) / float64(n)

	topPadding := 10.0
	bottomPadding := 5.0
	textSpace := 8.0 + bottomPadding
	maxBarHeight := h - (topPadding + bottomPadding + barWidth + textSpace) // Barwidth is 2*radius for rounded edge

	var maxPrice float64 = 0
	for _, record := range prices {
		if record.TotalDKK > maxPrice {
			maxPrice = record.TotalDKK
		}
	}

	for i, record := range prices {
		price := record.TotalDKK
		// TODO: Use delta instead?
		barHeight := (price / maxPrice) * maxBarHeight

		x0 := float64(i)*xSpacing + xSpacing + float64(i)*barWidth
		x1 := x0 + barWidth
		y0 := h - bottomPadding - barWidth/2 - textSpace
		y1 := y0 - barHeight

		// Draw bar
		gc.Save()
		draw2dkit.Rectangle(gc, x0, y0, x1, y1)
		gc.SetFillColor(COLOR_BLACK)
		gc.FillStroke()
		gc.Restore()

		// Draw rounded ends
		gc.Save()
		draw2dkit.Circle(gc, x0+barWidth/2, y0, barWidth/2)
		draw2dkit.Circle(gc, x0+barWidth/2, y1, barWidth/2)
		gc.SetFillColor(COLOR_BLACK)
		gc.FillStroke()
		gc.Restore()

		// Draw hour
		hour := strconv.Itoa(record.StartTimeUTC.In(location).Hour())
		gc.Save()
		gc.SetFillColor(COLOR_BLACK)
		gc.SetFontData(draw2d.FontData{Family: draw2d.FontFamilySans, Style: draw2d.FontStyleBold})
		gc.SetFontSize(8)
		left, _, right, _ := gc.GetStringBounds(hour)
		strWidth := right - left
		gc.FillStringAt(hour, x0+(barWidth/2)-(strWidth/2), h-bottomPadding)
		gc.Restore()

	}

	return dest
}
