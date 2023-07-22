package cli

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/teadove/awesome-fractals/internal/palette"
	"github.com/urfave/cli/v2"
)

var (
	widthFlag  = &cli.IntFlag{Name: "width", Usage: "Rendered image width", Value: 1_000}
	heightFlag = &cli.IntFlag{Name: "height", Usage: "Rendered image height", Value: 1_000}
	stepFlag   = &cli.IntFlag{
		Name: "step",
		Usage: "Color smooth step. " +
			"Value should be greater than iteration count, " +
			"otherwise the value will be adjusted to the iteration count.",
		Value: 6_000,
	}
	stepEndFlag = &cli.IntFlag{
		Name:  "step-end",
		Usage: "step end value",
		Value: 6_000,
	}
	xposFlag = &cli.Float64Flag{
		Name:  "xpos",
		Usage: "Point position on the real axis (defined on `x` axis)",
		Value: -0.00275,
	}
	xposEndFlag = &cli.Float64Flag{
		Name:  "xpos-end",
		Usage: "xpos end value",
		Value: -0.00275,
	}
	yposFlag = &cli.Float64Flag{
		Name:  "ypos",
		Usage: "Point position on the imaginary axis (defined on `y` axis)",
		Value: 0.78912,
	}
	yposEndFlag = &cli.Float64Flag{
		Name:  "ypos-end",
		Usage: "ypos end value",
		Value: 0.78912,
	}
	radiusFlag    = &cli.Float64Flag{Name: "radius", Usage: "Escape Radius", Value: .125689}
	radiusEndFlag = &cli.Float64Flag{
		Name:  "radius-end",
		Usage: "radius end value",
		Value: .125689,
	}
	iterationFlag    = &cli.IntFlag{Name: "iteration", Value: 800, Usage: "Iteration count"}
	iterationEndFlag = &cli.IntFlag{Name: "iteration-end", Value: 800, Usage: "iteration end value"}
	paletteFlag      = &cli.StringFlag{
		Name:  "palette",
		Value: "Hippi",
		Usage: strings.Join(palette.GetPaletteNames(), " | "),
	}
	imageFilenameFlag = &cli.StringFlag{
		Name:  "filename",
		Value: "fractal.png",
		Usage: "Path to save image",
	}
	videoFilenameFlag = &cli.StringFlag{
		Name:  "filename",
		Value: "fractal.avi",
		Usage: "Path to save image",
	}
	secondsFlag = &cli.IntFlag{
		Name:  "seconds",
		Value: 3,
		Usage: "Amount of to generate in video",
	}
	fpsFlag = &cli.IntFlag{
		Name:  "fps",
		Value: 30,
		Usage: "Frames per second",
	}

	maxprocsFlag = &cli.IntFlag{
		Name:  "maxprocs",
		Usage: "max amount of processes to use, by default is amount of cores in CPU munis one",
		Value: runtime.NumCPU() - 1,
	}
)

func Run() {
	captureInterrupt()

	app := &cli.App{
		Flags: []cli.Flag{maxprocsFlag},
		Commands: []*cli.Command{
			{
				Name:   "image",
				Action: drawImage,
				Flags: []cli.Flag{
					imageFilenameFlag,
					widthFlag,
					heightFlag,
					stepFlag,
					xposFlag,
					yposFlag,
					radiusFlag,
					iterationFlag,
					paletteFlag,
				},
			},
			{
				Name:   "video",
				Action: drawVideo,
				Flags: []cli.Flag{
					widthFlag,
					heightFlag,

					stepFlag,
					stepEndFlag,
					xposFlag,
					xposEndFlag,
					yposFlag,
					yposEndFlag,
					radiusFlag,
					radiusEndFlag,
					iterationFlag,
					iterationEndFlag,

					paletteFlag,
					secondsFlag,
					fpsFlag,
					videoFilenameFlag,
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}
}
