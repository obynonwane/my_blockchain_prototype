// /logger/logger.go
package logger

import (
	"os"

	"github.com/rs/zerolog"
)

// Global logger instance
var Log zerolog.Logger

// Initialize the logger (called once in main or server.go)
func Init() {
	Log = zerolog.New(os.Stdout).With().Timestamp().Logger()
}

// Utility functions to simplify logging
func Info(message string) {
	Log.Info().Msg(message)
}

func Error(message string, err error) {
	Log.Error().Err(err).Msg(message)
}

func Debug(message string) {
	Log.Debug().Msg(message)
}
