package plugins

import (
	"fmt"
	"github.com/dpwgc/easierweb"
	"log/slog"
	"net/http"
)

func JSONErrorHandle(ctx *easierweb.Context, err any) {
	slog.Error(fmt.Sprintf("%s -> [%s]%s | unexpected error: %s", ctx.Request.RemoteAddr, ctx.Request.Method, ctx.Request.RequestURI, err))
	res := make(map[string]string, 1)
	res["error"] = fmt.Sprintf("%s", err)
	ctx.WriteJSON(http.StatusInternalServerError, res)
}

func JSONNoDetailErrorHandle(ctx *easierweb.Context, err any) {
	slog.Error(fmt.Sprintf("%s -> [%s]%s | unexpected error: %s", ctx.Request.RemoteAddr, ctx.Request.Method, ctx.Request.RequestURI, err))
	res := make(map[string]string, 1)
	res["error"] = "unexpected error"
	ctx.WriteJSON(http.StatusInternalServerError, res)
}

func XMLErrorHandle(ctx *easierweb.Context, err any) {
	slog.Error(fmt.Sprintf("%s -> [%s]%s | unexpected error: %s", ctx.Request.RemoteAddr, ctx.Request.Method, ctx.Request.RequestURI, err))
	res := make(map[string]string, 1)
	res["error"] = fmt.Sprintf("%s", err)
	ctx.WriteXML(http.StatusInternalServerError, res)
}

func XMLNoDetailErrorHandle(ctx *easierweb.Context, err any) {
	slog.Error(fmt.Sprintf("%s -> [%s]%s | unexpected error: %s", ctx.Request.RemoteAddr, ctx.Request.Method, ctx.Request.RequestURI, err))
	res := make(map[string]string, 1)
	res["error"] = "unexpected error"
	ctx.WriteXML(http.StatusInternalServerError, res)
}

func YAMLErrorHandle(ctx *easierweb.Context, err any) {
	slog.Error(fmt.Sprintf("%s -> [%s]%s | unexpected error: %s", ctx.Request.RemoteAddr, ctx.Request.Method, ctx.Request.RequestURI, err))
	res := make(map[string]string, 1)
	res["error"] = fmt.Sprintf("%s", err)
	ctx.WriteYAML(http.StatusInternalServerError, res)
}

func YAMLNoDetailErrorHandle(ctx *easierweb.Context, err any) {
	slog.Error(fmt.Sprintf("%s -> [%s]%s | unexpected error: %s", ctx.Request.RemoteAddr, ctx.Request.Method, ctx.Request.RequestURI, err))
	res := make(map[string]string, 1)
	res["error"] = "unexpected error"
	ctx.WriteYAML(http.StatusInternalServerError, res)
}

func StringErrorHandle(ctx *easierweb.Context, err any) {
	slog.Error(fmt.Sprintf("%s -> [%s]%s | unexpected error: %s", ctx.Request.RemoteAddr, ctx.Request.Method, ctx.Request.RequestURI, err))
	ctx.WriteString(http.StatusInternalServerError, fmt.Sprintf("%s", err))
}

func StringNoDetailErrorHandle(ctx *easierweb.Context, err any) {
	slog.Error(fmt.Sprintf("%s -> [%s]%s | unexpected error: %s", ctx.Request.RemoteAddr, ctx.Request.Method, ctx.Request.RequestURI, err))
	ctx.WriteString(http.StatusInternalServerError, "unexpected error")
}
