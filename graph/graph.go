package main

import (
	"image"
	"image/color"

	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/unixpickle/sentigraph"
)

const (
	ImageWidth  = 500
	ImageHeight = 200
	PointCount  = 70
	LineWidth   = 2
)

func graph(points []*DataPoint) image.Image {
	yMean := make([]float64, PointCount)
	yCount := make([]float64, PointCount)
	for _, point := range points {
		xVal := int(point.Position * PointCount)
		if xVal == PointCount {
			xVal = PointCount - 1
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

	res := image.NewRGBA(image.Rect(0, 0, ImageWidth, ImageHeight))
	ctx := draw2dimg.NewGraphicContext(res)

	ctx.SetStrokeColor(color.RGBA{A: 0xff})
	ctx.SetLineWidth(LineWidth)

	ctx.BeginPath()
	for i, y := range yMean {
		x := float64(i) * ImageWidth / (PointCount - 1)
		if x == 0 {
			ctx.MoveTo(0, ImageHeight/2-y*(ImageHeight/2))
		} else {
			ctx.LineTo(float64(x), ImageHeight/2-y*(ImageHeight/2))
		}
	}
	ctx.Stroke()

	return res
}
