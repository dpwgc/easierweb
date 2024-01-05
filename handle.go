package easierweb

import (
	"errors"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/websocket"
	"io"
	"net/http"
	"reflect"
	"strings"
)

type Handle func(ctx *Context)

type RequestHandle func(ctx *Context, reqObj any) error

type ResponseHandle func(ctx *Context, result any, err error)

type ErrorHandle func(ctx *Context, err any)

func (r *Router) handle(route string, handle Handle, res http.ResponseWriter, req *http.Request, par httprouter.Params, ws *websocket.Conn, middlewares ...Handle) {

	ctx, err := r.buildContext(route, res, req, par, ws, middlewares...)

	defer func() {
		sErr := recover()
		if sErr != nil && r.errorHandle != nil {
			r.errorBottomUp(ctx, sErr)
		}
	}()

	if err != nil {
		panic(err)
	}

	// middleware execution
	ctx.handles = append(ctx.handles, handle)
	for ctx.index < len(ctx.handles) {
		ctx.handles[ctx.index](ctx)
		ctx.index++
	}

	// if a websocket connection exists, the websocket connection is automatically closed when the function returns
	if ws != nil {
		err = ws.Close()
		if err != nil {
			panic(err)
		}
	}
}

func (r *Router) easyHandle(easyHandle any) Handle {
	return func(ctx *Context) {
		// verify
		if r.requestHandle == nil {
			panic(errors.New("request handle is empty"))
		}
		if r.responseHandle == nil {
			panic(errors.New("response handle is empty"))
		}
		// reflection gets the type of function
		funcType := reflect.TypeOf(easyHandle)

		// create a slice of the parameter value
		var paramValues []reflect.Value

		var reqObj any = nil
		// if there is no second parameter, there is no auto-binding
		if funcType.NumIn() == 1 {
			paramValues = make([]reflect.Value, 1)
			paramValues[0] = reflect.ValueOf(ctx).Elem().Addr()
		} else if funcType.NumIn() == 2 {
			paramValues = make([]reflect.Value, 2)
			paramValues[0] = reflect.ValueOf(ctx).Elem().Addr()
			paramValues[1] = reflect.New(funcType.In(1)).Elem()
			reqObj = paramValues[1].Addr().Interface()
		} else {
			panic(errors.New("handle input parameters does not match"))
		}

		if reqObj != nil {
			err := r.requestHandle(ctx, reqObj)
			if err != nil {
				r.responseHandle(ctx, nil, err)
			}
		}

		// call the function
		returnValues := reflect.ValueOf(easyHandle).Call(paramValues)

		// no object return, no error return
		if len(returnValues) == 0 {
			r.responseHandle(ctx, nil, nil)
			return
		}

		if len(returnValues) > 2 {
			panic(errors.New("handle return values does not match"))
		}

		// process the return value
		var resultValue any = nil

		// if just one value return
		if len(returnValues) == 1 {
			firstValue, isErr := returnValues[0].Interface().(error)
			// if first return value is error
			if isErr {
				// return error
				r.responseHandle(ctx, nil, firstValue)
			} else {
				// return result
				r.responseHandle(ctx, resultValue, nil)
			}
			return
		} else if returnValues[0].IsValid() && returnValues[0].Kind() == reflect.Ptr && returnValues[0].Elem().IsValid() {
			resultValue = returnValues[0].Elem().Interface()
		} else if returnValues[0].IsValid() && returnValues[0].Kind() == reflect.Slice {
			resultValue = returnValues[0].Interface()
		}

		// has object return and error return
		errValue, _ := returnValues[1].Interface().(error)
		r.responseHandle(ctx, resultValue, errValue)
	}
}

func (r *Router) buildContext(route string, res http.ResponseWriter, req *http.Request, par httprouter.Params, ws *websocket.Conn, middlewares ...Handle) (*Context, error) {

	handles := append([]Handle(nil), r.middlewares...)
	handles = append(handles, middlewares...)
	ctx := Context{
		Route:          route,
		index:          0,
		handles:        handles,
		Header:         map[string]string{},
		Path:           map[string]string{},
		Query:          map[string]string{},
		Form:           map[string]string{},
		Body:           nil,
		Request:        req,
		ResponseWriter: res,
		WebsocketConn:  ws,
		Code:           0,
		Result:         nil,
		written:        false,
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

func (r *Router) errorBottomUp(ctx *Context, err any) {
	defer func() {
		_ = recover()
	}()
	r.errorHandle(ctx, err)
}
