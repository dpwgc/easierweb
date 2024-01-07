package easierweb

import (
	"bytes"
	"errors"
	"fmt"
	"golang.org/x/net/websocket"
	"io"
	"net/http"
	"testing"
	"time"
)

// router test

func TestRouter(t *testing.T) {

	fmt.Println("\n[TestRouter] start")

	router := New(RouterOptions{
		RootPath: "/test/router",
	}).Use(routerTestMiddleware)

	// set handles
	router.POST("/post", routerTestAPI)
	router.DELETE("/delete/:id", routerTestAPI)
	router.PUT("/put/:id", routerTestAPI)
	router.PATCH("/patch/:id", routerTestAPI)
	router.GET("/get/:id", routerTestAPI)
	router.OPTIONS("/options/:id", routerTestAPI)
	router.HEAD("/head/:id", routerTestAPI)

	router.WS("/ws", routerTestWebsocketConnect)

	router.GET("/error", routerTestErrorAPI)

	// set easy handles
	router.EasyPOST("/easy/post", routerTestEasySaveAPI)
	router.EasyDELETE("/easy/delete/:id", routerTestEasyDelAPI)
	router.EasyPUT("/easy/put/:id", routerTestEasySaveAPI)
	router.EasyPATCH("/easy/patch/:id", routerTestEasySaveAPI)
	router.EasyGET("/easy/get/:id", routerTestEasyQueryAPI)
	router.EasyOPTIONS("/easy/options/:id", routerTestEasyQueryAPI)
	router.EasyHEAD("/easy/head/:id", routerTestEasyQueryAPI)

	router.EasyGET("/easy/error", routerTestErrorAPI)
	router.EasyGET("/easy/error/return", routerTestErrorReturnAPI)

	go func() {
		time.Sleep(3 * time.Second)
		routerTestHttpSendExecute()
		fmt.Println()
		time.Sleep(1 * time.Second)
		routerTestWebsocketClientExecute()
		time.Sleep(2 * time.Second)
		fmt.Println("\n[TestRouter](go func) close router")
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

	fmt.Println("\n[TestRouter] end")
}

// middleware
func routerTestMiddleware(ctx *Context) {
	fmt.Println("[TestRouter](routerTestMiddleware) route ->", ctx.Route)
	fmt.Println("[TestRouter](routerTestMiddleware) before ->", time.Now().UnixMilli())
	ctx.Next()
	fmt.Println("[TestRouter](routerTestMiddleware) after ->", time.Now().UnixMilli())
}

func routerTestAPI(ctx *Context) {
	if ctx.Request.Method == "GET" || ctx.Request.Method == "HEAD" || ctx.Request.Method == "OPTIONS" || ctx.Request.Method == "PUT" || ctx.Request.Method == "PATCH" || ctx.Request.Method == "DELETE" {
		fmt.Println("[TestRouter](routerTestAPI) uri path id ->", ctx.Path.Int64("id"))
	}
	if ctx.Request.Method == "POST" || ctx.Request.Method == "PUT" || ctx.Request.Method == "PATCH" {
		dto := routerTestDTO{}
		err := ctx.BindJSON(&dto)
		if err != nil {
			panic(err)
		}
		if dto.Int64 == 0 {
			panic("bind data error")
		}
		fmt.Println("[TestRouter](routerTestAPI) request body ->", dto)
	}
	if ctx.Request.Method == "GET" {
		dto := routerTestDTO{}
		err := ctx.BindQuery(&dto)
		if err != nil {
			panic(err)
		}
		if dto.Int64 == 0 {
			panic("bind data error")
		}
		fmt.Println("[TestRouter](routerTestAPI) uri query parameters ->", dto)
	}
	if ctx.Request.Method == "HEAD" {
		ctx.NoContent(http.StatusNoContent)
		return
	}
	ctx.WriteJSON(http.StatusOK, routerTestDTO{
		Int:     1,
		Int32:   2,
		Int64:   3,
		String:  "test",
		Float32: 1.1,
		Float64: 2.2,
	})
}

func routerTestEasyQueryAPI(ctx *Context, dto routerTestDTO) *routerTestDTO {
	fmt.Println("[TestRouter](routerTestEasyQueryAPI) uri path id ->", ctx.Path.Int64("id"))
	fmt.Println("[TestRouter](routerTestEasyQueryAPI) uri query parameters ->", dto)
	if ctx.Request.Method == "HEAD" {
		return nil
	}
	return &routerTestDTO{
		Int:     1,
		Int32:   2,
		Int64:   3,
		String:  "test",
		Float32: 1.1,
		Float64: 2.2,
	}
}

func routerTestEasySaveAPI(ctx *Context, dto routerTestDTO) (*routerTestDTO, error) {
	if ctx.Request.Method == "PUT" || ctx.Request.Method == "PATCH" {
		fmt.Println("[TestRouter](routerTestEasySaveAPI) uri path id ->", ctx.Path.Int64("id"))
	}
	fmt.Println("[TestRouter](routerTestEasySaveAPI) request body ->", dto)
	if dto.Int64 == 0 {
		panic("bind data error")
	}
	return &routerTestDTO{
		Int:     1,
		Int32:   2,
		Int64:   3,
		String:  "test",
		Float32: 1.1,
		Float64: 2.2,
	}, nil
}

func routerTestEasyDelAPI(ctx *Context, dto routerTestDTO) {
	fmt.Println("[TestRouter](routerTestEasyDelAPI) uri path id ->", ctx.Path.Int64("id"))
}

func routerTestErrorAPI(ctx *Context) {
	panic("test error")
}

func routerTestErrorReturnAPI(ctx *Context) error {
	return errors.New("test error return")
}

func routerTestWebsocketConnect(ctx *Context) {
	msg, err := ctx.ReceiveString()
	if err != nil {
		panic(err)
	}
	fmt.Println("[TestRouter](routerTestWebsocketConnect) server read websocket msg ->", msg)
	err = ctx.SendString("test msg")
	if err != nil {
		panic(err)
	}
	return
}

func routerTestWebsocketClientExecute() {
	origin := "http://localhost/"
	url := "ws://localhost/test/router/ws"
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
	fmt.Println("[TestRouter](routerTestWebsocketClientExecute) client read websocket msg ->", string(buf[:n]))
}

func routerTestHttpSendExecute() {

	body := "{\"int\":1,\"int32\":2,\"int64\":3,\"string\":\"test\",\"float32\":1.1,\"float64\":2.2}"

	routerTestHttpClient("HEAD", "/head/123", "")
	routerTestHttpClient("HEAD", "/easy/head/123", "")

	routerTestHttpClient("OPTIONS", "/options/123", "")
	routerTestHttpClient("OPTIONS", "/easy/options/123", "")

	routerTestHttpClient("GET", "/get/123?int=1&int32=2&int64=3&string=test&float32=1.1&float64=2.2", "")
	routerTestHttpClient("GET", "/easy/get/123?int=1&int32=2&int64=3&string=test&float32=1.1&float64=2.2", "")

	routerTestHttpClient("POST", "/post", body)
	routerTestHttpClient("POST", "/easy/post", body)

	routerTestHttpClient("PUT", "/put/123", body)
	routerTestHttpClient("PUT", "/easy/put/123", body)

	routerTestHttpClient("PATCH", "/patch/123", body)
	routerTestHttpClient("PATCH", "/easy/patch/123", body)

	routerTestHttpClient("DELETE", "/delete/123", "")
	routerTestHttpClient("DELETE", "/easy/delete/123", "")

	routerTestHttpClient("GET", "/error", "")
	routerTestHttpClient("GET", "/easy/error", "")
	routerTestHttpClient("GET", "/easy/error/return", "")
}

func routerTestHttpClient(method, uri, body string) {
	fmt.Printf("\n[TestRouter](routerTestHttpClient) request method: %s, uri: %s, body: %s \n", method, uri, body)
	code, result, err := requestDo(method, "http://localhost/test/router"+uri, []byte(body))
	if err != nil {
		panic(err)
	}
	fmt.Printf("[TestRouter](routerTestHttpClient) response code: %v, data -> %s \n", code, string(result))
}

func requestDo(method, url string, body []byte, header ...map[string]string) (int, Data, error) {
	request, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return 0, nil, err
	}
	if len(header) > 0 {
		for _, h := range header {
			for k, v := range h {
				request.Header.Set(k, v)
			}
		}
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return 0, nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)
	result, err := io.ReadAll(response.Body)
	if err != nil {
		return 0, nil, err
	}
	return response.StatusCode, result, nil
}

type routerTestDTO struct {
	Int     int     `json:"int" mapstructure:"int"`
	Int32   int32   `json:"int32" mapstructure:"int32"`
	Int64   int64   `json:"int64" mapstructure:"int64"`
	String  string  `json:"string" mapstructure:"string"`
	Float32 float32 `json:"float32" mapstructure:"float32"`
	Float64 float32 `json:"float64" mapstructure:"float64"`
}
