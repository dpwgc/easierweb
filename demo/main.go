package main

import (
	"easierweb"
	"fmt"
	"net/http"
	"time"
)

// 路由示例
func main() {

	// 新建路由
	router := easierweb.NewRouter("/test")

	// 添加中间件
	router.AddMiddleware(DemoMiddleware)

	// 添加处理方法（GET接口，POST接口，Websocket处理器）
	router.GET("/demoGet/:id", DemoGet)
	router.POST("/demoPost", DemoPost)
	router.WS("/demoWS/:id", DemoWS)

	// 设置错误处理器
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

	// 序列化请求体
	command := Command{}
	err := ctx.BindJson(&command)
	if err != nil {
		panic(err)
	}

	fmt.Println("id:", command.Id, ", name:", command.Name)

	// 返回
	ctx.WriteJson(http.StatusOK, ResultDTO{
		Msg:  "hello world",
		Data: "POST Request",
	})
}

// DemoWS Websocket连接
// ws://127.0.0.1:8082/test/demoWS/123
func DemoWS(ctx *easierweb.Context) {

	// 获取URI参数
	fmt.Println("id:", ctx.Path.Int64("id"))

	go func() {
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
			err = ctx.SendJson(ResultDTO{
				Msg:  "hello world",
				Data: "Websocket Connect",
			})
			if err != nil {
				panic(err)
			}

			time.Sleep(3 * time.Second)

			// 关闭连接
			fmt.Println("close ws conn:")
			_ = ctx.WebsocketConn.Close()
			return
		}
	}()
}

// DemoMiddleware 中间件
func DemoMiddleware(ctx *easierweb.Context) {

	// 处理前-打印请求参数
	fmt.Println("request url:", ctx.Request.URL)
	fmt.Println("remote addr:", ctx.Request.RemoteAddr)

	// 跳到下一个方法
	ctx.Next()

	// 处理前-打印响应结果
	fmt.Println("result:", string(ctx.Result))
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
