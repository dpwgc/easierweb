package easierweb

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

// group test

func TestGroup(t *testing.T) {

	fmt.Println("\n[TestGroup] start")

	router := New(RouterOptions{
		RootPath: "/test/group",
	}).Use(groupTestMiddleware)

	group1 := router.Group("/1", groupTestMiddleware1)
	{
		group1.GET("/hello", groupTeatApi)
	}
	group2 := router.Group("/2", groupTestMiddleware2)
	{
		group2.GET("/hello", groupTeatApi)
	}
	group3 := router.Group("/3", groupTestMiddleware3)
	{
		group3.GET("/hello", groupTeatApi)
	}

	go func() {
		time.Sleep(3 * time.Second)
		groupTeatHttpClient("GET", "/1/hello", "")
		time.Sleep(1 * time.Second)
		groupTeatHttpClient("GET", "/2/hello", "")
		time.Sleep(1 * time.Second)
		groupTeatHttpClient("GET", "/3/hello", "")
		time.Sleep(1 * time.Second)
		fmt.Println("\n[TestGroup](go func) close router")
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

	fmt.Println("\n[TestGroup] end")
}

// middleware
func groupTestMiddleware(ctx *Context) {
	fmt.Println("[TestGroup](groupTestMiddleware) start")
	ctx.Next()
	fmt.Println("[TestGroup](groupTestMiddleware) end")
}

func groupTestMiddleware1(ctx *Context) {
	fmt.Println("[TestGroup](groupTestMiddleware1) start")
	ctx.Next()
	fmt.Println("[TestGroup](groupTestMiddleware1) end")
}

func groupTestMiddleware2(ctx *Context) {
	fmt.Println("[TestGroup](groupTestMiddleware2) start")
	ctx.Next()
	fmt.Println("[TestGroup](groupTestMiddleware2) end")
}

func groupTestMiddleware3(ctx *Context) {
	fmt.Println("[TestGroup](groupTestMiddleware3) start")
	ctx.Next()
	fmt.Println("[TestGroup](groupTestMiddleware3) end")
}

func groupTeatApi(ctx *Context) {
	fmt.Println("[TestGroup](groupTeatApi)", ctx.Method(), ctx.Route)
	ctx.WriteString(http.StatusOK, "hello")
}

func groupTeatHttpClient(method, uri, body string) {
	fmt.Printf("\n[TestGroup](groupTeatHttpClient) request method: %s, uri: %s, body: %s \n", method, uri, body)
	code, result, err := requestDo(method, "http://localhost/test/group"+uri, []byte(body))
	if err != nil {
		panic(err)
	}
	fmt.Printf("[TestGroup](groupTeatHttpClient) response code: %v, data -> %s \n", code, string(result))
}
