# EasierWeb

## 更方便使用、高性能、高可定制的 Go Web 框架

### 基于 [httprouter](https://github.com/julienschmidt/httprouter)

***

## 功能
* 具有更简洁的 API 代码编写方式，可以自动绑定请求数据、自动写入响应数据。
* 内置了 Websocket连接、SSE推送、文件服务、HTTP 客户端等功能。
* 提供 API 分组功能，自定义根级、组级、函数级中间件。
* 高度可定制，可以自定义错误捕获和请求/响应数据处理。
* 没有依赖过多的第三方包，架构简单。
* 支持 TLS 和 HTTP2.

***

一个简单的API示例

```go
// 自动绑定查询参数/表单参数/Body数据，自动写入响应数据
func helloAPI(ctx *easierweb.Context, request HelloRequest) (*HelloResponse, error) {

   // 打印请求数据
   fmt.Println("request data ->", request)

   // 只需要返回处理结果和错误值，框架会帮助你写入响应数据
   return &HelloResponse{Code: 1000, Msg:  "hello"}, nil
}
```

***

## 安装

```
go get github.com/dpwgc/easierweb
```

***

## 示例

### 框架提供了两种不同风格的API编写方式

### `1` 基础用法 : 类似 gin 和 echo（没有用到反射，性能较高）
### `2` 简单用法 : 类似 spring boot（更方便编写API代码）

#### [函数用法说明](./DESC.md)

***

## 基础用法

### API 函数只有一个上下文参数

### 需要手动绑定请求数据、手动写入响应数据

```go
package main

import (
   "github.com/dpwgc/easierweb"
   "github.com/dpwgc/easierweb/middlewares"
   "log"
   "net/http"
)

// 基础用法示例
func main() {

   // 创建一个路由并设置中间件
   router := easierweb.New().Use(middlewares.Logger())

   // 添加一个API
   router.GET("/hello/:name", hello)

   // 在80端口上运行
   log.Fatal(router.Run(":80"))
}

func hello(ctx *easierweb.Context) {

   // 获取URI路径参数
   name := ctx.Path.Get("name")

   // 写入响应，返回 'hello ' + 'name'
   ctx.WriteString(http.StatusOK, "hello "+name)
}

```

### 在浏览器中访问链接

> `GET` http://localhost/hello/easierweb

### 响应数据

```
hello easierweb
```

***

## 简单用法

### API 函数可以有传入参数和返回值

### 不需要手动绑定请求参数、手动写入响应参数，框架将替你完成这两步

```go
package main

import (
   "fmt"
   "github.com/dpwgc/easierweb"
   "github.com/dpwgc/easierweb/middlewares"
   "log"
)

// 简单用法
func main() {
	
   // 创建一个路由并设置中间件
   router := easierweb.New().Use(middlewares.Logger())
   
   // 添加一个 API (使用 'EasyXXX' 函数添加)
   router.EasyPOST("/submit", submit)
   
   // 在80端口上运行
   log.Fatal(router.Run(":80"))
}

func submit(ctx *easierweb.Context, req Request) *Response {
	
   // 打印请求数据
   fmt.Printf("post request data (json body) -> name: %s, mobile: %s \n", req.Name, req.Mobile)
   
   // 返回处理结果
   return &Response{Code: 1000, Msg:  "hello"}
}

// Request 请求Body数据对应的结构体
type Request struct {
   Name   string `json:"name"`
   Mobile string `json:"mobile"`
}

// Response 响应数据对应的结构体
type Response struct {
   Code int    `json:"code"`
   Msg  string `json:"msg"`
}
```

### 调用 HTTP POST API

> `POST` http://localhost/submit

### 请求体

```json
{
  "name": "hello",
  "mobile": "12345678"
}
```

### 响应数据

```json
{
  "code": 1000,
  "msg": "hello"
}
```

### 其他说明

* 如果想使用 'EasyXXX' 系列函数，传入的 API 函数必须遵循如下格式。

```go
// 第一个传入参数必须为 Context
// 请求结构体和响应结构体可以是切片 ([]Request/*[]Response)
func TestAPI(ctx *easierweb.Context, req Request) (*Response, error)
func TestAPI(ctx *easierweb.Context, req Request) *Response
func TestAPI(ctx *easierweb.Context, req Request) error
func TestAPI(ctx *easierweb.Context, req Request)
func TestAPI(ctx *easierweb.Context) (*Response, error)
func TestAPI(ctx *easierweb.Context) *Response
func TestAPI(ctx *easierweb.Context) error
func TestAPI(ctx *easierweb.Context)

// 添加 API
router.EasyPOST("/test", TestAPI)
```

* 框架默认使用 JSON 格式来处理请求和响应数据。
* 如果想改变数据的处理格式，可以像下文这样使用框架自带的插件，或者自己编写自定义请求/响应处理方法。

```go
// 使用 XML 格式解析请求数据、写入响应数据
router := easierweb.New(easierweb.RouterOptions{
   RequestHandle: plugins.XMLRequestHandle(),
   ResponseHandle: plugins.XMLResponseHandle(),
})
```

* 'EasyXXX' 系列函数兼容基础用法，也可以使用 'WriteXXX' 系列函数来写入响应数据。

```go
func TestAPI(ctx *easierweb.Context, req Request) {

   fmt.Println("request body ->", req)
   
   // 写入 XML 格式的响应体
   ctx.WriteXML(http.StatusOK, Response{Code: 1000, Msg:  "hello"})
}

// 添加 API
router.EasyPOST("/test", TestAPI)
```

***

## 演示程序

* demos
    * basic `基础用法演示`
        * main.go
    * easier `简单用法演示`
        * main.go
    * restful `restful应用演示`
        * main.go
    * customize `自定义框架配置演示`
        * main.go