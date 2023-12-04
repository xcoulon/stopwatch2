package cmd

import (
	"log/slog"
	"time"

	charmlog "github.com/charmbracelet/log"
)

func NewLogger(debug bool) *slog.Logger {
	charmlog.SetTimeFormat(time.Kitchen)
	if debug {
		charmlog.SetLevel(charmlog.DebugLevel)
	}
	return slog.New(charmlog.Default())
}
