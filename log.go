package main

import (
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

// getLogger returns new logger.
func getLogger(lvl string) zerolog.Logger {
	level, err := zerolog.ParseLevel(lvl)
	if err != nil {
		panic(err)
	}

	return zerolog.New(os.Stdout).
		Level(level).
		With().
		Timestamp().
		Logger()
}

func accessHandler(log zerolog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := log.With().Logger()
		l.UpdateContext(func(c zerolog.Context) zerolog.Context {
			return c.
				Str("user_agent", r.Header.Get("User-Agent")).
				Str("url", r.URL.Path).
				Str("method", r.Method).
				Str("request-id", uuid.New().String())
		})
		r = r.WithContext(l.WithContext(r.Context()))

		start := time.Now()
		lw := wrapWriter(w)
		next.ServeHTTP(lw, r)

		zerolog.Ctx(r.Context()).Info().
			Int("status", lw.Status()).
			Int("size", lw.BytesWritten()).
			Dur("duration", time.Since(start)).
			Msg(http.StatusText(lw.Status()))
	})
}
