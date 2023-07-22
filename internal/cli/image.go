package cli

import (
	"image/png"
	"os"
	"runtime"

	"github.com/rs/zerolog/log"
	"github.com/schollz/progressbar/v3"

	"github.com/teadove/awesome-fractals/internal/brot"
	"github.com/teadove/awesome-fractals/internal/palette"
	"github.com/urfave/cli/v2"
)

func drawImage(cCtx *cli.Context) error {
	runtime.GOMAXPROCS(cCtx.Int(maxprocsFlag.Name))

	colors, ok := palette.ColorPalettes[cCtx.String(paletteFlag.Name)]
	if !ok {
		return cli.Exit("Palette not found", 1)
	}

	done := make(chan struct{})
	buildInput := brot.BuilderInput{
		Colors:       colors,
		ColorStep:    cCtx.Int(stepFlag.Name),
		XPos:         cCtx.Float64(xposFlag.Name),
		YPos:         cCtx.Float64(yposFlag.Name),
		Width:        cCtx.Int(widthFlag.Name),
		Height:       cCtx.Int(heightFlag.Name),
		MaxIteration: cCtx.Int(iterationFlag.Name),
		EscapeRadius: cCtx.Float64(radiusFlag.Name),
		Done:         done,
	}

	input, iterations := brot.NewInputFromBuilderInput(&buildInput)
	go func() {
		bar := progressbar.Default(int64(iterations))
		for i := 0; i < iterations; i++ {
			<-done
			err := bar.Add(1)
			if err != nil {
				log.Error().Err(err).Str("status", "bar.add.error")
			}
		}
	}()
	img := brot.Render(input)

	output, err := os.Create(cCtx.String(imageFilenameFlag.Name))
	if err != nil {
		return err
	}
	err = png.Encode(output, img)
	if err != nil {
		return err
	}

	return nil
}
