package di

import (
	"log/slog"
	"os"
)

func (s *serviceImpl) logger() *slog.Logger {
	logHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	return slog.New(logHandler)
}
