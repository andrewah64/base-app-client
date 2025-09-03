package json

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

type Envelope map[string]any

func Read(ctx *context.Context, logger *slog.Logger, rw http.ResponseWriter, r *http.Request, dst any) error {
	const (
		maxBytes = 1_048_576
	)

	r.Body  = http.MaxBytesReader(rw, r.Body, int64(maxBytes))
	dec    := json.NewDecoder(r.Body)

	dec.DisallowUnknownFields()

	decErr := dec.Decode(dst)

	if decErr != nil {
		var (
			syntaxErr           *json.SyntaxError
			unmarshalTypeErr    *json.UnmarshalTypeError
			invalidUnmarshalErr *json.InvalidUnmarshalError
			maxBytesErr         *http.MaxBytesError
		)

		switch {
			case errors.As(decErr, &syntaxErr):
				logger.LogAttrs(*ctx, slog.LevelDebug, "syntax error",
					slog.String("error", decErr.Error()),
				)
				return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxErr.Offset)

			case errors.Is(decErr, io.ErrUnexpectedEOF):
				return errors.New("body contains badly-formed JSON")

			case errors.As(decErr, &unmarshalTypeErr):
				if unmarshalTypeErr.Field != "" {
					return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeErr.Field)
				}
				return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeErr.Offset)

			case errors.Is(decErr, io.EOF):
				return errors.New("body must not be empty")

			case errors.As(decErr, &maxBytesErr):
				return fmt.Errorf("body must not be larger than %d bytes", maxBytesErr.Limit)

			case errors.As(decErr, &invalidUnmarshalErr):
				slog.LogAttrs(*ctx, slog.LevelError, "invalidMarshalErr",
					slog.String("error", decErr.Error()),
				)
				panic(decErr)

			default:
				return decErr
		}
	}

	decErr = dec.Decode(&struct{}{})
	if !errors.Is(decErr, io.EOF) {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}

func Write (ctx *context.Context, logger *slog.Logger, rw http.ResponseWriter, status int, data Envelope, headers http.Header) error {
	js, jsErr := json.Marshal(data)
	if jsErr != nil {
		slog.LogAttrs(*ctx, slog.LevelError, "marshal json",
			slog.String("error", jsErr.Error()),
		)

		return jsErr
	}

	js = append(js, '\n')

	for k, v := range headers {
		rw.Header()[k] = v
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(status)
	rw.Write(js)

	return nil
}
