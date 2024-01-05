package middlewares

import (
	"encoding/json"
	"github.com/dpwgc/easierweb"
	"log/slog"
	"time"
)

func Logger() easierweb.Handle {
	return func(ctx *easierweb.Context) {
		start := time.Now().UnixMilli()
		ctx.Next()
		end := time.Now().UnixMilli()
		timeCost := end - start

		path := ""
		query := ""
		form := ""
		body := ""
		result := ""

		if len(ctx.Path) > 0 {
			marshal, err := json.Marshal(ctx.Path)
			if err != nil {
				path = string(marshal)
			}
		}
		if len(ctx.Query) > 0 {
			marshal, err := json.Marshal(ctx.Query)
			if err != nil {
				query = string(marshal)
			}
		}
		if len(ctx.Form) > 0 {
			marshal, err := json.Marshal(ctx.Form)
			if err != nil {
				form = string(marshal)
			}
		}
		sizeLimit := 1024 * 1024
		if len(ctx.Body) > 0 {
			if len(ctx.Body) > sizeLimit {
				body = "body is too large"
			} else {
				body = string(ctx.Body)
			}
		}
		if len(ctx.Result) > 0 {
			if len(ctx.Result) > sizeLimit {
				result = "result is too large"
			} else {
				result = string(ctx.Body)
			}
		}

		slog.Info("request", slog.String("method", ctx.Request.Method),
			slog.String("url", ctx.Request.URL.String()),
			slog.String("client", ctx.Request.RemoteAddr),
			slog.String("path", path),
			slog.String("query", query),
			slog.String("form", form),
			slog.String("body", body),
			slog.Int("code", ctx.Code),
			slog.String("result", result),
			slog.Int64("timeCost", timeCost))
	}
}
