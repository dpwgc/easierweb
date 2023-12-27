package main

import (
	"easierweb"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

// 路由示例
func main() {

	// 新建路由
	router := easierweb.New()

	// 设置根路径 /test
	router.SetContextPath("/test")

	// 添加中间件
	router.AddMiddleware(DemoMiddleware)

	// 常规服务（GET接口，POST接口，Websocket连接处理器，表单文件上传接口，文件下载接口）
	router.GET("/demoGet/:id", DemoGet)
	router.POST("/demoPost", DemoPost)
	router.WS("/demoWS/:id", DemoWS)
	router.POST("/demoUpload", DemoUpload)
	router.GET("/demoDownload/:fileName", DemoDownload)

	// 静态文件服务（访问demo目录） http://127.0.0.1:8082/test/demoStatic
	router.Static("/demoStatic/*filepath", "demo")

	// 设置错误处理器，捕获panic出来的异常
	router.SetErrorHandle(func(ctx *easierweb.Context, err any) {
		errMsg := fmt.Sprintf("%s", err)
		fmt.Println("err msg:", errMsg)
		// 返回code=500加异常信息
		ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		_, _ = ctx.ResponseWriter.Write([]byte(errMsg))
	})

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

	// 获取URI参数
	fmt.Println("id:", ctx.Path.Int64("id"))

	// 获取Query参数列表
	fmt.Println("query keys:", ctx.Query.Keys())
	fmt.Println("query values:", ctx.Query.Values())

	// 获取Query参数
	fmt.Println("type:", ctx.Query.Int("type"))
	fmt.Println("price:", ctx.Query.Float64("price"))
	fmt.Println("name:", ctx.Query["name"])

	// 返回
	ctx.WriteJson(http.StatusOK, DemoResultDTO{
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

	// 序列化请求体
	command := DemoCommand{}
	err := ctx.BindJson(&command)
	if err != nil {
		panic(err)
	}

	fmt.Println("body -> id:", command.Id, ", name:", command.Name)

	// 返回
	ctx.WriteJson(http.StatusOK, DemoResultDTO{
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
		// 读取消息
		var msg string
		// var msg []byte
		err := ctx.Receive(&msg)
		if err != nil {
			panic(err)
		}

		fmt.Println("read ws msg:", msg)

		// 发送消息
		err = ctx.SendJson(DemoResultDTO{
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
	ctx.WriteJson(http.StatusOK, DemoResultDTO{
		Msg:  "hello world",
		Data: "Upload File",
	})
}

// DemoDownload 文件下载接口
// http://127.0.0.1:8082/test/demoDownload/README.md
// 下载当前服务运行目录下的指定文件
func DemoDownload(ctx *easierweb.Context) {
	fileBytes, err := os.ReadFile(ctx.Path["fileName"])
	if err != nil {
		panic(err)
	}
	// 返回（Content-Type不传默认application/octet-stream）
	ctx.WriteFile("", fileBytes)
}

// DemoMiddleware 中间件
func DemoMiddleware(ctx *easierweb.Context) {

	// 处理前-打印URL
	fmt.Println("\nrequest url:", ctx.Request.URL.String())

	// 跳到下一个方法
	ctx.Next()

	// 处理后-打印响应结果
	fmt.Println("result:", string(ctx.Result))
}

// DemoCommand 请求命令
type DemoCommand struct {
	Name string `json:"name"`
	Id   int64  `json:"id"`
}

// DemoResultDTO 响应DTO
type DemoResultDTO struct {
	Msg  string `json:"msg"`
	Data string `json:"data"`
}
