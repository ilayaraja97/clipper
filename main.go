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
			// zerolog's Fatal() calls os.Exit(1) from within Msg(), which would
			// skip the message below. The logger writes to a file, so stderr is
			// the only place the user actually sees this.
			logger.Log.Error().Interface("panic", r).Msg("application panicked")
			fmt.Fprintf(os.Stderr, "\nApplication panicked: %v\n%s\n", r, logger.GetPanicMessage())
			os.Exit(1)
		}
	}()

	logger.Log.Info().Msg("clipper starting")

	input, err := ui.NewUIInput()
	if err != nil {
		fail("failed to read input", err)
	}

	if _, err := tea.NewProgram(ui.NewUi(input)).Run(); err != nil {
		fail("application error", err)
	}

	logger.Log.Info().Msg("clipper exiting")
}

func fail(message string, err error) {
	logger.Log.Error().Err(err).Msg(message)
	fmt.Fprintf(os.Stderr, "\n%s: %v\n%s\n", message, err, logger.GetPanicMessage())
	os.Exit(1)
}
