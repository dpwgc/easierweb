package plugins

import (
	"fmt"
	"github.com/dpwgc/easierweb"
	"log/slog"
	"net/http"
	"runtime/debug"
)

type ErrorHandleOptions struct {
	ShowError   bool
	OutputStack bool
}

type Err struct {
	Msg string `json:"msg" xml:"Msg" yaml:"msg"`
}

func JSONErrorHandle(opts ...ErrorHandleOptions) easierweb.ErrorHandle {
	return func(ctx *easierweb.Context, err any) {
		logError(ctx, err, opts...)
		res := Err{}
		if len(opts) > 0 && opts[0].ShowError {
			res.Msg = fmt.Sprintf("%s", err)
		} else {
			res.Msg = "unexpected error"
		}
		ctx.WriteJSON(http.StatusInternalServerError, res)
	}
}

func YAMLErrorHandle(opts ...ErrorHandleOptions) easierweb.ErrorHandle {
	return func(ctx *easierweb.Context, err any) {
		logError(ctx, err, opts...)
		res := Err{}
		if len(opts) > 0 && opts[0].ShowError {
			res.Msg = fmt.Sprintf("%s", err)
		} else {
			res.Msg = "unexpected error"
		}
		ctx.WriteYAML(http.StatusInternalServerError, res)
	}
}

func XMLErrorHandle(opts ...ErrorHandleOptions) easierweb.ErrorHandle {
	return func(ctx *easierweb.Context, err any) {
		logError(ctx, err, opts...)
		res := Err{}
		if len(opts) > 0 && opts[0].ShowError {
			res.Msg = fmt.Sprintf("%s", err)
		} else {
			res.Msg = "unexpected error"
		}
		ctx.WriteXML(http.StatusInternalServerError, res)
	}
}

func StringErrorHandle(opts ...ErrorHandleOptions) easierweb.ErrorHandle {
	return func(ctx *easierweb.Context, err any) {
		logError(ctx, err, opts...)
		if len(opts) > 0 && opts[0].ShowError {
			ctx.WriteString(http.StatusInternalServerError, fmt.Sprintf("%s", err))
		} else {
			ctx.WriteString(http.StatusInternalServerError, "unexpected error")
		}
	}
}

func logError(ctx *easierweb.Context, err any, opts ...ErrorHandleOptions) {
	if len(opts) > 0 && opts[0].OutputStack {
		ctx.Logger.Error(fmt.Sprintf("%s\n%s", err, string(debug.Stack())), slog.String("method", ctx.Request.Method), slog.String("route", ctx.Route))
	} else {
		ctx.Logger.Error(fmt.Sprintf("%s", err), slog.String("method", ctx.Request.Method), slog.String("route", ctx.Route))
	}
}
