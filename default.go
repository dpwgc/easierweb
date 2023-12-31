package easierweb

import (
	"fmt"
	"log/slog"
	"net/http"
)

func defaultRequestHandle(ctx *Context, reqObj any) error {
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

func defaultResponseHandle(ctx *Context, result any, err error) {
	if err != nil {
		errRes := make(map[string]string, 1)
		errRes["error"] = err.Error()
		ctx.WriteJSON(http.StatusBadRequest, errRes)
		return
	}
	if result == nil {
		ctx.Write(http.StatusNoContent, nil)
		return
	}
	ctx.WriteJSON(http.StatusOK, result)
}

func defaultErrorHandle(ctx *Context, err any) {
	slog.Error(fmt.Sprintf("%s -> [%s]%s | unexpected error: %s", ctx.Request.RemoteAddr, ctx.Request.Method, ctx.Request.RequestURI, err))
	ctx.WriteString(http.StatusInternalServerError, fmt.Sprintf("{\"error\":\"%s\"}", err))
}
