package notification 

import (
	"context"
	"log/slog"
	"net/http"
)

import (
	"github.com/andrewah64/base-app-client/internal/web/core/ui/data/page"
	"github.com/andrewah64/base-app-client/internal/web/core/ui/html"
)

func Toast(ctx context.Context, logger *slog.Logger, rw http.ResponseWriter, r *http.Request, ntfType string, msg *map[string]string, data *page.Data){
	data.NotificationData = &map[string]any{"Type": ntfType, "Messages" : msg}

	html.Fragment(ctx, logger, rw, r, "core/all/ntf/fragment/ntf", http.StatusCreated, data)
}

func Vrl (ctx context.Context, logger *slog.Logger, rw http.ResponseWriter, r *http.Request, title string, eq1ErrMg string, gt1ErrMsg string, msgs *[]string, data *page.Data) {
	data.NotificationData = &map[string]any{"Title": title, "Singular": eq1ErrMg, "Plural": gt1ErrMsg, "Msgs": msgs}

	html.Fragment(ctx, logger, rw, r, "core/all/ntf/fragment/vrl", http.StatusCreated, data)
}
