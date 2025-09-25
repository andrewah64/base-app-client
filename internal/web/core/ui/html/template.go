package html

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
	"path/filepath"
	"slices"
	"strings"
	"time"
)

import (
	e "github.com/andrewah64/base-app-client/internal/web/core/error"
	  "github.com/andrewah64/base-app-client/ui"
)

var (
	cache map[string]*template.Template = make(map[string]*template.Template)
)

func HiddenUtsFragment(rw http.ResponseWriter, cid string, id string, name string, uts time.Time, format string) {
	fragTxt := `<div hx-swap-oob="true" id="%v"><input id="%v" name="%v" value="%v" type="hidden"></div>`
	frag    := fmt.Sprintf(fragTxt, cid, id, name, uts.Format(format))

	rw.Write([]byte(frag))
}

func Fragment(ctx context.Context, logger *slog.Logger, rw http.ResponseWriter, r *http.Request, cacheKey string, status int, data any) {
	logger.LogAttrs(ctx, slog.LevelDebug, "get fragment",
		slog.String("cacheKey", cacheKey),
		slog.String("data"    , fmt.Sprintf("%v", data)),
	)

	if tmpl, ok := cache[cacheKey]; ok {
		err := template.Must(tmpl, nil).Execute(rw, data)
		if err != nil {
			e.IntSrv(ctx, rw, err)
			return
		}
	} else {
		panic(fmt.Sprintf("fragment '%v' not found", cacheKey))
	}
}

func Tmpl(ctx context.Context, logger *slog.Logger, rw http.ResponseWriter, r *http.Request, cacheKey string, status int, data any) {
	cacheKeySegments := strings.Split(cacheKey, "/")
	tmplType         := strings.Split(cacheKey, "/")[1]

	if f := slices.Index(cacheKeySegments, "template"); f != -1 {
		tmplType = cacheKeySegments[f + 1]
	}

	logger.LogAttrs(ctx, slog.LevelDebug, "get page",
		slog.String("cacheKey", cacheKey),
		slog.String("tmplType", tmplType),
	)

	if tmpl, ok := cache[cacheKey]; ok {
		buf := new(bytes.Buffer)

		err := template.Must(tmpl, nil).ExecuteTemplate(buf, tmplType, data)
		if err != nil {
			e.IntSrv(ctx, rw, err)
			return
		}

		rw.WriteHeader(status)

		buf.WriteTo(rw)
	} else {
		panic(fmt.Sprintf("template '%v' not found", cacheKey))
	}
}

func InitCache(ctx context.Context){
	tmpls := make(map[string][]string)

	const (
		tmplRoot     = "html/category"
		tmplTypePath = "html/template/*"
	)

	tmplTypes, err := fs.Glob(ui.Files, tmplTypePath)
	if err != nil {
		panic(fmt.Sprintf("%v", err.Error))
	}

	for _, tmplType := range tmplTypes {
		files, err := fs.Glob(ui.Files, fmt.Sprintf("%v/*.html", tmplType))
		if err != nil {
			panic(fmt.Sprintf("%v", err.Error))
		}

		tmpls[filepath.Base(tmplType)] = files
	}

	wdErr := fs.WalkDir(ui.Files, tmplRoot, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			switch d.Name() {
				case "content":
					pageKey, _ := filepath.Rel(tmplRoot, path)

					tmplType   := (strings.Split(pageKey, "/"))[1]

					files := append([]string{fmt.Sprintf("%v/index.html", path)}, tmpls[tmplType]...)

					tmpls, tmplsErr := fs.Glob(ui.Files, fmt.Sprintf("%v/template/*.html", filepath.Dir(path)))
					if tmplsErr != nil {
						panic(fmt.Sprintf("%v", err.Error))
					}

					if tmpls != nil {
						files = append(tmpls, files...)
					}

					tmpl, tmplErr := template.ParseFS(ui.Files, files...)
					if tmplErr != nil {
						panic(fmt.Sprintf("%v", tmplErr.Error()))
					}

					cache[pageKey] = tmpl
				case "fragment", "template":
					pageKey, _ := filepath.Rel(tmplRoot, path)

					tmpls, tmplsErr := fs.Glob(ui.Files, fmt.Sprintf("%v/*.html", path))
					if tmplsErr != nil {
						panic(fmt.Sprintf("%v", err.Error))
					}

					if tmpls != nil {
						for _, t := range tmpls {
							tmpl, tmplErr := template.ParseFS(ui.Files, t)
							if tmplErr != nil {
								panic(fmt.Sprintf("%v", tmplErr.Error()))
							}

							tmplKey := fmt.Sprintf("%v/%v", pageKey, strings.TrimSuffix(filepath.Base(t), ".html"))

							cache[tmplKey] = tmpl
						}
					}
			}
		}

		return nil
	})

	if wdErr != nil {
		slog.LogAttrs(ctx, slog.LevelError, "load templates",
			slog.String("error" , wdErr.Error()),
		)
		panic(fmt.Sprintf("failed to load templates %v", wdErr.Error()))
	}

	if len(cache) == 0 {
		panic("The template cache is empty")
	}
}
