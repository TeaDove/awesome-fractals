package palette

import (
	"image/color"
	"math"
)

func InterpolateColors(colors []Color, colorStep float64) []color.RGBA {
	var (
		factor             float64
		steps              []float64
		cols               []uint32
		interpolated       []uint32
		interpolatedColors []color.RGBA
	)

	factor = 1.0 / colorStep

	for index, col := range colors {
		if col.Step == 0.0 && index != 0 {
			stepRatio := float64(index+1) / float64(len(colors))
			step := float64(int(stepRatio*100)) / 100 // truncate to 2 decimal precision
			steps = append(steps, step)
		} else {
			steps = append(steps, col.Step)
		}

		r, g, b, a := col.Color.RGBA()
		r /= 0xff
		g /= 0xff
		b /= 0xff
		a /= 0xff
		uintColor := r<<24 | g<<16 | b<<8 | a
		cols = append(cols, uintColor)
	}

	var minStep, maxStep, minColor, maxColor float64

	for i := 0.0; i <= 1; i += factor {
		for j := range len(colors) - 1 {
			if !(i >= steps[j] && i < steps[j+1]) {
				continue
			}

			minStep = steps[j]
			maxStep = steps[j+1]
			minColor = float64(cols[j])
			maxColor = float64(cols[j+1])
			uintColor := cosineInterpolation(
				maxColor,
				minColor,
				(i-minStep)/(maxStep-minStep),
			)
			interpolated = append(interpolated, uint32(uintColor))
		}
	}

	for _, pixelValue := range interpolated {
		r := pixelValue >> 24 & 0xff
		g := pixelValue >> 16 & 0xff
		b := pixelValue >> 8 & 0xff
		a := 0xff

		interpolatedColors = append(
			interpolatedColors,
			color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)},
		)
	}

	return interpolatedColors
}

func cosineInterpolation(c1, c2, mu float64) float64 {
	mu2 := (1 - math.Cos(mu*math.Pi)) / 2.0
	return c1*(1-mu2) + c2*mu2
}
