package main

import (
	"image"
	"image/color"
	"math"

	"github.com/llgcode/draw2d/draw2dimg"
)

const (
	HeatPointCount  = 70
	HeatImageWidth  = 400
	HeatImageHeight = 100
)

func heatGraph(d []*DataPoint) image.Image {
	points := lineDataPoints(d, HeatPointCount)

	var mean float64
	var variance float64
	for _, p := range points {
		mean += p
		variance += p * p
	}
	mean /= float64(len(points))
	variance /= float64(len(points))
	variance -= mean * mean
	stddev := math.Sqrt(variance)

	res := image.NewRGBA(image.Rect(0, 0, HeatImageWidth, HeatImageHeight))
	ctx := draw2dimg.NewGraphicContext(res)

	for i, p := range points {
		normalized := (p - mean) / stddev
		intensity := math.Tanh(normalized)
		if intensity > 0 {
			ctx.SetFillColor(color.RGBA{
				R: uint8(0xff*(1-intensity) + 0.5),
				G: 0xff,
				B: uint8(0xff*(1-intensity) + 0.5),
				A: 0xff,
			})
		} else {
			ctx.SetFillColor(color.RGBA{
				R: 0xff,
				G: uint8(0xff*(1+intensity) + 0.5),
				B: uint8(0xff*(1+intensity) + 0.5),
				A: 0xff,
			})
		}
		x := float64(i) * HeatImageWidth / HeatPointCount
		nextX := float64(i+1) * HeatImageWidth / HeatPointCount
		ctx.BeginPath()
		ctx.MoveTo(x, 0)
		ctx.LineTo(nextX, 0)
		ctx.LineTo(nextX, HeatImageHeight)
		ctx.LineTo(x, HeatImageHeight)
		ctx.Close()
		ctx.Fill()
	}

	return res
}
