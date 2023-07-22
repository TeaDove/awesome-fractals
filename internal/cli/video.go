package cli

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"runtime"

	"github.com/icza/mjpeg"
	"github.com/rs/zerolog/log"
	"github.com/schollz/progressbar/v3"
	"github.com/teadove/awesome-fractals/internal/brot"
	"github.com/teadove/awesome-fractals/internal/palette"

	"github.com/urfave/cli/v2"
)

func getStepFloat64(start, end float64, count int) float64 {
	return (end - start) / float64(count)
}

func getStepInt(start, end, count int) int {
	return (end - start) / count
}

func drawVideo(cCtx *cli.Context) error {
	runtime.GOMAXPROCS(cCtx.Int(maxprocsFlag.Name))

	colors, ok := palette.ColorPalettes[cCtx.String(paletteFlag.Name)]
	if !ok {
		return cli.Exit("Palette not found", 1)
	}

	xPos, xPosEnd := cCtx.Float64(xposFlag.Name), cCtx.Float64(xposEndFlag.Name)
	yPos, yPosEnd := cCtx.Float64(yposFlag.Name), cCtx.Float64(yposEndFlag.Name)
	colorStep, _ := cCtx.Int(stepFlag.Name), cCtx.Int(stepEndFlag.Name)
	escapeRadius, escapeRadiusEnd := cCtx.Float64(radiusFlag.Name), cCtx.Float64(radiusEndFlag.Name)
	iteration, iterationEnd := cCtx.Int(iterationFlag.Name), cCtx.Int(iterationEndFlag.Name)

	buildInput := brot.BuilderInput{
		Colors:       colors,
		ColorStep:    colorStep,
		XPos:         xPos,
		YPos:         yPos,
		Width:        cCtx.Int(widthFlag.Name),
		Height:       cCtx.Int(heightFlag.Name),
		MaxIteration: iteration,
		EscapeRadius: escapeRadius,
	}

	frames := cCtx.Int(secondsFlag.Name) * cCtx.Int(fpsFlag.Name)
	input, _ := brot.NewInputFromBuilderInput(&buildInput)

	aw, err := mjpeg.New(
		cCtx.String(videoFilenameFlag.Name),
		int32(cCtx.Int(widthFlag.Name)),
		int32(cCtx.Int(heightFlag.Name)),
		int32(cCtx.Int(fpsFlag.Name)),
	)
	if err != nil {
		return cli.Exit(fmt.Sprintf("Error creating video file: %s", err.Error()), 1)
	}

	xposStep := getStepFloat64(xPos, xPosEnd, frames)
	yposStep := getStepFloat64(yPos, yPosEnd, frames)
	escapeRadiusStep := getStepFloat64(escapeRadius, escapeRadiusEnd, frames)
	maxIterationStep := getStepInt(iteration, iterationEnd, frames)
	log.Info().Str("status", "start.compiling.video").
		Float64("xpos_step", xposStep).
		Float64("ypos_step", yposStep).
		Float64("escape_radius_step", escapeRadiusStep).
		Int("max_iteration_step", maxIterationStep).Send()
	bar := progressbar.Default(int64(frames))
	for i := 0; i < frames; i++ {
		input.XPos += xposStep
		input.YPos += yposStep
		input.EscapeRadius += escapeRadiusStep
		input.MaxIteration += maxIterationStep

		img := brot.Render(input)

		buf := &bytes.Buffer{}
		err = jpeg.Encode(buf, img, nil)
		if err != nil {
			return cli.Exit(fmt.Sprintf("Error encoding jpeg: %s", err.Error()), 1)
		}
		err = aw.AddFrame(buf.Bytes())
		if err != nil {
			return cli.Exit(fmt.Sprintf("Error adding frame: %s", err.Error()), 1)
		}

		err := bar.Add(1)
		if err != nil {
			log.Error().Err(err).Str("status", "bar.add.error")
		}
	}
	err = aw.Close()
	if err != nil {
		return cli.Exit(fmt.Sprintf("Error closing video file: %s", err.Error()), 1)
	}
	return nil
}
