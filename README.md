# EasierWeb

## A more user-friendly, high performance, high customizable Go Web framework

### Based on [httprouter](https://github.com/julienschmidt/httprouter)

### [中文文档](./README_CH.md)

***

## Features
* Have a more concise way to write API. Can automatic binding query/form/body data and serialization writes to the response body.
* Built-in websocket, server-sent events (SSE), file server.
* Group APIs. Custom root-level, group-level, function-level middleware.
* Highly customizable. Custom error, request and response handle functions.
* No dependencies on too many third-party packages. Architecture is simple.
* Support TLS and HTTP2.

***

Simple example of API handle

```go
// automatic binding query/form/body data and serialization writes to the response body
func helloAPI(ctx *easierweb.Context, request HelloRequest) (*HelloResponse, error) {

   // print the request data
   fmt.Println("request data ->", request)

   // return result and error. framework will help you write the response
   return &HelloResponse{Code: 1000, Msg:  "hello"}, nil
}
```

***

## Installation

```
go get github.com/dpwgc/easierweb
```

* Specify version

```
go get github.com/dpwgc/easierweb@v1.0.4.15
```

***

## Example

### Framework provides two different styles of usage

### `1` Basic usage : like gin and echo ( no reflect, fast )
### `2` Easier usage : like spring boot ( more concise way to write API handle )

### [Function Usage Document](./DOC.md)

***

## Basic usage

### API handle function have only one context parameter

### Need to manually bind request data and write response data

```go
package main

import (
   "github.com/dpwgc/easierweb"
   "github.com/dpwgc/easierweb/middlewares"
   "log"
   "net/http"
)

// basic usage example
func main() {

   // create a router and set middleware
   router := easierweb.New().Use(middlewares.Logger())

   // set api handle
   router.GET("/hello/:name", hello)

   // runs on port 80
   log.Fatal(router.Run(":80"))
}

// api handle
func hello(ctx *easierweb.Context) {

   // get the path parameters
   name := ctx.Path.Get("name")

   // Write response, return 'hello ' + 'name'
   ctx.WriteString(http.StatusOK, "hello "+name)
}

```

### Access the HTTP URL in your browser

> `GET` http://localhost/hello/easierweb

### Response data

```
hello easierweb
```

***

## Easier usage

### API handle function can have input object and return values

### Don't need to manually bind request data and write response data. framework will help you do this

```go
package main

import (
   "fmt"
   "github.com/dpwgc/easierweb"
   "github.com/dpwgc/easierweb/middlewares"
   "log"
)

// easier usage example
func main() {
	
   // create a router and set middleware
   router := easierweb.New().Use(middlewares.Logger())
   
   // set api handle (use function 'EasyXXX')
   router.EasyPOST("/submit", submit)
   
   // runs on port 80
   log.Fatal(router.Run(":80"))
}

// api handle
func submit(ctx *easierweb.Context, req Request) *Response {
	
   // print the request data
   fmt.Printf("post request data (json body) -> name: %s, mobile: %s \n", req.Name, req.Mobile)
   
   // return result
   return &Response{Code: 1000, Msg:  "hello"}
}

// Request json body data
type Request struct {
   Name   string `json:"name"`
   Mobile string `json:"mobile"`
}

// Response json result data
type Response struct {
   Code int    `json:"code"`
   Msg  string `json:"msg"`
}
```

### Invoke HTTP POST API

> `POST` http://localhost/submit

### Request body

```json
{
  "name": "hello",
  "mobile": "12345678"
}
```

### Response data

```json
{
  "code": 1000,
  "msg": "hello"
}
```

### Other notes

* If you want to use 'EasyXXX' series functions. The api handle function must be in the following formats.

```go
// first input parameter must be Context
// request struct and response struct can be slices ([]Request/*[]Response)
func TestAPI(ctx *easierweb.Context, req Request) (*Response, error)
func TestAPI(ctx *easierweb.Context, req Request) *Response
func TestAPI(ctx *easierweb.Context, req Request) error
func TestAPI(ctx *easierweb.Context, req Request)
func TestAPI(ctx *easierweb.Context) (*Response, error)
func TestAPI(ctx *easierweb.Context) *Response
func TestAPI(ctx *easierweb.Context) error
func TestAPI(ctx *easierweb.Context)

// set TestAPI handle
router.EasyPOST("/test", TestAPI)
```

* Framework default use json format to process request and response data.
* If you want to change the format. You can use the plugins provided by the framework as described below. Or you can implement them by yourself.

```go
// use xml format to process request and response data
router := easierweb.New(easierweb.RouterOptions{
   RequestHandle: plugins.XMLRequestHandle(),
   ResponseHandle: plugins.XMLResponseHandle(),
})
```

* The 'EasyXXX' series functions are compatible with basic usage. You can use the 'WriteXXX' series functions to write response.

```go
func TestAPI(ctx *easierweb.Context, req Request) {

   fmt.Println("request body ->", req)
   
   // write xml response
   ctx.WriteXML(http.StatusOK, Response{Code: 1000, Msg:  "hello"})
}

// set TestAPI handle
router.EasyPOST("/test", TestAPI)
```

***

## Demo program

### If you want to learn more about how to use it, read the demo program

* demos
  * basic `basic usage demo（includes websocket, SSE, upload and download file）`
    * main.go
  * easier `easier usage demo`
    * main.go
  * restful `restful application demo`
    * main.go
  * customize `customize framework configuration demo`
    * main.go