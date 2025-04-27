package logger

import (
	"os"

	"github.com/rs/zerolog"
)

// Initialize the logger
var Log zerolog.Logger

func Init() {
	// Set the global logger to output JSON logs by default
	Log = zerolog.New(os.Stdout).With().Timestamp().Logger()
}

// Utility function for easier logging
func Info(message string) {
	Log.Info().Msg(message)
}

func Error(message string, err error) {
	Log.Error().Err(err).Msg(message)
}
