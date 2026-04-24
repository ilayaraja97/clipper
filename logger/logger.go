package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"

	"github.com/rs/zerolog"
)

var (
	Log         zerolog.Logger
	logFilePath string
)

func Init(sessionID string) {
	logDir := filepath.Join(getLogDir(), "clipper")
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "failed to create log directory: %v\n", err)
		InitToStderr()
		return
	}

	logFilePath = filepath.Join(logDir, fmt.Sprintf("clipper-%s.log", sessionID))
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open log file: %v\n", err)
		InitToStderr()
		return
	}

	multiWriter := io.MultiWriter(file)
	Log = zerolog.New(multiWriter).
		With().
		Timestamp().
		Caller().
		Logger()
}

func InitToStderr() {
	Log = zerolog.New(os.Stderr).
		With().
		Timestamp().
		Caller().
		Logger()
	logFilePath = ""
}

func GetLogFilePath() string {
	return logFilePath
}

func getLogDir() string {
	if runtime.GOOS == "windows" {
		return os.Getenv("LOCALAPPDATA")
	}
	return filepath.Join(os.Getenv("HOME"), ".local", "share")
}

func GetPanicMessage() string {
	if logFilePath == "" {
		return "check stderr for details"
	}
	return fmt.Sprintf("check log file: %s", logFilePath)
}
