package main

import (
	"image"
	"image/color"

	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/unixpickle/sentigraph"
)

const (
	LineImageWidth  = 500
	LineImageHeight = 200
	LinePointCount  = 70
	LineWidth       = 2
)

func lineGraph(points []*DataPoint) image.Image {
	res := image.NewRGBA(image.Rect(0, 0, LineImageWidth, LineImageHeight))
	ctx := draw2dimg.NewGraphicContext(res)

	ctx.SetStrokeColor(color.RGBA{A: 0xff})
	ctx.SetLineWidth(LineWidth)

	ctx.BeginPath()
	for i, y := range lineDataPoints(points, LinePointCount) {
		x := float64(i) * LineImageWidth / (LinePointCount - 1)
		if x == 0 {
			ctx.MoveTo(0, LineImageHeight/2-y*(LineImageHeight/2))
		} else {
			ctx.LineTo(float64(x), LineImageHeight/2-y*(LineImageHeight/2))
		}
	}
	ctx.Stroke()

	return res
}

func lineDataPoints(points []*DataPoint, count int) []float64 {
	yMean := make([]float64, count)
	yCount := make([]float64, count)
	for _, point := range points {
		xVal := int(point.Position * float64(count))
		if xVal == count {
			xVal = count - 1
		}
		switch point.Sentiment {
		case sentigraph.Positive:
			yMean[xVal]++
		case sentigraph.Negative:
			yMean[xVal]--
		}
		yCount[xVal]++
	}
	for i, c := range yCount {
		yMean[i] /= c
	}
	return yMean
}
