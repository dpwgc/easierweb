package main

import (
	"fmt"
	"github.com/dpwgc/easierweb"
	"github.com/dpwgc/easierweb/plugins"
)

// 示例程序
func main() {

	// 新建路由
	router := easierweb.New()

	// 设置根路径 /test
	router.SetContextPath("/test")

	// 添加中间件
	router.Use(DemoMiddleware)

	// 更简单的接口写法，不需要手写请求Body/Query参数/Form参数绑定逻辑以及Write响应逻辑，由框架代理实现
	// 需要配置请求处理插件和响应处理插件，可在router上全局配置，也可在单个方法上定制化配置插件
	// 请求处理插件-JSONRequestHandle：解析Json Body数据与Query参数，并将两者的参数自动绑定到函数的第二个入参上
	// 响应处理插件-JSONResponseHandle：自动接收函数返回的响应对象，并将其序列化成Json字符串回传给客户端
	router.SetEasyHandlePlugins(plugins.JSONRequestHandle, plugins.JSONResponseHandle)

	// 使用router设置的全局请求/响应处理插件（Json + Json）
	router.EasyPOST("/demoEasyPost", DemoEasyPost)
	router.EasyGET("/demoEasyGet/:id", DemoEasyGet)

	// 局部定制化配置方法的请求/响应处理插件
	// 响应处理插件-BytesResponseHandle：自动接收函数返回的字节数组，并将其回传给客户端，该响应处理器适用于：返回字符串、返回文件、空响应体等场景
	router.ReEasyGET("/demoEasyGetString", DemoEasyGetString, plugins.JSONRequestHandle, plugins.BytesResponseHandle)
	router.ReEasyGET("/DemoEasyGetHtml", DemoEasyGetHTML, plugins.JSONRequestHandle, plugins.BytesResponseHandle)
	router.ReEasyGET("/demoEasyGetNoContent", DemoEasyGetNoContent, plugins.JSONRequestHandle, plugins.BytesResponseHandle)

	// 启动路由
	err := router.Run(":8081")
	if err != nil {
		panic(err)
	}
}

// DemoEasyPost POST接口(同时绑定URI Query参数和Json Body数据)-简易写法
// http://127.0.0.1:8081/test/demoEasyPost?type=1
/*
{
  "id": 123,
  "name": "dpwgc"
}
*/
func DemoEasyPost(ctx *easierweb.Context, command DemoCommand) (*DemoResultDTO, error) {

	// 打印Body数据
	fmt.Println("body -> id:", command.Id, ", name:", command.Name, ", type:", command.Type)
	// 打印Query参数
	fmt.Println("query -> type:", command.Type)

	// 返回
	return &DemoResultDTO{
		Msg:  "hello world",
		Data: "POST Request (Easy)",
	}, nil
}

// DemoEasyGet GET接口(无入参)-简易写法
// 如果该方法不需要自动绑定参数，第二个入参应该去除，只保留第一个入参Context
// http://127.0.0.1:8081/test/demoEasyGet/123
func DemoEasyGet(ctx *easierweb.Context) (*DemoResultDTO, error) {

	// 获取URI参数
	fmt.Println("id:", ctx.Path.Int64("id"))

	// 返回
	return &DemoResultDTO{
		Msg:  "hello world",
		Data: "GET Request (Easy)",
	}, nil
}

// DemoEasyGetString GET接口(返回字符串)-简易写法
// http://127.0.0.1:8081/test/demoEasyGetString
func DemoEasyGetString(ctx *easierweb.Context) (*[]byte, error) {
	// 以字节数组形式返回
	res := []byte("hello world")
	return &res, nil
}

// DemoEasyGetHTML GET接口(返回HTML)-简易写法
// http://127.0.0.1:8081/test/demoEasyGetHtml
func DemoEasyGetHTML(ctx *easierweb.Context) (*[]byte, error) {
	// 以字节数组形式返回
	res := []byte("<!DOCTYPE html>\n<html lang=\"en\">\n<head>\n    <meta charset=\"UTF-8\">\n    <title>Title</title>\n</head>\n<body>\n    hello world\n</body>\n</html>")
	return &res, nil
}

// DemoEasyGetNoContent GET接口(响应数据为空)-简易写法
// http://127.0.0.1:8081/test/demoEasyGetNoContent
func DemoEasyGetNoContent(ctx *easierweb.Context) (*[]byte, error) {
	return nil, nil
}

// DemoMiddleware 中间件
func DemoMiddleware(ctx *easierweb.Context) {

	// 处理前-打印URL
	fmt.Println("\nrequest url:", ctx.Request.URL.String())

	// 跳到下一个方法
	ctx.Next()

	// 处理后-打印响应结果
	fmt.Println("result:", ctx.Result.String())
}

// DemoCommand 命令请求
type DemoCommand struct {
	// URI Query参数接收，需配置schema标签
	Type int `schema:"type"`
	// Json Body数据接收，需配置json标签
	Name string `json:"name"`
	Id   int64  `json:"id"`
}

// DemoResultDTO 响应DTO
type DemoResultDTO struct {
	Msg  string `json:"msg"`
	Data string `json:"data"`
}
