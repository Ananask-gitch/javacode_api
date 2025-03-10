package logger

import (
	"io"
	"log/slog"
	"os"
)

var (
	Log     *slog.Logger
	logFile io.Writer
)

func Init(logFilePath string) error {
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	logFile = file

	jsonHandler := slog.NewJSONHandler(logFile, nil)
	Log = slog.New(jsonHandler).With(
		slog.String("app-version", "v0.0.1-beta"),
	)

	return nil
}

func Sync() error {
	if file, ok := logFile.(*os.File); ok {
		return file.Sync()
	}
	return nil
}
