package middlewares

import (
	"github.com/dpwgc/easierweb"
	"net/http"
)

func CORS() easierweb.Handle {
	return func(ctx *easierweb.Context) {
		if ctx.Request.Header.Get("Origin") != "" {
			ctx.SetHeader("Access-Control-Allow-Origin", "*")
			ctx.SetHeader("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, PATCH, DELETE, HEAD")
			ctx.SetHeader("Access-Control-Allow-Headers", "*")
			ctx.SetHeader("Access-Control-Expose-Headers", "*")
			ctx.SetHeader("Access-Control-Allow-Credentials", "true")
		}
		if ctx.Request.Method == "OPTIONS" {
			ctx.WriteString(http.StatusOK, "Options Request!")
		}
		ctx.Next()
	}
}
