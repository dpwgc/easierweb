package main

import (
	"easierweb"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

// 路由示例
func main() {

	// 新建路由
	router := easierweb.NewRouter("/test")

	// 添加中间件
	router.AddMiddleware(DemoMiddleware)

	// 添加处理方法
	// GET接口
	router.GET("/demoGet/:id", DemoGet)
	// POST接口
	router.POST("/demoPost", DemoPost)
	// Websocket连接
	router.GET("/demoWS/:id", DemoWS)

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
	fmt.Println("id:", ctx.Path.GetInt64("id"))

	// 获取Query参数列表
	fmt.Println("query keys:", ctx.Query.Keys())
	fmt.Println("query values:", ctx.Query.Values())

	// 获取Query参数
	fmt.Println("type:", ctx.Query.GetInt("type"))
	fmt.Println("price:", ctx.Query.GetFloat64("price"))
	fmt.Println("name:", ctx.Query["name"])

	// 返回
	ctx.WriteJsonResult(http.StatusOK, ResultDTO{
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
	err := json.Unmarshal(ctx.Body, &command)
	if err != nil {
		panic(err)
	}

	fmt.Println("id:", command.Id, ", name:", command.Name)

	// 返回
	ctx.WriteJsonResult(http.StatusOK, ResultDTO{
		Msg:  "hello world",
		Data: "POST Request",
	})
}

// DemoWS Websocket连接
// ws://127.0.0.1:8082/test/demoWS/123
func DemoWS(ctx *easierweb.Context) {

	// 获取URI参数
	fmt.Println("id:", ctx.Path.GetInt64("id"))

	// 将HTTP升级为WebSocket
	upGrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upGrader.Upgrade(ctx.ResponseWriter, ctx.Request, nil)
	if err != nil {
		panic(err)
	}

	// 处理WebSocket连接
	go func() {
		for {
			// 读取消息
			messageType, p, err := conn.ReadMessage()
			if err != nil {
				panic(err)
			}

			fmt.Println("read ws msg:", string(p))

			// 发送消息
			err = conn.WriteMessage(messageType, []byte("Hello"))
			if err != nil {
				panic(err)
			}

			time.Sleep(3 * time.Second)

			// 关闭连接
			fmt.Println("close ws conn:")
			_ = conn.Close()
			return
		}
	}()
}

// DemoMiddleware 中间件
func DemoMiddleware(ctx *easierweb.Context) {

	fmt.Println("request url:", ctx.Request.URL)
	fmt.Println("remote addr:", ctx.Request.RemoteAddr)

	// 跳到下一个方法
	ctx.Next()

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
