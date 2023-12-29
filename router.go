package easierweb

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/websocket"
	"io"
	"net/http"
	"reflect"
	"strings"
)

type Router struct {
	contextPath            string
	multipartFormMaxMemory int64
	router                 *httprouter.Router
	server                 http.Server
	middlewares            []Handle
	errorHandle            ErrorHandle
	requestHandle          RequestHandle
	responseHandle         ResponseHandle
}

type Handle func(ctx *Context)

type EasyHandle func(ctx *Context, reqObj any) (any, error)

type RequestHandle func(ctx *Context, paramValues []reflect.Value) error

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

func (r *Router) SetEasyHandlePlugins(requestHandle RequestHandle, responseHandle ResponseHandle) *Router {
	r.requestHandle = requestHandle
	r.responseHandle = responseHandle
	return r
}

// --

func (r *Router) EasyGET(path string, easyHandle any) *Router {
	return r.ReEasyGET(path, easyHandle, r.requestHandle, r.responseHandle)
}

func (r *Router) EasyHEAD(path string, easyHandle any) *Router {
	return r.ReEasyHEAD(path, easyHandle, r.requestHandle, r.responseHandle)
}

func (r *Router) EasyOPTIONS(path string, easyHandle any) *Router {
	return r.ReEasyOPTIONS(path, easyHandle, r.requestHandle, r.responseHandle)
}

func (r *Router) EasyPOST(path string, easyHandle any) *Router {
	return r.ReEasyPOST(path, easyHandle, r.requestHandle, r.responseHandle)
}

func (r *Router) EasyPUT(path string, easyHandle any) *Router {
	return r.ReEasyPUT(path, easyHandle, r.requestHandle, r.responseHandle)
}

func (r *Router) EasyPATCH(path string, easyHandle any) *Router {
	return r.ReEasyPATCH(path, easyHandle, r.requestHandle, r.responseHandle)
}

func (r *Router) EasyDELETE(path string, easyHandle any) *Router {
	return r.ReEasyDELETE(path, easyHandle, r.requestHandle, r.responseHandle)
}

// --

func (r *Router) ReEasyGET(path string, easyHandle any, requestHandle RequestHandle, responseHandle ResponseHandle) *Router {
	return r.GET(path, r.easyHandle2handle(easyHandle, requestHandle, responseHandle))
}

func (r *Router) ReEasyHEAD(path string, easyHandle any, requestHandle RequestHandle, responseHandle ResponseHandle) *Router {
	return r.HEAD(path, r.easyHandle2handle(easyHandle, requestHandle, responseHandle))
}

func (r *Router) ReEasyOPTIONS(path string, easyHandle any, requestHandle RequestHandle, responseHandle ResponseHandle) *Router {
	return r.OPTIONS(path, r.easyHandle2handle(easyHandle, requestHandle, responseHandle))
}

func (r *Router) ReEasyPOST(path string, easyHandle any, requestHandle RequestHandle, responseHandle ResponseHandle) *Router {
	return r.POST(path, r.easyHandle2handle(easyHandle, requestHandle, responseHandle))
}

func (r *Router) ReEasyPUT(path string, easyHandle any, requestHandle RequestHandle, responseHandle ResponseHandle) *Router {
	return r.PUT(path, r.easyHandle2handle(easyHandle, requestHandle, responseHandle))
}

func (r *Router) ReEasyPATCH(path string, easyHandle any, requestHandle RequestHandle, responseHandle ResponseHandle) *Router {
	return r.PATCH(path, r.easyHandle2handle(easyHandle, requestHandle, responseHandle))
}

func (r *Router) ReEasyDELETE(path string, easyHandle any, requestHandle RequestHandle, responseHandle ResponseHandle) *Router {
	return r.DELETE(path, r.easyHandle2handle(easyHandle, requestHandle, responseHandle))
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

func (r *Router) easyHandle2handle(easyHandle any, requestHandle RequestHandle, responseHandle ResponseHandle) Handle {
	return func(ctx *Context) {
		// 如果为空
		if requestHandle == nil {
			panic(errors.New("request handle is empty"))
		}
		if responseHandle == nil {
			panic(errors.New("response handle is empty"))
		}
		// 反射获取函数类型
		funcType := reflect.TypeOf(easyHandle)

		// 创建参数值的切片
		var paramValues []reflect.Value

		// 如果没有第二个参数，就不进行自动绑定了
		if funcType.NumIn() == 1 {
			paramValues = make([]reflect.Value, 1)
			paramValues[0] = reflect.ValueOf(ctx).Elem().Addr()
		} else if funcType.NumIn() == 2 {
			paramValues = make([]reflect.Value, 2)
			paramValues[0] = reflect.ValueOf(ctx).Elem().Addr()
			paramValues[1] = reflect.New(funcType.In(1)).Elem()
		} else {
			panic(errors.New("response handle parameters does not match"))
		}

		err := requestHandle(ctx, paramValues)
		if err != nil {
			responseHandle(ctx, nil, err)
		}

		// 调用函数
		returnValues := reflect.ValueOf(easyHandle).Call(paramValues)

		// 无对象返回，无错误返回
		if len(returnValues) == 0 {
			responseHandle(ctx, nil, nil)
			return
		}

		// 处理返回值
		var resultValue any
		if returnValues[0].IsValid() && returnValues[0].Kind() == reflect.Ptr && returnValues[0].Elem().IsValid() {
			resultValue = returnValues[0].Elem().Interface()
		} else if returnValues[0].IsValid() && returnValues[0].Kind() == reflect.Slice {
			resultValue = returnValues[0].Interface()
		}

		// 无错误返回
		if len(returnValues) == 1 {
			responseHandle(ctx, resultValue, nil)
			return
		}

		// 有对象返回，有错误返回
		errValue, _ := returnValues[1].Interface().(error)
		responseHandle(ctx, resultValue, errValue)
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
