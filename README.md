# EasierWeb

## A minimalist Go web framework based on httprouter

***

## Features
* Easier to handle http request and response.
* Custom middleware framework.
* Easier to obtain parameters and bind data.
* Can auto bind query/form/body data.
* Easier to write websocket service.
* Easier to write file services.
* Centralized error capture.
* Highly customizable.
* Support TLS.

***

## Installation

```
go get github.com/dpwgc/easierweb
```

***

## Example

framework provides two different styles of API code writing

* basic
* easier

### Basic API code example

* like gin and echo, only one context parameter

```go
package main

import (
   "fmt"
   "github.com/dpwgc/easierweb"
   "log"
   "net/http"
   "time"
)

// basic example
func main() {
   // create a router and started on port 80
   log.Fatal(easierweb.New().Use(timeCost).GET("/", hello).Run(":80"))
}

// middleware handle
func timeCost(ctx *easierweb.Context) {
   start := time.Now().UnixMilli()
   // next handle
   ctx.Next()
   end := time.Now().UnixMilli()
   fmt.Printf("time cost: %vms\n", end-start)
}

// get handle
func hello(ctx *easierweb.Context) {
   time.Sleep(1 * time.Second)
   // Write response
   ctx.WriteString(http.StatusOK, "hello")
}
```

* you can use the bind method to obtain the request data

```go
// struct
request := Request{}

// bind uri query parameters
ctx.BindQuery(&request)

// bind json body data
ctx.BindJSON(&request)
```

```go
// obtain the uri path parameter
id := ctx.Path.Int64("id")

// obtain the uri query parameter
name := ctx.Path.Get("name")

// obtain the post form parameter
mobile := ctx.Form.Get("mobile")
```

### Easier API code example

* like spring boot, has input object and return values
* easier to write api code, don't need to write logic for binding data and writing response data. framework will help you do this

```go
package main

import (
   "fmt"
   "github.com/dpwgc/easierweb"
   "log"
)

// easier example
func main() {
   // create a router and set a handle
   router := easierweb.New().EasyPOST("/", hello)
   // started on port 80
   log.Fatal(router.Run(":80"))
}

// post request handle
func hello(ctx *easierweb.Context, request Request) *Response {
   // print the request data
   fmt.Printf("post request data (json body) -> name: %s, mobile: %s \n", request.Name, request.Mobile)
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

* framework default use json format to process request and response data
* if you want to change the format, you can use the plugin, framework comes with multiple plug-ins
* use method 'SetEasyHandlePlugins' to set up the plug-ins

```go
// use xml format to process request and response data (global configuration, takes effect for all api)
router.SetEasyHandlePlugins(plugins.XMLRequestHandle, plugins.XMLResponseHandle)
```

* if you want to change the request and response format for a single api
* use method 'ReEasyGET' to set up the path, handle and plug-ins

```go
// use xml format to process request and response data (takes effect only for this api)
router.ReEasyGET("/test", TestHandle, plugins.XMLRequestHandle, plugins.XMLResponseHandle)
```

***

### Demo program

* demo
  * base `basic usage demo`
    * main.go
  * easier `easier usage demo`
    * main.go
  * restful `restful application demo`
    * app
    * main.go