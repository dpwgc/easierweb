package easierweb

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/websocket"
	"log/slog"
	"net/http"
	"sync"
)

type RouterOptions struct {
	RootPath               string
	MultipartFormMaxMemory int64
	ErrorHandle            ErrorHandle
	RequestHandle          RequestHandle
	ResponseHandle         ResponseHandle
	Logger                 *slog.Logger
	CloseConsolePrint      bool
}

type Router struct {
	rootPath               string
	multipartFormMaxMemory int64
	router                 *httprouter.Router
	server                 *http.Server
	middlewares            []Handle
	errorHandle            ErrorHandle
	requestHandle          RequestHandle
	responseHandle         ResponseHandle
	logger                 *slog.Logger
	contextPool            *sync.Pool
	closeConsolePrint      bool
}

func New(opts ...RouterOptions) *Router {
	r := &Router{
		multipartFormMaxMemory: 32 << 20,
		router:                 httprouter.New(),
		errorHandle:            defaultErrorHandle(),
		requestHandle:          defaultRequestHandle(),
		responseHandle:         defaultResponseHandle(),
		logger:                 slog.Default(),
		contextPool: &sync.Pool{
			New: func() any {
				return new(Context)
			},
		},
	}
	for _, v := range opts {
		if v.RootPath != "" {
			r.rootPath = v.RootPath
		}
		if v.MultipartFormMaxMemory > 0 {
			r.multipartFormMaxMemory = v.MultipartFormMaxMemory
		}
		if v.ErrorHandle != nil {
			r.errorHandle = v.ErrorHandle
		}
		if v.RequestHandle != nil {
			r.requestHandle = v.RequestHandle
		}
		if v.ResponseHandle != nil {
			r.responseHandle = v.ResponseHandle
		}
		if v.Logger != nil {
			r.logger = v.Logger
		}
		r.closeConsolePrint = v.CloseConsolePrint
	}
	return r
}

// easier usage function

func (r *Router) EasyGET(path string, easyHandle any, middlewares ...Handle) *Router {
	return r.EasyHandle(MethodGET, path, easyHandle, middlewares...)
}

func (r *Router) EasyHEAD(path string, easyHandle any, middlewares ...Handle) *Router {
	return r.EasyHandle(MethodHEAD, path, easyHandle, middlewares...)
}

func (r *Router) EasyOPTIONS(path string, easyHandle any, middlewares ...Handle) *Router {
	return r.EasyHandle(MethodOPTIONS, path, easyHandle, middlewares...)
}

func (r *Router) EasyPOST(path string, easyHandle any, middlewares ...Handle) *Router {
	return r.EasyHandle(MethodPOST, path, easyHandle, middlewares...)
}

func (r *Router) EasyPUT(path string, easyHandle any, middlewares ...Handle) *Router {
	return r.EasyHandle(MethodPUT, path, easyHandle, middlewares...)
}

func (r *Router) EasyPATCH(path string, easyHandle any, middlewares ...Handle) *Router {
	return r.EasyHandle(MethodPATCH, path, easyHandle, middlewares...)
}

func (r *Router) EasyDELETE(path string, easyHandle any, middlewares ...Handle) *Router {
	return r.EasyHandle(MethodDELETE, path, easyHandle, middlewares...)
}

func (r *Router) EasyHandle(method, path string, easyHandle any, middlewares ...Handle) *Router {
	return r.Handle(method, path, r.easyHandle(easyHandle), middlewares...)
}

func (r *Router) EasyAny(path string, easyHandle any, middlewares ...Handle) *Router {
	return r.Any(path, r.easyHandle(easyHandle), middlewares...)
}

// basic usage function

func (r *Router) GET(path string, handle Handle, middlewares ...Handle) *Router {
	return r.Handle(MethodGET, path, handle, middlewares...)
}

func (r *Router) HEAD(path string, handle Handle, middlewares ...Handle) *Router {
	return r.Handle(MethodHEAD, path, handle, middlewares...)
}

func (r *Router) OPTIONS(path string, handle Handle, middlewares ...Handle) *Router {
	return r.Handle(MethodOPTIONS, path, handle, middlewares...)
}

func (r *Router) POST(path string, handle Handle, middlewares ...Handle) *Router {
	return r.Handle(MethodPOST, path, handle, middlewares...)
}

func (r *Router) PUT(path string, handle Handle, middlewares ...Handle) *Router {
	return r.Handle(MethodPUT, path, handle, middlewares...)
}

func (r *Router) PATCH(path string, handle Handle, middlewares ...Handle) *Router {
	return r.Handle(MethodPATCH, path, handle, middlewares...)
}

func (r *Router) DELETE(path string, handle Handle, middlewares ...Handle) *Router {
	return r.Handle(MethodDELETE, path, handle, middlewares...)
}

var methodNames = []string{MethodGET, MethodHEAD, MethodOPTIONS, MethodPOST, MethodPUT, MethodPATCH, MethodDELETE}

func (r *Router) Any(path string, handle Handle, middlewares ...Handle) *Router {
	for _, method := range methodNames {
		r.Handle(method, path, handle, middlewares...)
	}
	return r
}

func (r *Router) Handle(method, path string, handle Handle, middlewares ...Handle) *Router {
	route := r.rootPath + path
	r.router.Handle(method, route, func(res http.ResponseWriter, req *http.Request, par httprouter.Params) {
		r.handle(route, handle, res, req, par, nil, false, middlewares...)
	})
	return r
}

func (r *Router) WS(path string, handle Handle, middlewares ...Handle) *Router {
	route := r.rootPath + path
	r.router.GET(route, func(res http.ResponseWriter, req *http.Request, par httprouter.Params) {
		websocket.Server{
			Handler: func(ws *websocket.Conn) {
				r.handle(route, handle, res, req, par, ws, false, middlewares...)
			},
			Handshake: func(config *websocket.Config, req *http.Request) error {
				// 解决跨域
				return nil
			},
		}.ServeHTTP(res, req)
	})
	return r
}

func (r *Router) SSE(path string, handle Handle, middlewares ...Handle) *Router {
	route := r.rootPath + path
	r.router.GET(route, func(res http.ResponseWriter, req *http.Request, par httprouter.Params) {
		r.handle(route, handle, res, req, par, nil, true, middlewares...)
	})
	return r
}

func (r *Router) Static(path, dir string) *Router {
	return r.StaticFS(path, http.Dir(dir))
}

func (r *Router) StaticFS(path string, fs http.FileSystem) *Router {
	route := r.rootPath + path
	r.router.ServeFiles(route, fs)
	return r
}

func (r *Router) Use(middlewares ...Handle) *Router {
	r.middlewares = append(r.middlewares, middlewares...)
	return r
}

func (r *Router) Run(addr string) error {
	return r.Serve(&http.Server{
		Addr: addr,
	})
}

func (r *Router) RunTLS(addr string, certFile string, keyFile string, tlsConfig *tls.Config) error {
	return r.ServeTLS(&http.Server{
		Addr:      addr,
		TLSConfig: tlsConfig,
	}, certFile, keyFile)
}

func (r *Router) Serve(server *http.Server) error {
	r.server = server
	r.server.Handler = r.router
	r.consoleStartPrint(r.server.Addr)
	return r.server.ListenAndServe()
}

func (r *Router) ServeTLS(server *http.Server, certFile string, keyFile string) error {
	r.server = server
	r.server.Handler = r.router
	r.consoleStartPrint(r.server.Addr)
	return r.server.ListenAndServeTLS(certFile, keyFile)
}

func (r *Router) Close() error {
	return r.server.Shutdown(context.Background())
}

func (r *Router) consoleStartPrint(addr string) {
	if r.closeConsolePrint {
		return
	}
	fmt.Println("  ______          _        __          __  _     \n |  ____|        (_)       \\ \\        / / | |    \n | |__   __ _ ___ _  ___ _ _\\ \\  /\\  / /__| |__  \n |  __| / _` / __| |/ _ \\ '__\\ \\/  \\/ / _ \\ '_ \\ \n | |___| (_| \\__ \\ |  __/ |   \\  /\\  /  __/ |_) |\n |______\\__,_|___/_|\\___|_|    \\/  \\/ \\___|_.__/")
	fmt.Printf("\033[1;32;40m%s\033[0m\n", fmt.Sprintf(" >>> server runs on [%s] ", addr))
}
