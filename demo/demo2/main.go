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

	// 更简单的接口写法，不需要手写请求Body、Query参数、Form参数绑定逻辑以及Write响应逻辑，由框架代理实现
	// 需要配置请求处理插件和响应处理插件，可在router上全局配置(如下面这一行，设置了两个插件：自动绑定Json请求数据、自动响应Json数据)
	// 也可在单个方法上单独配置插件
	router.SetEasyHandleDefaultPlugins(plugins.JSONRequestHandle, plugins.JSONResponseHandle)

	// POST接口
	// EasyPOST函数-使用router设置的全局请求响应处理插件
	// 服务处理函数-DemoEasyPost
	router.EasyPOST("/demoEasyPost", DemoEasyPost)

	// GET接口
	// ReEasyGET函数-单独配置这个方法的请求响应处理插件
	// 服务处理函数-DemoEasyGet
	// 请求处理插件-NoActionRequestHandle：不进行任何自动绑定操作
	// 响应处理插件-JSONResponseHandle：自动接收响应结果，并将结果序列化成Json字符串回传给客户端
	router.ReEasyGET("/demoEasyGet/:id", DemoEasyGet, plugins.NoActionRequestHandle, plugins.JSONResponseHandle)

	// 启动路由
	err := router.Run(":8081")
	if err != nil {
		panic(err)
	}
}

// DemoEasyPost POST接口-简易写法
// http://127.0.0.1:8081/test/demoEasyPost
/*
{
  "id": 123,
  "name": "dpwgc"
}
*/
func DemoEasyPost(ctx *easierweb.Context, command DemoCommand) (*DemoResultDTO, error) {

	// 打印参数
	fmt.Println("body -> id:", command.Id, ", name:", command.Name)

	// 返回
	return &DemoResultDTO{
		Msg:  "hello world",
		Data: "POST Request (Easy)",
	}, nil
}

// DemoEasyGet GET接口-简易写法
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

// DemoCommand 命令请求
type DemoCommand struct {
	Name string `json:"name"`
	Id   int64  `json:"id"`
}

// DemoResultDTO 响应DTO
type DemoResultDTO struct {
	Msg  string `json:"msg"`
	Data string `json:"data"`
}
