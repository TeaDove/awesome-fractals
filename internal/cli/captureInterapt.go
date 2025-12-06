package cli

import (
	"os"
	"os/signal"
	"runtime/pprof"

	"github.com/rs/zerolog/log"
)

func captureInterrupt() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for sig := range c {
			log.Info().
				Str("status", "interrupt.captured.exiting").
				Str("signal", sig.String()).
				Send()
			pprof.StopCPUProfile()
			os.Exit(0)
		}
	}()
}
