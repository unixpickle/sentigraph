package main

import (
	"image"
	"image/color"
	"math"

	"github.com/llgcode/draw2d/draw2dimg"
)

const (
	HeatPointCount  = 70
	HeatImageWidth  = 600
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

	xValues := make([]float64, 0, len(points))
	colors := make([]color.RGBA, 0, len(points))

	for i, p := range points {
		normalized := (p - mean) / stddev
		intensity := math.Tanh(normalized)
		if intensity > 0 {
			colors = append(colors, color.RGBA{
				R: uint8(0xff*(1-intensity) + 0.5),
				G: 0xff,
				B: uint8(0xff*(1-intensity) + 0.5),
				A: 0xff,
			})
		} else {
			colors = append(colors, color.RGBA{
				R: 0xff,
				G: uint8(0xff*(1+intensity) + 0.5),
				B: uint8(0xff*(1+intensity) + 0.5),
				A: 0xff,
			})
		}
		x := float64(i) * HeatImageWidth / HeatPointCount
		xValues = append(xValues, x)
	}

	var colorIdx int
	for x := 0.0; x < HeatImageWidth; x++ {
		for colorIdx < len(xValues)-1 && x > xValues[colorIdx+1] {
			colorIdx++
		}
		if colorIdx == len(colors)-1 {
			ctx.SetFillColor(colors[colorIdx])
		} else {
			first := colors[colorIdx]
			second := colors[colorIdx+1]
			t := (x - xValues[colorIdx]) / (xValues[colorIdx+1] - xValues[colorIdx])
			ctx.SetFillColor(color.RGBA{
				R: uint8(0.5 + float64(first.R)*(1-t) + float64(second.R)*t),
				G: uint8(0.5 + float64(first.G)*(1-t) + float64(second.G)*t),
				B: uint8(0.5 + float64(first.B)*(1-t) + float64(second.B)*t),
				A: 0xff,
			})
		}
		ctx.BeginPath()
		ctx.MoveTo(x, 0)
		ctx.LineTo(x+1, 0)
		ctx.LineTo(x+1, HeatImageHeight)
		ctx.LineTo(x, HeatImageHeight)
		ctx.Close()
		ctx.Fill()
	}

	return res
}
