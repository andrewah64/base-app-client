package i18n

import (
	"context"
)

import (
	"golang.org/x/text/language"
)

import (
	"github.com/andrewah64/base-app-client/internal/common/core/i18n"
)

func InitCache(ctx context.Context, dl language.Tag) (error) {
	if err := i18n.InitCache(ctx, dl, "i18n/api"); err != nil {
		return err
	}

	return nil
}
