# EasierWeb

## A minimalist Go web framework based on httprouter

***

### Features
* Easier to handle http request and response.
* Custom middleware framework.
* Easier to obtain path/query/form parameters and convert their type.
* Easier to bind json/yaml/xml body data.
* Easier to write websocket service.
* Easier to write file services.
* Centralized error capture.
* Highly customizable.
* Support TLS.

***

### Installation

```
go get github.com/dpwgc/easierweb
```

***

### Simple example

```go
package main

import (
   "fmt"
   "github.com/dpwgc/easierweb"
   "log"
   "net/http"
   "time"
)

// Simple example server
func main() {
   // Build router and start server on port 8080
   log.Fatal(easierweb.New().Use(timeCost).GET("/", hello).Run(":8080"))
}

// Middleware method
func timeCost(ctx *easierweb.Context) {
   start := time.Now().UnixMilli()
   // Next method
   ctx.Next()
   end := time.Now().UnixMilli()
   fmt.Printf("time cost: %vms\n", end-start)
}

// Handler method
func hello(ctx *easierweb.Context) {
   time.Sleep(1 * time.Second)
   // Write response
   ctx.WriteString(http.StatusOK, "hello")
}
```

***

### Demo program

* demo
  * main.go