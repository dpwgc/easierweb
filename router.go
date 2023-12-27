package easierweb

import (
	"crypto/tls"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/websocket"
	"io"
	"net/http"
)

type Router struct {
	contextPath string
	router      *httprouter.Router
	middlewares []Method
	errorHandle ErrorHandle
}

type ErrorHandle func(ctx *Context, err any)

func NewRouter(contextPath string) *Router {
	return &Router{
		contextPath: contextPath,
		router:      httprouter.New(),
	}
}

func (r *Router) GET(path string, method Method) {
	r.router.GET(r.contextPath+path, func(res http.ResponseWriter, req *http.Request, par httprouter.Params) {
		r.handle(method, res, req, par)
	})
}

func (r *Router) HEAD(path string, method Method) {
	r.router.HEAD(r.contextPath+path, func(res http.ResponseWriter, req *http.Request, par httprouter.Params) {
		r.handle(method, res, req, par)
	})
}

func (r *Router) OPTIONS(path string, method Method) {
	r.router.OPTIONS(r.contextPath+path, func(res http.ResponseWriter, req *http.Request, par httprouter.Params) {
		r.handle(method, res, req, par)
	})
}

func (r *Router) POST(path string, method Method) {
	r.router.POST(r.contextPath+path, func(res http.ResponseWriter, req *http.Request, par httprouter.Params) {
		r.handle(method, res, req, par)
	})
}

func (r *Router) PUT(path string, method Method) {
	r.router.PUT(r.contextPath+path, func(res http.ResponseWriter, req *http.Request, par httprouter.Params) {
		r.handle(method, res, req, par)
	})
}

func (r *Router) PATCH(path string, method Method) {
	r.router.PATCH(r.contextPath+path, func(res http.ResponseWriter, req *http.Request, par httprouter.Params) {
		r.handle(method, res, req, par)
	})
}

func (r *Router) DELETE(path string, method Method) {
	r.router.DELETE(r.contextPath+path, func(res http.ResponseWriter, req *http.Request, par httprouter.Params) {
		r.handle(method, res, req, par)
	})
}

func (r *Router) WS(path string) {
	r.router.GET(r.contextPath+path, func(res http.ResponseWriter, req *http.Request, par httprouter.Params) {
		websocket.Handler(func(ws *websocket.Conn) {
			// TODO
		}).ServeHTTP(res, req)
	})
}

func (r *Router) StaticFiles(path string, fs http.FileSystem) {
	r.router.ServeFiles(r.contextPath+path, fs)
}

// AddMiddleware
// append the middleware
func (r *Router) AddMiddleware(middleware Method) *Router {
	r.middlewares = append(r.middlewares, middleware)
	return r
}

// AddMiddlewares
// append the middleware
func (r *Router) AddMiddlewares(middlewares ...Method) *Router {
	r.middlewares = append(r.middlewares, middlewares...)
	return r
}

func (r *Router) SetErrorHandle(errorHandle ErrorHandle) {
	r.errorHandle = errorHandle
}

func (r *Router) Run(addr string) error {
	fmt.Println("Run ->", addr)
	return http.ListenAndServe(addr, r.router)
}

func (r *Router) RunTLS(addr string, certFile string, keyFile string, tlsConfig *tls.Config) error {
	fmt.Println("RunTLS ->", addr)
	server := http.Server{
		Addr:      addr,
		Handler:   r.router,
		TLSConfig: tlsConfig,
	}
	return server.ListenAndServeTLS(certFile, keyFile)
}

func (r *Router) handle(method Method, res http.ResponseWriter, req *http.Request, par httprouter.Params) {

	ctx := Context{
		index:          0,
		methods:        append([]Method(nil), r.middlewares...),
		Header:         make(map[string]string),
		Path:           make(map[string]string),
		Query:          make(map[string]string),
		Form:           make(map[string]string),
		Request:        req,
		ResponseWriter: res,
		Code:           -1,
		Result:         nil,
	}

	defer func() {
		err := recover()
		if err != nil {
			r.errorHandle(&ctx, err)
		}
	}()

	bodyBytes, err := io.ReadAll(req.Body)
	if err == nil && len(bodyBytes) > 0 {
		ctx.Body = bodyBytes
	}
	_ = req.Body.Close()

	for k, v := range req.Header {
		if len(v) > 0 {
			ctx.Header[k] = v[0]
		}
	}

	for _, v := range par {
		ctx.Path[v.Key] = v.Value
	}

	for k, v := range req.URL.Query() {
		if len(v) > 0 {
			ctx.Query[k] = v[0]
		}
	}

	for k, v := range req.Form {
		if len(v) > 0 {
			ctx.Form[k] = v[0]
		}
	}

	if len(r.middlewares) > 0 {
		ctx.methods = append(ctx.methods, method)
		for ctx.index < len(ctx.methods) {
			ctx.methods[ctx.index](&ctx)
			ctx.index++
		}
	} else {
		method(&ctx)
	}
}
