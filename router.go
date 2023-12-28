package easierweb

import (
	"context"
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
	server                 http.Server
	middlewares            []Handle
	errorHandle            ErrorHandle
}

type Handle func(ctx *Context)

type SimpleHandle func(ctx *Context) (any, error)

type ResponseHandle func(ctx *Context, result any, err error)

type ErrorHandle func(ctx *Context, err any)

func New() *Router {
	return &Router{
		contextPath:            "",
		multipartFormMaxMemory: 1024,
		router:                 httprouter.New(),
		errorHandle: func(ctx *Context, err any) {
			fmt.Println("[ERROR]", ctx.Request.RemoteAddr, "->", ctx.Request.Method, ctx.Request.RequestURI, "::", err)
			// 返回code=500加异常信息
			ctx.WriteString(http.StatusInternalServerError, fmt.Sprintf("{\"error\":\"%s\"}", err))
		},
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

// --

func (r *Router) SimpleGET(path string, simpleHandle SimpleHandle, responseHandle ResponseHandle) *Router {
	return r.GET(path, r.simpleHandle2handle(simpleHandle, responseHandle))
}

func (r *Router) SimpleHEAD(path string, simpleHandle SimpleHandle, responseHandle ResponseHandle) *Router {
	return r.HEAD(path, r.simpleHandle2handle(simpleHandle, responseHandle))
}

func (r *Router) SimpleOPTIONS(path string, simpleHandle SimpleHandle, responseHandle ResponseHandle) *Router {
	return r.OPTIONS(path, r.simpleHandle2handle(simpleHandle, responseHandle))
}

func (r *Router) SimplePOST(path string, simpleHandle SimpleHandle, responseHandle ResponseHandle) *Router {
	return r.POST(path, r.simpleHandle2handle(simpleHandle, responseHandle))
}

func (r *Router) SimplePUT(path string, simpleHandle SimpleHandle, responseHandle ResponseHandle) *Router {
	return r.PUT(path, r.simpleHandle2handle(simpleHandle, responseHandle))
}

func (r *Router) SimplePATCH(path string, simpleHandle SimpleHandle, responseHandle ResponseHandle) *Router {
	return r.PATCH(path, r.simpleHandle2handle(simpleHandle, responseHandle))
}

func (r *Router) SimpleDELETE(path string, simpleHandle SimpleHandle, responseHandle ResponseHandle) *Router {
	return r.DELETE(path, r.simpleHandle2handle(simpleHandle, responseHandle))
}

// --

func (r *Router) GET(path string, handle Handle) *Router {
	r.router.GET(r.contextPath+path, func(res http.ResponseWriter, req *http.Request, par httprouter.Params) {
		r.handle(handle, res, req, par, nil)
	})
	return r
}

func (r *Router) HEAD(path string, handle Handle) *Router {
	r.router.HEAD(r.contextPath+path, func(res http.ResponseWriter, req *http.Request, par httprouter.Params) {
		r.handle(handle, res, req, par, nil)
	})
	return r
}

func (r *Router) OPTIONS(path string, handle Handle) *Router {
	r.router.OPTIONS(r.contextPath+path, func(res http.ResponseWriter, req *http.Request, par httprouter.Params) {
		r.handle(handle, res, req, par, nil)
	})
	return r
}

func (r *Router) POST(path string, handle Handle) *Router {
	r.router.POST(r.contextPath+path, func(res http.ResponseWriter, req *http.Request, par httprouter.Params) {
		r.handle(handle, res, req, par, nil)
	})
	return r
}

func (r *Router) PUT(path string, handle Handle) *Router {
	r.router.PUT(r.contextPath+path, func(res http.ResponseWriter, req *http.Request, par httprouter.Params) {
		r.handle(handle, res, req, par, nil)
	})
	return r
}

func (r *Router) PATCH(path string, handle Handle) *Router {
	r.router.PATCH(r.contextPath+path, func(res http.ResponseWriter, req *http.Request, par httprouter.Params) {
		r.handle(handle, res, req, par, nil)
	})
	return r
}

func (r *Router) DELETE(path string, handle Handle) *Router {
	r.router.DELETE(r.contextPath+path, func(res http.ResponseWriter, req *http.Request, par httprouter.Params) {
		r.handle(handle, res, req, par, nil)
	})
	return r
}

func (r *Router) WS(path string, handle Handle) *Router {
	r.router.GET(r.contextPath+path, func(res http.ResponseWriter, req *http.Request, par httprouter.Params) {
		websocket.Server{
			Handler: func(ws *websocket.Conn) {
				r.handle(handle, res, req, par, ws)
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
	return r.StaticFS(path, http.Dir(dir))
}

func (r *Router) StaticFS(path string, fs http.FileSystem) *Router {
	r.router.ServeFiles(r.contextPath+path, fs)
	return r
}

func (r *Router) Use(middlewares ...Handle) *Router {
	r.middlewares = append(r.middlewares, middlewares...)
	return r
}

func (r *Router) CustomHandle(method, path string, handle httprouter.Handle) *Router {
	r.router.Handle(method, r.contextPath+path, handle)
	return r
}

func (r *Router) Run(addr string) error {
	r.consoleStartPrint(addr)
	r.server = http.Server{
		Addr:    addr,
		Handler: r.router,
	}
	return r.server.ListenAndServe()
}

func (r *Router) RunTLS(addr string, certFile string, keyFile string, tlsConfig *tls.Config) error {
	r.consoleStartPrint(addr)
	r.server = http.Server{
		Addr:      addr,
		Handler:   r.router,
		TLSConfig: tlsConfig,
	}
	return r.server.ListenAndServeTLS(certFile, keyFile)
}

func (r *Router) Close() error {
	return r.server.Shutdown(context.Background())
}

func (r *Router) newContext(res http.ResponseWriter, req *http.Request, par httprouter.Params, ws *websocket.Conn) (*Context, error) {

	ctx := Context{
		index:          0,
		handles:        append([]Handle(nil), r.middlewares...),
		Header:         map[string]string{},
		Path:           map[string]string{},
		Query:          map[string]string{},
		Form:           map[string]string{},
		CustomCache:    make(map[string]any),
		Request:        req,
		ResponseWriter: res,
		WebsocketConn:  ws,
		Code:           -1,
		Result:         nil,
	}

	if strings.Contains(strings.ToLower(req.Header.Get("Content-Type")), "multipart/form-data") ||
		strings.Contains(strings.ToLower(req.Header.Get("content-type")), "multipart/form-data") {
		err := req.ParseMultipartForm(r.multipartFormMaxMemory)
		if err != nil {
			return nil, err
		}
	} else if strings.Contains(strings.ToLower(req.Header.Get("Content-Type")), "application/x-www-form-urlencoded") ||
		strings.Contains(strings.ToLower(req.Header.Get("content-type")), "application/x-www-form-urlencoded") {
		err := req.ParseForm()
		if err != nil {
			return nil, err
		}
	} else {
		bodyBytes, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		ctx.Body = bodyBytes
	}

	if len(req.Header) > 0 {
		ctx.Header = make(map[string]string, len(req.Header))
		for k, v := range req.Header {
			if len(v) > 0 {
				ctx.Header[k] = v[0]
			}
		}
	}

	if len(par) > 0 {
		ctx.Path = make(map[string]string, len(par))
		for _, v := range par {
			ctx.Path[v.Key] = v.Value
		}
	}

	if len(req.URL.Query()) > 0 {
		ctx.Query = make(map[string]string, len(req.URL.Query()))
		for k, v := range req.URL.Query() {
			if len(v) > 0 {
				ctx.Query[k] = v[0]
			}
		}
	}

	if len(req.PostForm) > 0 {
		ctx.Form = make(map[string]string, len(req.PostForm))
		for k, v := range req.PostForm {
			if len(v) > 0 {
				ctx.Form[k] = v[0]
			}
		}
	}

	return &ctx, nil
}

func (r *Router) handle(handle Handle, res http.ResponseWriter, req *http.Request, par httprouter.Params, ws *websocket.Conn) {

	ctx, err := r.newContext(res, req, par, nil)

	defer func() {
		sErr := recover()
		if sErr != nil && r.errorHandle != nil {
			r.errorBottomUp(ctx, sErr)
		}
	}()

	if err != nil {
		panic(err)
	}

	// 中间件
	if len(r.middlewares) > 0 {
		ctx.handles = append(ctx.handles, handle)
		for ctx.index < len(ctx.handles) {
			ctx.handles[ctx.index](ctx)
			ctx.index++
		}
	} else {
		handle(ctx)
	}

	// 如果ws存在，自动关闭ws连接
	if ws != nil {
		err = ws.Close()
		if err != nil {
			panic(err)
		}
	}
}

func (r *Router) simpleHandle2handle(simpleHandle SimpleHandle, responseHandle ResponseHandle) Handle {
	return func(ctx *Context) {
		result, err := simpleHandle(ctx)
		responseHandle(ctx, result, err)
	}
}

func (r *Router) errorBottomUp(ctx *Context, err any) {
	defer func() {
		_ = recover()
	}()
	r.errorHandle(ctx, err)
}

func (r *Router) consoleStartPrint(addr string) {
	fmt.Println("  ______          _        __          __  _     \n |  ____|        (_)       \\ \\        / / | |    \n | |__   __ _ ___ _  ___ _ _\\ \\  /\\  / /__| |__  \n |  __| / _` / __| |/ _ \\ '__\\ \\/  \\/ / _ \\ '_ \\ \n | |___| (_| \\__ \\ |  __/ |   \\  /\\  /  __/ |_) |\n |______\\__,_|___/_|\\___|_|    \\/  \\/ \\___|_.__/")
	fmt.Printf("\033[1;32;40m%s\033[0m\n", fmt.Sprintf(" >>> http server started on [%s] ", addr))
}
