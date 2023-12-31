# EasierWeb

## A minimalist Go web framework based on httprouter

***

## Features
* Easier to handle http request and response.
* Custom middleware framework.
* Easier to obtain parameters and bind data. Can auto bind query/form/body data.
* Easier to write websocket service and file services.
* Highly customizable. Custom error capture and request/response data handling.
* Offers two different styles of use.
* Support TLS.

***

## Installation

```
go get github.com/dpwgc/easierweb
```

***

## Example

### Framework offers two different styles of use

#### `Basic usage`: like gin and echo
#### `Easier usage`: like spring boot

### Basic usage

* api handle function have only one context parameter

```go
package main

import (
   "fmt"
   "github.com/dpwgc/easierweb"
   "log"
   "net/http"
   "time"
)

// basic usage example
func main() {
   // create a router and started on port 80
   log.Fatal(easierweb.New().Use(timeCost).GET("/hello", hello).Run(":80"))
}

// get handle
func hello(ctx *easierweb.Context) {
   time.Sleep(1 * time.Second)
   // Write response
   ctx.WriteString(http.StatusOK, "hello")
}

// middleware handle
func timeCost(ctx *easierweb.Context) {
  start := time.Now().UnixMilli()
  // next handle
  ctx.Next()
  end := time.Now().UnixMilli()
  fmt.Printf("time cost: %vms\n", end-start)
}
```

* access the http url

> `GET` http://localhost/hello

* you can use the bind function to obtain the request data

```go
// struct
request := Request{}

// bind uri query parameters
ctx.BindQuery(&request)

// bind json body data
ctx.BindJSON(&request)
```

* get the parameters individually

```go
// obtain the uri path parameter
id := ctx.Path.Int64("id")

// obtain the uri query parameter
name := ctx.Path.Get("name")

// obtain the post form parameter
mobile := ctx.Form.Get("mobile")
```

### Easier usage

* api handle function has input object and return values
* easier to write api code, don't need to write logic for binding data and writing response data. framework will help you do this

```go
package main

import (
   "fmt"
   "github.com/dpwgc/easierweb"
   "log"
)

// easier usage example
func main() {
   // create a router and set a handle (use function EasyPOST)
   router := easierweb.New().EasyPOST("/submit", submit)
   // started on port 80
   log.Fatal(router.Run(":80"))
}

// post request handle
func submit(ctx *easierweb.Context, req Request) *Response {
   // print the request data
   fmt.Printf("post request data (json body) -> name: %s, mobile: %s \n", req.Name, req.Mobile)
   // return result
   return &Response{
      Code: 1000,
      Msg:  "hello",
   }
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

* request the http post api

> `POST` http://localhost/submit

* request body

```json
{
  "name": "hello",
  "mobile": "12345678"
}
```

* framework default use json format to process request and response data
* if you want to change the format, you can use the plugin, framework comes with multiple plugins
* when creating a router, use 'RouterOptions' to set up the plugins

```go
// use xml format to process request and response data (global configuration, takes effect for all api)
router := easierweb.New(easierweb.RouterOptions{
   RequestHandle: plugins.XMLRequestHandle,
   ResponseHandle: plugins.XMLResponseHandle,
})
```

* if you want to change the request and response format for a single api
* use 'PluginOptions' to set up the plugins

```go
// use xml format to process request and response data (takes effect only for this api)
router.EasyGET("/test", test, easierweb.PluginOptions{
   RequestHandle: plugins.XMLRequestHandle,
   ResponseHandle: plugins.XMLResponseHandle,
})
```

***

## Demo program

* demo
  * base `basic usage demo`
    * main.go
  * easier `easier usage demo`
    * main.go
  * restful `restful application demo`
    * app
    * main.go