package logger

import (
	"io"
	"log"
	"log/slog"
	"os"
)

type PrettyHandlerOptions struct {
	SlogOpts slog.HandlerOptions
}

type PrettyHandler struct {
	slog.Handler
	l *log.Logger
}

func New() (*slog.Logger, error) {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}

	file, err := os.OpenFile("logger.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}

	logger := slog.New(slog.NewJSONHandler(file, opts))
	log.Println("Logger initialized successfully")
	return logger, nil
}



func NewPrettyHandler(out io.Writer, opts PrettyHandlerOptions) *PrettyHandler {
	h := &PrettyHandler{
		Handler: slog.NewJSONHandler(out, &opts.SlogOpts),
		l:       log.New(out, "", 0),
	}

	return h
}
