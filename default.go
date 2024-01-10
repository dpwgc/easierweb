package easierweb

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
)

// default function

func defaultRequestHandle() RequestHandle {
	return func(ctx *Context, reqObj any) error {
		if len(ctx.Form) > 0 {
			err := ctx.BindForm(reqObj)
			if err != nil {
				return err
			}
		} else if len(ctx.Body) > 0 {
			err := ctx.BindJSON(reqObj)
			if err != nil {
				return err
			}
		}
		if len(ctx.Query) > 0 {
			err := ctx.BindQuery(reqObj)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func defaultResponseHandle() ResponseHandle {
	return func(ctx *Context, result any, err error) {
		if err != nil {
			if result != nil {
				ctx.WriteJSON(http.StatusBadRequest, result)
				return
			}
			panic(err)
		}
		if result == nil {
			ctx.NoContent(http.StatusNoContent)
			return
		}
		ctx.WriteJSON(http.StatusOK, result)
	}
}

func defaultErrorHandle() ErrorHandle {
	return func(ctx *Context, err any) {
		ctx.Logger.Error(fmt.Sprintf("%s\n%s", err, string(debug.Stack())), slog.String("method", ctx.Request.Method), slog.String("route", ctx.Route))
		ctx.WriteString(http.StatusInternalServerError, fmt.Sprintf("{\"error\":\"%s\"}", err))
	}
}
