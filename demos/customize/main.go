package main

import (
	"fmt"
	"github.com/dpwgc/easierweb"
	"log/slog"
	"net/http"
)

func main() {
	easierweb.New(easierweb.RouterOptions{
		// customize unexpected error handling logic
		ErrorHandle: customErrorHandle(),
		// customize the auto-binding request data logic
		RequestHandle: customRequestHandle(),
		// customize the auto-writing response data logic
		ResponseHandle: customResponseHandle(),
		// form request body size limit
		MultipartFormMaxMemory: 4096,
		// whether to turn off console output
		CloseConsolePrint: false,
	})
}

func customErrorHandle() easierweb.ErrorHandle {
	return func(ctx *easierweb.Context, err any) {
		slog.Error("unexpected error", slog.String("error", fmt.Sprintf("%s", err)))
		ctx.WriteString(http.StatusBadRequest, "unexpected error!!!")
	}
}

func customRequestHandle() easierweb.RequestHandle {
	return func(ctx *easierweb.Context, reqObj any) error {
		if len(ctx.Body) > 0 {
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

func customResponseHandle() easierweb.ResponseHandle {
	return func(ctx *easierweb.Context, result any, err error) {
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
