package cli

import (
	"flag"
	"os"
	"runtime/pprof"
	"sync"

	"github.com/schollz/progressbar/v3"
	"github.com/teadove/awesome-fractals/internal/brot"
)

var service brot.Service

func init() {
	service.WG = &sync.WaitGroup{}

	flag.Float64Var(
		&service.ColorStep,
		"step",
		6000,
		"Color smooth step. Value should be greater than iteration count, otherwise the value will be adjusted to the iteration count.",
	)
	flag.IntVar(&service.Width, "width", 1000, "Rendered image width")
	flag.IntVar(&service.Height, "height", 1000, "Rendered image height")
	flag.Float64Var(
		&service.XPos,
		"xpos",
		-0.00275,
		"Point position on the real axis (defined on `x` axis)",
	)
	flag.Float64Var(
		&service.YPos,
		"ypos",
		0.78912,
		"Point position on the imaginary axis (defined on `y` axis)",
	)
	flag.Float64Var(&service.EscapeRadius, "radius", .125689, "Escape Radius")
	flag.IntVar(&service.MaxIteration, "iteration", 800, "Iteration count")
	flag.StringVar(
		&service.ColorPalette,
		"palette",
		"Hippi",
		"Hippi | Plan9 | AfternoonBlue | SummerBeach | Biochimist | Fiesta",
	)
	flag.StringVar(
		&service.OutputFile,
		"file",
		"mandelbrot.png",
		"The rendered mandelbrot image filname",
	)
	flag.Parse()
}

func Run() {
	file, err := os.OpenFile("main.prof", os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	err = pprof.StartCPUProfile(file)
	if err != nil {
		panic(err)
	}

	done := make(chan struct{})
	iterations := service.Init()

	go func() {
		bar := progressbar.Default(int64(iterations))
		for i := 0; i <= iterations; i++ {
			<-done
			err := bar.Add(1)
			if err != nil {
				println(err.Error())
			}
		}
	}()

	service.Run(done)

	pprof.StopCPUProfile()
}
