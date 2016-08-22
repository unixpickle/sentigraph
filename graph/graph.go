package main

import (
	"image"
	"image/color"
	"sort"

	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/unixpickle/sentigraph"
)

const (
	ImageWidth  = 500
	ImageHeight = 200
	PointCount  = 100
	LineWidth   = 2
)

func graph(points []*DataPoint) image.Image {
	sort.Sort(PointSorter(points))

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
		x := float64(i) * ImageWidth / PointCount
		if x == 0 {
			ctx.MoveTo(0, ImageHeight/2-y*(ImageHeight/2))
		} else {
			ctx.LineTo(float64(x), ImageHeight/2-y*(ImageHeight/2))
		}
	}
	ctx.Stroke()

	return res
}

type PointSorter []*DataPoint

func (p PointSorter) Len() int {
	return len(p)
}

func (p PointSorter) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p PointSorter) Less(i, j int) bool {
	return p[i].Position < p[j].Position
}
