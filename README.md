# EasierWeb

## A more user-friendly, highly customizable Go Web framework

### Based on httprouter

***

## Features
* Easier to obtain parameters and bind data. Can auto bind query/form/body data.
* Offers two different styles of use. Have a more concise way to write API.
* Highly customizable. Custom error capture and request/response data handling. Custom middleware.
* Easier to write websocket service and file service.
* Support TLS.

***

## Installation

```
go get github.com/dpwgc/easierweb
```

***

## Example

### Framework offers two different styles of use

#### `1`  Basic usage: like gin and echo
#### `2` Easier usage: like spring boot ( more concise way to write API handle )

***

### Basic usage

#### API handle function have only one context parameter

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
   // create a router
   router := easierweb.New()
   // set middleware handle
   router.Use(timeCost)
   // set api handle
   router.GET("/hello", hello)
   // runs on port 80
   log.Fatal(router.Run(":80"))
}

// get handle
func hello(ctx *easierweb.Context) {
   time.Sleep(1 * time.Second)
   // Write response, return 'hello'
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

#### Access the http url

> `GET` http://localhost/hello

#### Response data

```
hello
```

* You can use the bind function to obtain the request data.

```go
// struct
request := Request{}

// bind uri query parameters
ctx.BindQuery(&request)

// bind json body data
ctx.BindJSON(&request)
```

* Get the parameters individually.

```go
// obtain the uri path parameter
id := ctx.Path.Int64("id")

// obtain the uri query parameter
name := ctx.Path.Get("name")

// obtain the post form parameter
mobile := ctx.Form.Get("mobile")
```

***

### Easier usage

#### API handle function has input object and return values

#### More concise way to write API handle, don't need to write logic for binding data and writing response data. framework will help you do this

```go
package main

import (
   "fmt"
   "github.com/dpwgc/easierweb"
   "log"
)

// easier usage example
func main() {
   // create a router
   router := easierweb.New()
   // set api handle (use function 'EasyPOST')
   router.EasyPOST("/submit", submit)
   // runs on port 80
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

#### Invoke http api

> `POST` http://localhost/submit

#### Request body

```json
{
  "name": "hello",
  "mobile": "12345678"
}
```

#### Response data

```json
{
  "code": 1000,
  "msg": "hello"
}
```

* If you want to use function 'EasyGET', 'EasyPOST', 'EasyPUT'... the api handle function must be in the following formats.

```go
// input: Request
// output: *Response, error
func TestAPI(ctx *easierweb.Context, req Request) (*Response, error)

// input: Request
// output: *Response
func TestAPI(ctx *easierweb.Context, req Request) *Response

// input: Request
// output: empty
func TestAPI(ctx *easierweb.Context, req Request)

// input: empty
// output: *Response, error
func TestAPI(ctx *easierweb.Context) (*Response, error)

// input: empty
// output: *Response
func TestAPI(ctx *easierweb.Context) *Response

// input: empty
// output: empty
func TestAPI(ctx *easierweb.Context)

// ----- ----- -----

// set TestAPI handle
router.EasyPOST("/test", TestAPI)
```

* Framework default use json format to process request and response data.
* If you want to change the format, you can use the plugin, framework comes with multiple plugins.
* When creating a router, use 'RouterOptions' to set up the plugins.

```go
// use xml format to process request and response data (global configuration, takes effect for all api)
router := easierweb.New(easierweb.RouterOptions{
   RequestHandle: plugins.XMLRequestHandle,
   ResponseHandle: plugins.XMLResponseHandle,
})
```

* If you want to change the request and response format for a single api.
* Use 'PluginOptions' to set up the plugins.

```go
// use xml format to process request and response data (takes effect only for this api)
router.EasyGET("/test", test, easierweb.PluginOptions{
   RequestHandle: plugins.XMLRequestHandle,
   ResponseHandle: plugins.XMLResponseHandle,
})
```

***

## Demo program

### If you want to learn more about how to use it, read the demo program

* demo
  * basic `basic usage demo`
    * main.go
  * easier `easier usage demo`
    * main.go
  * restful `restful application demo`
    * main.go