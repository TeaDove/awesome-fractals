package brot

import (
	"context"
	"image"
	"image/color"
	"math"
	"runtime"
	"sync"

	"golang.org/x/sync/semaphore"

	"github.com/teadove/awesome-fractals/internal/palette"
)

type BuilderInput struct {
	Colors        []palette.Color
	ColorStep     int
	XPos          float64
	YPos          float64
	Width, Height int
	MaxIteration  int
	EscapeRadius  float64
	Done          chan struct{}
}

type Input struct {
	Done          chan struct{}
	XPos, YPos    float64
	Width, Height int
	MaxIteration  int
	EscapeRadius  float64
	Colors        []color.RGBA
}

func NewInputFromBuilderInput(input *BuilderInput) (*Input, int) {
	r := Input{
		XPos:         input.XPos,
		YPos:         input.YPos,
		Width:        input.Width,
		Height:       input.Height,
		MaxIteration: input.MaxIteration,
		EscapeRadius: input.EscapeRadius,
		Done:         input.Done,
	}

	if input.ColorStep < input.MaxIteration {
		input.ColorStep = input.MaxIteration
	}

	r.Colors = palette.InterpolateColors(input.Colors, float64(input.ColorStep))

	return &r, r.Height
}

func Render(input *Input) *image.RGBA {
	ratio := float64(input.Height) / float64(input.Width)
	xmin, xmax := input.XPos-input.EscapeRadius/2.0, math.Abs(input.XPos+input.EscapeRadius/2.0)
	ymin, ymax := input.YPos-input.EscapeRadius*ratio/2.0, math.Abs(
		input.YPos+input.EscapeRadius*ratio/2.0,
	)

	rgbaImageComplied := image.NewRGBA(
		image.Rectangle{Min: image.Point{}, Max: image.Point{X: input.Width, Y: input.Height}},
	)

	var (
		weighted = semaphore.NewWeighted(int64(runtime.GOMAXPROCS(0) + 2))
		wg       sync.WaitGroup
	)

	for iy := range input.Height {
		_ = weighted.Acquire(context.Background(), 1)

		wg.Go(func() {
			defer weighted.Release(1)

			defer func() {
				if input.Done != nil {
					input.Done <- struct{}{}
				}
			}()

			for ix := range input.Width {
				x := xmin + (xmax-xmin)*float64(ix)/float64(input.Width-1)
				y := ymin + (ymax-ymin)*float64(iy)/float64(input.Width-1)
				norm, it := mandelIteration(x, y, input.MaxIteration)
				iteration := float64(input.MaxIteration-it) + math.Log(norm)

				if int(math.Abs(iteration)) >= len(input.Colors)-1 {
					continue
				}

				color1 := input.Colors[int(math.Abs(iteration))]
				color2 := input.Colors[int(math.Abs(iteration))+1]
				compiledColor := linearInterpolation(
					rgbaToUint(color1),
					rgbaToUint(color2),
					uint32(iteration),
				)

				rgbaImageComplied.Set(ix, iy, uint32ToRgba(compiledColor))
			}
		})
	}

	wg.Wait()

	return rgbaImageComplied
}

func linearInterpolation(c1, c2, mu uint32) uint32 {
	return c1*(1-mu) + c2*mu
}

func mandelIteration(cx, cy float64, maxIter int) (float64, int) {
	var x, y, xx, yy float64

	for i := range maxIter {
		xy := x * y
		xx = x * x

		yy = y * y
		if xx+yy > 4 {
			return xx + yy, i
		}

		x = xx - yy + cx
		y = 2*xy + cy
	}

	logZn := (x*x + y*y) / 2

	return logZn, maxIter
}

func rgbaToUint(color color.RGBA) uint32 {
	r, g, b, a := color.RGBA()
	r /= 0xff
	g /= 0xff
	b /= 0xff
	a /= 0xff

	return r<<24 | g<<16 | b<<8 | a
}

func uint32ToRgba(col uint32) color.RGBA {
	r := col >> 24 & 0xff
	g := col >> 16 & 0xff
	b := col >> 8 & 0xff
	a := 0xff

	return color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
}
