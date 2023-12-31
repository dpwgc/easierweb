package easierweb

import (
	"fmt"
	"golang.org/x/net/websocket"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

// simple test

func TestRouterSimple(t *testing.T) {

	fmt.Println("\n[TestRouterSimple] start")

	router := New(RouterOptions{
		RootPath: "/test/simple",
	}).Use(simpleTestMiddleware)

	// set handles
	router.POST("/post", simpleTestAPI)
	router.DELETE("/delete/:id", simpleTestAPI)
	router.PUT("/put/:id", simpleTestAPI)
	router.PATCH("/patch/:id", simpleTestAPI)
	router.GET("/get/:id", simpleTestAPI)
	router.OPTIONS("/options/:id", simpleTestAPI)
	router.HEAD("/head/:id", simpleTestAPI)

	router.WS("/ws", simpleTestWebsocketConnect)

	router.GET("/error", simpleTestErrorAPI)

	// set easy handles
	router.EasyPOST("/easy/post", simpleTestEasySaveAPI)
	router.EasyDELETE("/easy/delete/:id", simpleTestEasyDelAPI)
	router.EasyPUT("/easy/put/:id", simpleTestEasySaveAPI)
	router.EasyPATCH("/easy/patch/:id", simpleTestEasySaveAPI)
	router.EasyGET("/easy/get/:id", simpleTestEasyQueryAPI)
	router.EasyOPTIONS("/easy/options/:id", simpleTestEasyQueryAPI)
	router.EasyHEAD("/easy/head/:id", simpleTestEasyQueryAPI)

	router.EasyGET("/easy/error", simpleTestErrorAPI)

	go func() {
		time.Sleep(3 * time.Second)
		simpleTestHttpSendExecute()
		fmt.Println()
		time.Sleep(1 * time.Second)
		simpleTestWebsocketClientExecute()
		time.Sleep(2 * time.Second)
		fmt.Println("\n[TestRouterSimple](go func) close router")
		err := router.Close()
		if err != nil {
			panic(err)
		}
	}()

	// started on port 80
	err := router.Run(":80")
	if err != nil && err.Error() != "http: Server closed" {
		panic(err)
	}

	fmt.Println("\n[TestRouterSimple] end")
}

// middleware
func simpleTestMiddleware(ctx *Context) {
	fmt.Println("[TestRouterSimple](simpleTestMiddleware) before ->", time.Now().UnixMilli())
	ctx.Next()
	fmt.Println("[TestRouterSimple](simpleTestMiddleware) after ->", time.Now().UnixMilli())
}

func simpleTestAPI(ctx *Context) {
	if ctx.Request.Method == "GET" || ctx.Request.Method == "HEAD" || ctx.Request.Method == "OPTIONS" || ctx.Request.Method == "PUT" || ctx.Request.Method == "PATCH" || ctx.Request.Method == "DELETE" {
		fmt.Println("[TestRouterSimple](simpleTestAPI) uri path id ->", ctx.Path.Int64("id"))
	}
	if ctx.Request.Method == "POST" || ctx.Request.Method == "PUT" || ctx.Request.Method == "PATCH" {
		dto := simpleTestDTO{}
		err := ctx.BindJSON(&dto)
		if err != nil {
			panic(err)
		}
		fmt.Println("[TestRouterSimple](simpleTestAPI) request body ->", dto)
	}
	if ctx.Request.Method == "GET" {
		dto := simpleTestDTO{}
		err := ctx.BindQuery(&dto)
		if err != nil {
			panic(err)
		}
		fmt.Println("[TestRouterSimple](simpleTestAPI) uri query parameters ->", dto)
	}
	if ctx.Request.Method == "HEAD" {
		ctx.NoContent(http.StatusNoContent)
		return
	}
	ctx.WriteJSON(http.StatusOK, simpleTestDTO{
		Int:     1,
		Int32:   2,
		Int64:   3,
		String:  "test",
		Float32: 1.1,
		Float64: 2.2,
	})
}

func simpleTestEasyQueryAPI(ctx *Context, dto simpleTestDTO) *simpleTestDTO {
	fmt.Println("[TestRouterSimple](simpleTestEasyQueryAPI) uri path id ->", ctx.Path.Int64("id"))
	fmt.Println("[TestRouterSimple](simpleTestEasyQueryAPI) uri query parameters ->", dto)
	if ctx.Request.Method == "HEAD" {
		return nil
	}
	return &simpleTestDTO{
		Int:     1,
		Int32:   2,
		Int64:   3,
		String:  "test",
		Float32: 1.1,
		Float64: 2.2,
	}
}

func simpleTestEasySaveAPI(ctx *Context, dto simpleTestDTO) (*simpleTestDTO, error) {
	if ctx.Request.Method == "PUT" || ctx.Request.Method == "PATCH" {
		fmt.Println("[TestRouterSimple](simpleTestEasySaveAPI) uri path id ->", ctx.Path.Int64("id"))
	}
	fmt.Println("[TestRouterSimple](simpleTestEasySaveAPI) request body ->", dto)
	return &simpleTestDTO{
		Int:     1,
		Int32:   2,
		Int64:   3,
		String:  "test",
		Float32: 1.1,
		Float64: 2.2,
	}, nil
}

func simpleTestEasyDelAPI(ctx *Context, dto simpleTestDTO) {
	fmt.Println("[TestRouterSimple](simpleTestEasyDelAPI) uri path id ->", ctx.Path.Int64("id"))
}

func simpleTestErrorAPI(ctx *Context) {
	panic("test error")
}

func simpleTestWebsocketConnect(ctx *Context) {
	msg, err := ctx.ReceiveString()
	if err != nil {
		panic(err)
	}
	fmt.Println("[TestRouterSimple](simpleTestWebsocketConnect) server read websocket msg ->", msg)
	err = ctx.SendString("test msg")
	if err != nil {
		panic(err)
	}
	return
}

func simpleTestWebsocketClientExecute() {
	origin := "http://localhost/"
	url := "ws://localhost/test/simple/ws"
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		panic(err)
	}
	_, err = ws.Write([]byte("test msg"))
	if err != nil {
		panic(err)
	}
	var buf = make([]byte, 100)
	n, err := ws.Read(buf)
	if err != nil {
		panic(err)
	}
	fmt.Println("[TestRouterSimple](simpleTestWebsocketClientExecute) client read websocket msg ->", string(buf[:n]))
}

func simpleTestHttpSendExecute() {

	body := "{\"int\":1,\"int32\":2,\"int64\":3,\"string\":\"test\",\"float32\":1.1,\"float64\":2.2}"

	simpleTestHttpClient("HEAD", "/head/123", "")
	simpleTestHttpClient("HEAD", "/easy/head/123", "")

	simpleTestHttpClient("OPTIONS", "/options/123", "")
	simpleTestHttpClient("OPTIONS", "/easy/options/123", "")

	simpleTestHttpClient("GET", "/get/123?int=1&int32=2&int64=3&string=test&float32=1.1&float64=2.2", "")
	simpleTestHttpClient("GET", "/easy/get/123?int=1&int32=2&int64=3&string=test&float32=1.1&float64=2.2", "")

	simpleTestHttpClient("POST", "/post", body)
	simpleTestHttpClient("POST", "/easy/post", body)

	simpleTestHttpClient("PUT", "/put/123", body)
	simpleTestHttpClient("PUT", "/easy/put/123", body)

	simpleTestHttpClient("PATCH", "/patch/123", body)
	simpleTestHttpClient("PATCH", "/easy/patch/123", body)

	simpleTestHttpClient("DELETE", "/delete/123", "")
	simpleTestHttpClient("DELETE", "/easy/delete/123", "")

	simpleTestHttpClient("GET", "/error", "")
	simpleTestHttpClient("GET", "/easy/error", "")
}

func simpleTestHttpClient(method, uri, body string) {
	fmt.Printf("\n[TestRouterSimple](simpleTestHttpClient) request method: %s, uri: %s, body: %s \n", method, uri, body)
	var payload = strings.NewReader(body)
	request, err := http.NewRequest(method, "http://localhost/test/simple"+uri, payload)
	if err != nil {
		panic(err)
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			panic(err)
		}
	}(response.Body)
	result, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	fmt.Printf("[TestRouterSimple](simpleTestHttpClient) response code: %v, data -> %s \n", response.StatusCode, string(result))
}

type simpleTestDTO struct {
	Int     int     `json:"int" schema:"int"`
	Int32   int32   `json:"int32" schema:"int32"`
	Int64   int64   `json:"int64" schema:"int64"`
	String  string  `json:"string" schema:"string"`
	Float32 float32 `json:"float32" schema:"float32"`
	Float64 float32 `json:"float64" schema:"float64"`
}
