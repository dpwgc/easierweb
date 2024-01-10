package main

import (
	"crypto/tls"
	"fmt"
	"github.com/dpwgc/easierweb"
	"github.com/dpwgc/easierweb/plugins"
	"log/slog"
	"net/http"
)

// you can customize error, request, and response handle functions
// or use the plugin provided by the framework
func main() {

	// customize
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

	// use framework plugins
	router := easierweb.New(easierweb.RouterOptions{
		// respond to error messages in xml format
		ErrorHandle: plugins.XMLErrorHandle(plugins.ErrorHandleOptions{
			// render an error in the response body
			ShowError: true,
			// output stack info in logs
			OutputStack: true,
		}),
		// request data is parsed and automatically bound using xml format
		RequestHandle: plugins.XMLRequestHandle(),
		// write response data in xml format
		ResponseHandle: plugins.XMLResponseHandle(),
	})

	// configure TLS and start router
	err := router.RunTLS("127.0.0.1:8080", "cert.pem", "private.key", &tls.Config{})
	if err != nil {
		panic(err)
	}
}

func customErrorHandle() easierweb.ErrorHandle {
	return func(ctx *easierweb.Context, err any) {
		ctx.Logger.Error("unexpected error", slog.String("error", fmt.Sprintf("%s", err)))
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
