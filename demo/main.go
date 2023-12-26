package main

import (
	"easierweb"
	"fmt"
	"net/http"
)

func main() {

	// 新建路由
	router := easierweb.NewRouter("/test")

	// 添加中间件
	router.AddMiddleware(DemoMiddleware)

	// 添加处理方法
	router.GET("/demoGet/:id", DemoGet)
	router.POST("/demoPost", DemoPost)

	// 启动路由
	err := router.Run(":8082")
	if err != nil {
		panic(err)
	}
}

// DemoGet GET接口
// http://127.0.0.1:8082/test/demoGet/123?type=1&price=10.24&name=dpwgc
func DemoGet(ctx *easierweb.Context) {

	fmt.Println("id:", ctx.Path.GetInt64("id"))

	fmt.Println("type:", ctx.Query.GetInt("type"))
	fmt.Println("price:", ctx.Query.GetFloat64("price"))
	fmt.Println("name:", ctx.Query.GetString("name"))

	ctx.WriteJson(http.StatusOK, ResultDTO{
		Msg:  "hello world",
		Data: "GET Request",
	})
}

// DemoPost POST接口
// http://127.0.0.1:8082/test/demoPost
/*
{
  "id": 123,
  "name": "dpwgc"
}
*/
func DemoPost(ctx *easierweb.Context) {

	command := Command{}
	err := ctx.BindJson(&command)
	if err != nil {
		panic(err)
	}

	fmt.Println("id:", command.Id, ", name:", command.Name)

	ctx.WriteJson(http.StatusOK, ResultDTO{
		Msg:  "hello world",
		Data: "POST Request",
	})
}

// DemoMiddleware 中间件
func DemoMiddleware(ctx *easierweb.Context) {
	fmt.Println("request url:", ctx.Request.URL)
	ctx.Next()
}

// Command 请求命令
type Command struct {
	Name string `json:"name"`
	Id   int64  `json:"id"`
}

// ResultDTO 响应DTO
type ResultDTO struct {
	Msg  string `json:"msg"`
	Data string `json:"data"`
}
