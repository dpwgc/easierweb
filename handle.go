package easierweb

import (
	"errors"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/websocket"
	"net/http"
	"reflect"
)

type Handle func(ctx *Context)

type RequestHandle func(ctx *Context, reqObj any) error

type ResponseHandle func(ctx *Context, result any, err error)

type ErrorHandle func(ctx *Context, err any)

func (r *Router) handle(route string, handle Handle, res http.ResponseWriter, req *http.Request, par httprouter.Params, ws *websocket.Conn, sse bool, middlewares ...Handle) {

	ctx := r.contextPool.Get().(*Context)

	err := setContext(ctx, r, route, res, req, par, ws, middlewares...)

	defer func() {
		if ctx != nil {
			r.contextPool.Put(ctx)
		}
		sErr := recover()
		if sErr != nil && r.errorHandle != nil {
			r.errorBottomUp(ctx, sErr)
		}
	}()

	if err != nil {
		panic(err)
	}

	if sse {
		res.Header().Set("Content-Type", "text/event-stream")
		res.Header().Set("Cache-Control", "no-cache")
		res.Header().Set("Connection", "keep-alive")
		res.Header().Set("Access-Control-Allow-Origin", "*")
		flusher, ok := res.(http.Flusher)
		if !ok {
			panic("client does not support server-sent events")
		}
		ctx.Flusher = flusher
	}

	// middleware execution
	ctx.handles = append(ctx.handles, handle)
	for ctx.index < len(ctx.handles) {
		ctx.handles[ctx.index](ctx)
		ctx.index++
	}

	// if a websocket connection exists, the websocket connection is automatically closed when the function returns
	if ws != nil {
		err = ctx.Close()
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

		// if just one value return
		if len(returnValues) == 1 {
			firstValue, isErr := returnValues[0].Interface().(error)
			// if first return value is error
			if isErr {
				// return error
				r.responseHandle(ctx, nil, firstValue)
				return
			}
		}

		// get the result value
		var resultValue any = nil

		if returnValues[0].IsValid() && returnValues[0].Kind() == reflect.Ptr && returnValues[0].Elem().IsValid() {
			resultValue = returnValues[0].Elem().Interface()
		} else if returnValues[0].IsValid() && returnValues[0].Kind() == reflect.Slice {
			resultValue = returnValues[0].Interface()
		}

		// just result return
		if len(returnValues) == 1 {
			r.responseHandle(ctx, resultValue, nil)
			return
		}

		// has result return and error return
		errValue, _ := returnValues[1].Interface().(error)
		r.responseHandle(ctx, resultValue, errValue)
	}
}

func (r *Router) errorBottomUp(ctx *Context, err any) {
	defer func() {
		_ = recover()
	}()
	r.errorHandle(ctx, err)
}
