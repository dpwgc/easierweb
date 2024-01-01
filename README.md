# EasierWeb

## A more user-friendly, high performance, high customizable Go Web framework

### Based on [httprouter](https://github.com/julienschmidt/httprouter)

***

## Features
* Easier to obtain parameters and bind data. Can auto bind query/form/body data.
* Have a more concise way to write API. Easier to handle request and response
* Highly customizable. Custom error capture and request/response data handling. Custom middleware.
* No dependencies on too many third-party packages. Architecture is simple
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

### `1` Basic usage: like gin and echo
### `2` Easier usage: like spring boot ( more concise way to write API handle )

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
   // create a router
   router := easierweb.New()
   // set middleware handle
   router.Use(middlewares.Logger)
   // set api handle
   router.GET("/hello", hello)
   // runs on port 80
   log.Fatal(router.Run(":80"))
}

// get method request handle
func hello(ctx *easierweb.Context) {
   // Write response, return 'hello'
   ctx.WriteString(http.StatusOK, "hello")
}
```

### Access the HTTP URL in your browser

> `GET` http://localhost/hello

### Other notes

* You can use context bind function to obtain the request data.

```go
// struct
request := Request{}

// bind uri query parameters (based on mapstructure)
ctx.BindQuery(&request)

// bind uri path parameters (based on mapstructure)
ctx.BindPath(&request)

// bind post form parameters (based on mapstructure)
ctx.BindForm(&request)

// bind json body data
ctx.BindJSON(&request)
```

* Take a single parameter and convert its type.

```go
// obtain the uri path parameter
id := ctx.Path.Int64("id")

// obtain the uri query parameter
name := ctx.Query.Float32("weight")

// obtain the post form parameter
mobile := ctx.Form.Get("mobile")
```

***

## Easier usage

### API handle function has input object and return values

### Don't need to manually bind request data and write response data. framework will help you do this

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
   // set api handle (use function 'EasyXXX')
   router.EasyPOST("/submit", submit)
   // runs on port 80
   log.Fatal(router.Run(":80"))
}

// post method request handle
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

* demos
  * basic `basic usage demo`
    * main.go
  * easier `easier usage demo`
    * main.go
  * restful `restful application demo`
    * main.go