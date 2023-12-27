package easierweb

import (
	"crypto/tls"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/websocket"
	"io"
	"net/http"
	"strings"
)

type Router struct {
	contextPath            string
	multipartFormMaxMemory int64
	router                 *httprouter.Router
	middlewares            []Method
	errorHandle            ErrorHandle
}

type ErrorHandle func(ctx *Context, err any)

func New() *Router {
	return &Router{
		contextPath:            "",
		multipartFormMaxMemory: 1024,
		router:                 httprouter.New(),
	}
}

func (r *Router) SetErrorHandle(errorHandle ErrorHandle) *Router {
	r.errorHandle = errorHandle
	return r
}

func (r *Router) SetContextPath(contextPath string) *Router {
	r.contextPath = contextPath
	return r
}

func (r *Router) SetMultipartFormMaxMemory(maxMemory int64) *Router {
	r.multipartFormMaxMemory = maxMemory
	return r
}

func (r *Router) GET(path string, method Method) *Router {
	r.router.GET(r.contextPath+path, func(res http.ResponseWriter, req *http.Request, par httprouter.Params) {
		r.handle(method, res, req, par, nil)
	})
	return r
}

func (r *Router) HEAD(path string, method Method) *Router {
	r.router.HEAD(r.contextPath+path, func(res http.ResponseWriter, req *http.Request, par httprouter.Params) {
		r.handle(method, res, req, par, nil)
	})
	return r
}

func (r *Router) OPTIONS(path string, method Method) *Router {
	r.router.OPTIONS(r.contextPath+path, func(res http.ResponseWriter, req *http.Request, par httprouter.Params) {
		r.handle(method, res, req, par, nil)
	})
	return r
}

func (r *Router) POST(path string, method Method) *Router {
	r.router.POST(r.contextPath+path, func(res http.ResponseWriter, req *http.Request, par httprouter.Params) {
		r.handle(method, res, req, par, nil)
	})
	return r
}

func (r *Router) PUT(path string, method Method) *Router {
	r.router.PUT(r.contextPath+path, func(res http.ResponseWriter, req *http.Request, par httprouter.Params) {
		r.handle(method, res, req, par, nil)
	})
	return r
}

func (r *Router) PATCH(path string, method Method) *Router {
	r.router.PATCH(r.contextPath+path, func(res http.ResponseWriter, req *http.Request, par httprouter.Params) {
		r.handle(method, res, req, par, nil)
	})
	return r
}

func (r *Router) DELETE(path string, method Method) *Router {
	r.router.DELETE(r.contextPath+path, func(res http.ResponseWriter, req *http.Request, par httprouter.Params) {
		r.handle(method, res, req, par, nil)
	})
	return r
}

func (r *Router) Handle(method, path string, handle httprouter.Handle) *Router {
	r.router.Handle(method, r.contextPath+path, handle)
	return r
}

func (r *Router) WS(path string, method Method) *Router {
	r.router.GET(r.contextPath+path, func(res http.ResponseWriter, req *http.Request, par httprouter.Params) {
		websocket.Server{
			Handler: func(ws *websocket.Conn) {
				r.handle(method, res, req, par, ws)
			},
			Handshake: func(config *websocket.Config, req *http.Request) error {
				// 解决跨域
				return nil
			},
		}.ServeHTTP(res, req)
	})
	return r
}

func (r *Router) Static(path, dir string) *Router {
	r.StaticFS(path, http.Dir(dir))
	return r
}

func (r *Router) StaticFS(path string, fs http.FileSystem) *Router {
	r.router.ServeFiles(r.contextPath+path, fs)
	return r
}

func (r *Router) AddMiddleware(middleware Method) *Router {
	r.middlewares = append(r.middlewares, middleware)
	return r
}

func (r *Router) AddMiddlewares(middlewares ...Method) *Router {
	r.middlewares = append(r.middlewares, middlewares...)
	return r
}

func (r *Router) Run(addr string) error {
	r.consoleStartPrint(addr)
	return http.ListenAndServe(addr, r.router)
}

func (r *Router) RunTLS(addr string, certFile string, keyFile string, tlsConfig *tls.Config) error {
	r.consoleStartPrint(addr)
	server := http.Server{
		Addr:      addr,
		Handler:   r.router,
		TLSConfig: tlsConfig,
	}
	return server.ListenAndServeTLS(certFile, keyFile)
}

func (r *Router) handle(method Method, res http.ResponseWriter, req *http.Request, par httprouter.Params, ws *websocket.Conn) {

	ctx := Context{
		index:          0,
		methods:        append([]Method(nil), r.middlewares...),
		Header:         make(map[string]string),
		Path:           make(map[string]string),
		Query:          make(map[string]string),
		Form:           make(map[string]string),
		Request:        req,
		ResponseWriter: res,
		WebsocketConn:  ws,
		Code:           -1,
		Result:         nil,
	}

	defer func() {
		err := recover()
		if err != nil && r.errorHandle != nil {
			r.errorHandle(&ctx, err)
		}
	}()

	if !strings.Contains(strings.ToLower(req.Header.Get("Content-Type")), "multipart/form-data") &&
		!strings.Contains(strings.ToLower(req.Header.Get("content-type")), "multipart/form-data") {
		bodyBytes, err := io.ReadAll(req.Body)
		if err == nil {
			ctx.Body = bodyBytes
		}
	} else {
		err := req.ParseMultipartForm(r.multipartFormMaxMemory)
		if err != nil {
			panic(err)
		}
	}

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

	for k, v := range req.PostForm {
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
	// 如果websocket连接不为空，关闭它
	if ws != nil {
		err := ws.Close()
		if err != nil {
			panic(err)
		}
	}
}

func (r *Router) consoleStartPrint(addr string) {
	fmt.Println("  ______          _        __          __  _     \n |  ____|        (_)       \\ \\        / / | |    \n | |__   __ _ ___ _  ___ _ _\\ \\  /\\  / /__| |__  \n |  __| / _` / __| |/ _ \\ '__\\ \\/  \\/ / _ \\ '_ \\ \n | |___| (_| \\__ \\ |  __/ |   \\  /\\  /  __/ |_) |\n |______\\__,_|___/_|\\___|_|    \\/  \\/ \\___|_.__/")
	fmt.Printf("\033[1;32;40m%s\033[0m\n", fmt.Sprintf(" >>> http server started on [%s] ", addr))
}
