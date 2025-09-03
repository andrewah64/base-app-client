package log

import (
	"context"
	"os"
	"log/slog"
	"runtime/debug"
)

import (
	"github.com/andrewah64/base-app-client/internal/common/core/session"
)

type CtxDataHandler struct {
	slog.Handler
}

func (h *CtxDataHandler) Handle(ctx context.Context, r slog.Record) error {
	c := r.Clone()

	if mwd, ok := session.FromContext(ctx); ok {
		c.AddAttrs(slog.String("requestId" , mwd.RequestId))
	}

	return h.Handler.Handle(ctx, c)
}

func (h *CtxDataHandler) WithGroup(name string) slog.Handler {
	return &CtxDataHandler{h.Handler.WithGroup(name)}
}

func (h *CtxDataHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &CtxDataHandler{h.Handler.WithAttrs(attrs)}
}

func Setup (level slog.Level) *slog.Logger {
	bi, _ := debug.ReadBuildInfo()

	opts := &slog.HandlerOptions {
		AddSource: true,
		Level    : level,
	}

	h := &CtxDataHandler{
		slog.NewJSONHandler(os.Stdout, opts),
	}

	logger := slog.New(h).With(
		slog.Group("program",
			slog.Int("os-pid", os.Getpid()),
			slog.String("go-version", bi.GoVersion),
		),
	)

	return logger.WithGroup("request")
}
