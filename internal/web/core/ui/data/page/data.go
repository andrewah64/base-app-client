package page

import (
	"context"
	"log/slog"
	"slices"
	"time"
)

import (
	"github.com/andrewah64/base-app-client/internal/web/core/session"
)

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type Data struct {
	CSRFToken        string
	User             *session.AuthSessionUser
	FormOpts         *map[string]any
	NotificationData *map[string]any
	ResultSet        *map[string]any
	Localiser        *i18n.Localizer
}

type key int

var dataKey key

// NewContext returns a new Context that carries value u.
func NewContext(ctx context.Context, data *Data) context.Context {
	slog.LogAttrs(ctx, slog.LevelDebug, "data added to context")

	return context.WithValue(ctx, dataKey, data)
}

// FromContext returns the data value stored in ctx, if any.
func FromContext(ctx context.Context) (*Data, bool) {
	slog.LogAttrs(ctx, slog.LevelDebug, "get data from context")

	p, ok := ctx.Value(dataKey).(*Data)

	if ok {
		slog.LogAttrs(ctx, slog.LevelDebug, "data found in context")
	} else {
		slog.LogAttrs(ctx, slog.LevelError, "data not found in context")
	}

	return p, ok
}

func (D Data) HasRole(role string) bool {
	return slices.Contains(D.User.Roles, role)
}

func (D Data) T(id string, params ...string) string{
	td := make(map[string]interface{})

	if params != nil {
		if len(params) % 2 == 0 {
			for i := 1; i < len(params); i = i + 2 {
				td[params[i-1]] = params[i]
			}
		} else {
			panic("incorrect parameters")
		}
	}

	return D.Localiser.MustLocalize(
		&i18n.LocalizeConfig{
			MessageID   : id,
			TemplateData: td,
		},
	)
}

func (D Data) TFT() string {
	return time.RFC3339Nano
}
