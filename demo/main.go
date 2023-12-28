package main

import (
	"fmt"
	"github.com/dpwgc/easierweb"
	"github.com/dpwgc/easierweb/utils"
	"mime/multipart"
	"net/http"
	"time"
)

// 示例程序
func main() {

	// 新建路由
	router := easierweb.New()

	// 设置根路径 /test
	router.SetContextPath("/test")

	// 添加中间件
	router.Use(DemoMiddleware)

	// 常规服务（GET接口，POST接口，Websocket连接处理器，表单文件上传接口，文件下载接口）
	router.GET("/demoGet/:id", DemoGet)
	router.POST("/demoPost", DemoPost)
	router.WS("/demoWS/:id", DemoWS)
	router.POST("/demoUpload", DemoUpload)
	router.GET("/demoDownload/:fileName", DemoDownload)

	// 静态文件服务（访问demo目录） http://127.0.0.1:8082/test/demoStatic
	router.Static("/demoStatic/*filepath", "demo")

	// 更简单的接口写法
	// 接口处理函数直接return结果，响应处理器接收结果并响应客户端
	router.SimpleGET("/demoSimpleGet/:id", DemoSimpleGet, utils.JSONResponseHandle)

	// 启动路由
	err := router.Run(":8082")
	// 启动TLS路由
	/*
		router.RunTLS("0.0.0.0:443", "cert.pem", "private.key", &tls.Config{
			ClientAuth: tls.NoClientCert,
		})
	*/
	if err != nil {
		panic(err)
	}
}

// DemoGet GET接口
// http://127.0.0.1:8082/test/demoGet/123?type=1&price=10.24&name=dpwgc
func DemoGet(ctx *easierweb.Context) {

	// 上下文缓存数据读取
	cache := ctx.CustomCache.Get("test_cache2")
	if cache != nil {
		fmt.Println("context cache:", cache.(int))
	}

	// 获取URI参数
	fmt.Println("id:", ctx.Path.Int64("id"))

	// 获取Query参数列表
	fmt.Println("query keys:", ctx.Query.Keys())
	fmt.Println("query values:", ctx.Query.Values())

	// 将Query参数绑定到结构体上
	query := DemoQuery{}
	err := ctx.BindQuery(&query)
	if err != nil {
		panic(err)
	}
	fmt.Println("type:", query.Type)
	fmt.Println("price:", query.Price)
	fmt.Println("name:", query.Name)

	// 返回
	ctx.WriteJSON(http.StatusOK, DemoResultDTO{
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

	// 将Body数据绑定到结构体上
	command := DemoCommand{}
	err := ctx.BindJSON(&command)
	if err != nil {
		panic(err)
	}

	fmt.Println("body -> id:", command.Id, ", name:", command.Name)

	// 返回
	ctx.WriteJSON(http.StatusOK, DemoResultDTO{
		Msg:  "hello world",
		Data: "POST Request",
	})
}

// DemoWS Websocket连接
// ws://127.0.0.1:8082/test/demoWS/123
func DemoWS(ctx *easierweb.Context) {

	// 获取URI参数
	fmt.Println("id:", ctx.Path.Int64("id"))

	// 处理WebSocket连接
	for {
		// 读取字符串消息
		msg, err := ctx.ReceiveString()
		// 读取字节消息
		// msg, err := ctx.Receive()
		if err != nil {
			panic(err)
		}

		fmt.Println("read ws msg:", msg)

		// 发送消息
		err = ctx.SendJSON(DemoResultDTO{
			Msg:  "hello world",
			Data: "Websocket Connect",
		})
		if err != nil {
			panic(err)
		}

		time.Sleep(3 * time.Second)

		// 函数return后会自动关闭连接
		fmt.Println("close ws conn:")
		return
	}
}

// DemoUpload 文件上传接口
// http://127.0.0.1:8082/test/demoUpload
// Form表单参数
// 'file' -> '文件'
func DemoUpload(ctx *easierweb.Context) {

	fmt.Println("file keys:", ctx.FileKeys())

	// 获取表单文件
	file, err := ctx.GetFile("file")
	if err != nil {
		panic(err)
	}
	defer func(file multipart.File) {
		err = file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	fmt.Println("file:", file)

	// 返回
	ctx.WriteJSON(http.StatusOK, DemoResultDTO{
		Msg:  "hello world",
		Data: "Upload File",
	})
}

// DemoDownload 文件下载接口
// http://127.0.0.1:8082/test/demoDownload/README.md
// 下载当前服务运行目录下的指定文件
func DemoDownload(ctx *easierweb.Context) {

	// 获取本地文件并返回数据（contentType参数不传：默认application/octet-stream，fileName参数不传：下载后文件名默认为时间戳）
	ctx.WriteLocalFile("", ctx.Path.Get("fileName"), ctx.Path.Get("fileName"))

	// 直接返回文件字节数据
	// ctx.WriteFile("", ctx.Path.Get("fileName"), []byte{})
}

// DemoSimpleGet GET接口-简易写法
// http://127.0.0.1:8082/test/demoSimpleGet/123
func DemoSimpleGet(ctx *easierweb.Context) (any, error) {

	// 获取URI参数
	fmt.Println("id:", ctx.Path.Int64("id"))

	// 返回
	return DemoResultDTO{
		Msg:  "hello world",
		Data: "GET Request (Simple)",
	}, nil
}

// DemoMiddleware 中间件
func DemoMiddleware(ctx *easierweb.Context) {

	// 自定义缓存，可跟随Context传递到下层
	ctx.CustomCache.Set("test_cache1", "aaa").Set("test_cache2", 222)

	// 处理前-打印URL
	fmt.Println("\nrequest url:", ctx.Request.URL.String())

	// 跳到下一个方法
	ctx.Next()

	// 处理后-打印响应结果
	fmt.Println("result:", ctx.Result.String())
}

// DemoCommand 命令请求
type DemoCommand struct {
	Name string `json:"name"`
	Id   int64  `json:"id"`
}

// DemoQuery 查询请求
type DemoQuery struct {
	Type  int     `schema:"type"`
	Price float64 `schema:"price"`
	Name  string  `schema:"name"`
}

// DemoResultDTO 响应DTO
type DemoResultDTO struct {
	Msg  string `json:"msg"`
	Data string `json:"data"`
}
