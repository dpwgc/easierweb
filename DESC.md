# Function Description

***

## easierweb.Router

### Create

```go
// create a router
router := easierweb.New()
// create a router (with RouterOptions)
router := easierweb.New(easierweb.RouterOptions{
   RootPath: "/api",
})
```

### Set Middlewares

```go
// set middlewares
router.Use(middlewares.Logger())
```

### Set APIs Handle

```go
// APIs (basic usage)
router.GET("/hello", hello)
router.HEAD("/hello", hello)
router.OPTIONS("/hello", hello)
router.POST("/hello", hello)
router.PUT("/hello", hello)
router.PATCH("/hello", hello)
router.DELETE("/hello", hello)
router.Any("/hello", hello)
router.Handle("GET", "/hello", hello)

// APIs (easier usage)
router.EasyGET("/hello", hello)
router.EasyHEAD("/hello", hello)
router.EasyOPTIONS("/hello", hello)
router.EasyPOST("/hello", hello)
router.EasyPUT("/hello", hello)
router.EasyPATCH("/hello", hello)
router.EasyDELETE("/hello", hello)
router.EasyAny("/hello", hello)
router.EasyHandle("GET", "/hello", hello)
```

### Set Other Handle

```go
// websocket connect
router.WS("/hello", hello)
// server-sent events (SSE)
router.SSE("/hello", hello)
// static file server
router.Static("/hello", "demo")
router.StaticFS("/hello", http.Dir("demo"))
```

### Start And Close

```go
// start server
router.Run(":80")
// start server with TLS & HTTP2
router.RunTLS("0.0.0.0:443", "cert.pem", "private.key", &tls.Config{})
// custom HTTP server and start server
router.Serve(&http.Server{})
router.ServeTLS(&http.Server{}, "cert.pem", "private.key")
// close server
router.Close()
```

***

## easierweb.Group

### Create

```go
// create a api group
group := router.Group("/group")
// create a api group and set middleware
group := router.Group("/group", middlewares.Logger())
```

### Set APIs Handle

```go
// set handle (setting method is the same as router)
group.GET("/hello", hello)
group.EasyGET("/hello", hello)
```

***

## easierweb.Context

### Middleware Related Operations

```go
// go to the next handle
ctx.Next()
// process termination
ctx.Abort()
```

### Bind Request Data

```go
// bind uri query parameters (based on mapstructure)
ctx.BindQuery(&request)
// bind uri path parameters (based on mapstructure)
ctx.BindPath(&request)
// bind form parameters (based on mapstructure)
ctx.BindForm(&request)
// bind header parameters (based on mapstructure)
ctx.BindHeader(&request)
// bind body data
ctx.BindJSON(&request)
ctx.BindYAML(&request)
ctx.BindXML(&request)
```

### Write Response

```go
// write response
ctx.WriteJSON(http.StatusOK, Response{Msg:  "hello world"})
ctx.WriteYAML(http.StatusOK, Response{Msg:  "hello world"})
ctx.WriteXML(http.StatusOK, Response{Msg:  "hello world"})
ctx.WriteLocalFile("application/octet-stream", "name", "demo")
ctx.WriteFile("application/octet-stream", "name", []byte("hello world"))
ctx.WriteString(http.StatusOK, "hello world")
ctx.NoContent(http.StatusNoContent)
ctx.Write(http.StatusOK, []byte("hello world"))
ctx.Redirect(http.StatusOK, "http://127.0.0.1/hello")
```

### Set Response Header

```go
// set response header
ctx.SetContentType("application/json")
ctx.SetContentDisposition(fmt.Sprintf("attachment; filename=\"%v\"", time.Now().Unix()))
ctx.SetHeader("Content-Type", "application/json")
ctx.AddHeader("Content-Type", "application/json")
```

### Websocket Connect

```go
// prerequisites for using these function: router.WS("/hello", hello)

// receive websocket message
ctx.ReceiveJSON(&message)
ctx.ReceiveYAML(&message)
ctx.ReceiveXML(&message)
ctx.ReceiveString()
ctx.Receive()

// send websocket message
ctx.SendJSON(Message{Msg:  "hello world"})
ctx.SendYAML(Message{Msg:  "hello world"})
ctx.SendXML(Message{Msg:  "hello world"})
ctx.SendString("hello world")
ctx.Send([]byte("hello world"))

// close websocket connect
ctx.Close()
```

### Server-Sent Events (SSE)

```go
// prerequisites for using these function: router.SSE("/hello", hello)

// server-sent events (SSE) push message
ctx.PushJSON(Message{Msg:  "hello world"}, "\n\n")
ctx.PushYAML(Message{Msg:  "hello world"}, "\n\n")
ctx.PushXML(Message{Msg:  "hello world"}, "\n\n")
ctx.Push("hello world", "\n\n")
```

### File

```go
// get all the file keys in the form
ctx.FileKeys()
// get form file by key
ctx.GetFile("hello")
```

### Other Request Parameters

```go
// get all cookies
ctx.Cookies()
// get cookie by name
ctx.GetCookie("hello")
// get other request parameters
ctx.URI()
ctx.Method()
ctx.URL()
ctx.RemoteAddr()
ctx.Host()
ctx.Proto()
```

### Logger

```go
ctx.Info("hello")
ctx.Debug("hello")
ctx.Warn("hello")
ctx.Error(errors.New("hello"))
// log error and panic error
ctx.Panic(errors.New("hello"))
```

***

## easierweb.Params 

### `ctx.Path` `ctx.Query` `ctx.Form` `ctx.Header`

### Base

```go
// get all keys
ctx.Query.Keys()
// get all values
ctx.Query.Values()
// get value (string) by key
ctx.Query.Get("hello")
// check whether the key exists
ctx.Query.Has("hello")
// set key-value
ctx.Query.Set("hello", "world")
// delete key
ctx.Query.Del("hello")
// bind struct
ctx.Query.Bind(&request)
```

### Value Type Conversion

```go
// when parsing errors, error will be returned
ctx.Query.ParseInt("hello")
ctx.Query.ParseInt32("hello")
ctx.Query.ParseInt64("hello")
ctx.Query.ParseFloat32("hello")
ctx.Query.ParseFloat64("hello")

// when parsing errors, default value will be returned
ctx.Query.GetInt("hello")
ctx.Query.GetInt32("hello")
ctx.Query.GetInt64("hello")
ctx.Query.GetFloat32("hello")
ctx.Query.GetFloat64("hello")

// when parsing errors, panic error
ctx.Query.Int("hello")
ctx.Query.Int32("hello")
ctx.Query.Int64("hello")
ctx.Query.Float32("hello")
ctx.Query.Float64("hello")
```

***

## easierweb.Data

### `ctx.Body` `ctx.Result`

### Parse

```go
// parse response data
ctx.Result.ParseJSON(&response)
ctx.Result.ParseYAML(&response)
ctx.Result.ParseXML(&response)
```

### Save

```go
// overwrite request body data
ctx.Body.SaveJSON(request)
ctx.Body.SaveYAML(request)
ctx.Body.SaveXML(request)
ctx.Body.Save([]byte("hello"))
```