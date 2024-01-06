package plugins

import (
	"fmt"
	"github.com/dpwgc/easierweb"
	"net/http"
)

func JSONErrorHandle(returnError bool) easierweb.ErrorHandle {
	return func(ctx *easierweb.Context, err any) {
		ctx.Error(fmt.Sprintf("%s -> [%s]%s | unexpected error: %s", ctx.Request.RemoteAddr, ctx.Request.Method, ctx.Request.RequestURI, err))
		res := make(map[string]string, 1)
		if returnError {
			res["error"] = fmt.Sprintf("%s", err)
		} else {
			res["error"] = "unexpected error"
		}
		ctx.WriteJSON(http.StatusInternalServerError, res)
	}
}

func YAMLErrorHandle(returnError bool) easierweb.ErrorHandle {
	return func(ctx *easierweb.Context, err any) {
		ctx.Error(fmt.Sprintf("%s -> [%s]%s | unexpected error: %s", ctx.Request.RemoteAddr, ctx.Request.Method, ctx.Request.RequestURI, err))
		res := make(map[string]string, 1)
		if returnError {
			res["error"] = fmt.Sprintf("%s", err)
		} else {
			res["error"] = "unexpected error"
		}
		ctx.WriteYAML(http.StatusInternalServerError, res)
	}
}

func XMLErrorHandle(returnError bool) easierweb.ErrorHandle {
	return func(ctx *easierweb.Context, err any) {
		ctx.Error(fmt.Sprintf("%s -> [%s]%s | unexpected error: %s", ctx.Request.RemoteAddr, ctx.Request.Method, ctx.Request.RequestURI, err))
		res := make(map[string]string, 1)
		if returnError {
			res["error"] = fmt.Sprintf("%s", err)
		} else {
			res["error"] = "unexpected error"
		}
		ctx.WriteXML(http.StatusInternalServerError, res)
	}
}

func StringErrorHandle(returnError bool) easierweb.ErrorHandle {
	return func(ctx *easierweb.Context, err any) {
		ctx.Error(fmt.Sprintf("%s -> [%s]%s | unexpected error: %s", ctx.Request.RemoteAddr, ctx.Request.Method, ctx.Request.RequestURI, err))
		res := "unexpected error"
		if returnError {
			res = fmt.Sprintf("%s", err)
		}
		ctx.WriteString(http.StatusInternalServerError, res)
	}
}
