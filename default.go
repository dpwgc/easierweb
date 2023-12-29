package easierweb

import (
	"fmt"
	"net/http"
	"reflect"
)

func defaultRequestHandle(ctx *Context, paramValues []reflect.Value) error {
	if len(ctx.Form) > 0 {
		err := ctx.BindForm(paramValues[1].Addr().Interface())
		if err != nil {
			return err
		}
	} else if len(ctx.Body) > 0 {
		err := ctx.BindJSON(paramValues[1].Addr().Interface())
		if err != nil {
			return err
		}
	}
	if len(ctx.Query) > 0 {
		err := ctx.BindQuery(paramValues[1].Addr().Interface())
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
	fmt.Println("[ERROR]", ctx.Request.RemoteAddr, "->", ctx.Request.Method, ctx.Request.RequestURI, "::", err)
	// 返回code=500加异常信息
	ctx.WriteString(http.StatusInternalServerError, fmt.Sprintf("{\"error\":\"%s\"}", err))
}
