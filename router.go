package easierweb

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/websocket"
	"log/slog"
	"net/http"
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
	server                 http.Server
	middlewares            []Handle
	errorHandle            ErrorHandle
	requestHandle          RequestHandle
	responseHandle         ResponseHandle
	logger                 *slog.Logger
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
	return r.GET(path, r.easyHandle(easyHandle), middlewares...)
}

func (r *Router) EasyHEAD(path string, easyHandle any, middlewares ...Handle) *Router {
	return r.HEAD(path, r.easyHandle(easyHandle), middlewares...)
}

func (r *Router) EasyOPTIONS(path string, easyHandle any, middlewares ...Handle) *Router {
	return r.OPTIONS(path, r.easyHandle(easyHandle), middlewares...)
}

func (r *Router) EasyPOST(path string, easyHandle any, middlewares ...Handle) *Router {
	return r.POST(path, r.easyHandle(easyHandle), middlewares...)
}

func (r *Router) EasyPUT(path string, easyHandle any, middlewares ...Handle) *Router {
	return r.PUT(path, r.easyHandle(easyHandle), middlewares...)
}

func (r *Router) EasyPATCH(path string, easyHandle any, middlewares ...Handle) *Router {
	return r.PATCH(path, r.easyHandle(easyHandle), middlewares...)
}

func (r *Router) EasyDELETE(path string, easyHandle any, middlewares ...Handle) *Router {
	return r.DELETE(path, r.easyHandle(easyHandle), middlewares...)
}

func (r *Router) EasyAny(path string, easyHandle any, middlewares ...Handle) *Router {
	return r.Any(path, r.easyHandle(easyHandle), middlewares...)
}

func (r *Router) EasyHandle(method, path string, easyHandle any, middlewares ...Handle) *Router {
	return r.Handle(method, path, r.easyHandle(easyHandle), middlewares...)
}

// basic usage function

func (r *Router) GET(path string, handle Handle, middlewares ...Handle) *Router {
	return r.Handle("GET", path, handle, middlewares...)
}

func (r *Router) HEAD(path string, handle Handle, middlewares ...Handle) *Router {
	return r.Handle("HEAD", path, handle, middlewares...)
}

func (r *Router) OPTIONS(path string, handle Handle, middlewares ...Handle) *Router {
	return r.Handle("OPTIONS", path, handle, middlewares...)
}

func (r *Router) POST(path string, handle Handle, middlewares ...Handle) *Router {
	return r.Handle("POST", path, handle, middlewares...)
}

func (r *Router) PUT(path string, handle Handle, middlewares ...Handle) *Router {
	return r.Handle("PUT", path, handle, middlewares...)
}

func (r *Router) PATCH(path string, handle Handle, middlewares ...Handle) *Router {
	return r.Handle("PATCH", path, handle, middlewares...)
}

func (r *Router) DELETE(path string, handle Handle, middlewares ...Handle) *Router {
	return r.Handle("DELETE", path, handle, middlewares...)
}

var methodNames = []string{"GET", "HEAD", "OPTIONS", "POST", "PUT", "PATCH", "DELETE"}

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

func (r *Router) consoleStartPrint(addr string) {
	if r.closeConsolePrint {
		return
	}
	fmt.Println("  ______          _        __          __  _     \n |  ____|        (_)       \\ \\        / / | |    \n | |__   __ _ ___ _  ___ _ _\\ \\  /\\  / /__| |__  \n |  __| / _` / __| |/ _ \\ '__\\ \\/  \\/ / _ \\ '_ \\ \n | |___| (_| \\__ \\ |  __/ |   \\  /\\  /  __/ |_) |\n |______\\__,_|___/_|\\___|_|    \\/  \\/ \\___|_.__/")
	fmt.Printf("\033[1;32;40m%s\033[0m\n", fmt.Sprintf(" >>> http server runs on [%s] ", addr))
}
