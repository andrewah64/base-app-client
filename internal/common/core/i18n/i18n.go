package i18n

import (
	"context"
	"io/fs"
	"log/slog"
	"path/filepath"
)

import (
	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

import (
	"golang.org/x/text/language"
)

import (
	"github.com/andrewah64/base-app-client/ui"
)

var (
	cache *i18n.Bundle
)

func Localiser(ctx context.Context, logger *slog.Logger, l string) (*i18n.Localizer) {
	logger.LogAttrs(ctx, slog.LevelDebug, "get localiser",
		slog.String("language", l),
	)
	
	return i18n.NewLocalizer(cache, l)
}

func InitCache(ctx context.Context, dl language.Tag, root string) (error) {
	slog.LogAttrs(ctx, slog.LevelInfo, "start")

	translations := i18n.NewBundle(dl)
	translations.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	wdErr := fs.WalkDir(ui.Files, root, func(path string, d fs.DirEntry, err error) error{
		if ! d.IsDir() && filepath.Ext(path) == ".toml" {
			slog.LogAttrs(ctx, slog.LevelInfo, "load i18nFile",
				slog.String("i18nFile", path),
			)

			if _, i18nFileErr := translations.LoadMessageFileFS(ui.Files, path); i18nFileErr != nil {
				slog.LogAttrs(ctx, slog.LevelError, "load i18nFile",
					slog.String("i18nFile", path),
					slog.String("error"   , i18nFileErr.Error()),
				)

				return i18nFileErr
			}
		}
		return nil
	})

	if wdErr != nil {
		slog.LogAttrs(ctx, slog.LevelError, "load i18nFiles",
			slog.String("error" , wdErr.Error()),
		)

		return wdErr
	}

	cache = translations

	return nil
}
