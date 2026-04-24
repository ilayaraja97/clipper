package main

import (
	"fmt"
	"os"
	"time"

	"github.com/ilayaraja97/clipper/logger"
	"github.com/ilayaraja97/clipper/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	logger.Init(time.Now().Format("20060102-150405"))
	defer func() {
		if r := recover(); r != nil {
			logger.Log.Fatal().Interface("panic", r).Msg("application panicked")
			fmt.Fprintf(os.Stderr, "\nApplication panicked: %v\n%s\n", r, logger.GetPanicMessage())
			os.Exit(1)
		}
	}()

	logger.Log.Info().Msg("clipper starting")

	input, err := ui.NewUIInput()
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("failed to create UI input")
	}

	if _, err := tea.NewProgram(ui.NewUi(input)).Run(); err != nil {
		logger.Log.Fatal().Err(err).Msg("application error")
	}

	logger.Log.Info().Msg("clipper exiting")
}
