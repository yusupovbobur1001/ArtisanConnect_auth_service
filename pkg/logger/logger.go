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

// func (h *PrettyHandler) Handle(ctx context.Context, r slog.Record) error {

// 	level := r.Level.String() + ":"

// 	switch r.Level {
// 	case slog.LevelDebug:
// 		level = color.MagentaString(level)
// 	case slog.LevelInfo:
// 		level = color.BlueString(level)
// 	case slog.LevelWarn:
// 		level = color.YellowString(level)
// 	case slog.LevelError:
// 		level = color.RedString(level)
// 	}

// 	feilds := make(map[string]interface{}, r.NumAttrs())
// 	r.Attrs(func(a slog.Attr) bool {
// 		feilds[a.Key] = a.Value.Any()

// 		return true
// 	})

// 	b, err := json.MarshalIndent(feilds, "", "  ")
// 	if err != nil {
// 		return err
// 	}

// 	timeStr := r.Time.Format("[15:05:05.000]")
// 	msg := color.CyanString(r.Message)

// 	h.l.Println(timeStr, level, msg, color.WhiteString(string(b)))

// 	return nil
// }

func NewPrettyHandler(out io.Writer, opts PrettyHandlerOptions) *PrettyHandler {
	h := &PrettyHandler{
		Handler: slog.NewJSONHandler(out, &opts.SlogOpts),
		l:       log.New(out, "", 0),
	}

	return h
}
