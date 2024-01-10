package main

import (
	"fmt"
	"github.com/dpwgc/easierweb"
	"log"
	"mime/multipart"
	"net/http"
	"time"
)

// basic usage demo
// basic usage, like gin and echo
// includes common APIs, websocket connection, server-send events(SSE), file upload and download, static file service, middleware
func main() {

	// create a router, root path is /test
	router := easierweb.New(easierweb.RouterOptions{
		RootPath:               "/test",
		MultipartFormMaxMemory: 4096,
	})

	// set middleware
	router.Use(DemoMiddleware())

	// set handles
	router.GET("/demoGet/:id", DemoGet)
	router.POST("/demoPost", DemoPost)
	router.POST("/demoUpload", DemoUpload)
	router.GET("/demoDownload/:fileName", DemoDownload)

	// set websocket handle
	router.WS("/demoWS/:id", DemoWS)

	// set server-sent events handle
	router.SSE("/demoSSE/:id", DemoSSE)

	// static file service (access the demo directory)
	router.Static("/demoStatic/*filepath", "demo")

	// started on port 80
	log.Fatal(router.Run(":80"))
}

// DemoGet get request handle
// http://localhost/test/demoGet/123?type=1&price=10.24&name=dpwgc
func DemoGet(ctx *easierweb.Context) {

	// print uri path parameter
	fmt.Println("id:", ctx.Path.Int64("id"))

	// print keys for all uri query parameters
	fmt.Println("query keys:", ctx.Query.Keys())

	// print values for all uri query parameters
	fmt.Println("query values:", ctx.Query.Values())

	// bind the uri query parameters to the DemoQuery struct
	query := DemoQuery{}
	err := ctx.BindQuery(&query)
	if err != nil {
		panic(err)
	}

	// print parameters
	fmt.Println("type:", query.Type)
	fmt.Println("price:", query.Price)
	fmt.Println("name:", query.Name)

	// write response
	ctx.WriteJSON(http.StatusOK, DemoResultDTO{
		Msg:  "hello world",
		Data: "GET Request",
	})
}

// DemoPost post request handle
// http://localhost/test/demoPost
/*
{
  "id": 123,
  "name": "dpwgc"
}
*/
func DemoPost(ctx *easierweb.Context) {

	command := DemoCommand{}

	// bind the json body data to the DemoCommand struct
	err := ctx.BindJSON(&command)

	// bind the yaml body data to the DemoCommand struct
	// err := ctx.BindYAML(&command)

	// bind the xml body data to the DemoCommand struct
	// err := ctx.BindXML(&command)

	if err != nil {
		panic(err)
	}

	// print body data
	fmt.Println("body -> id:", command.Id, ", name:", command.Name)

	// write response
	ctx.WriteJSON(http.StatusOK, DemoResultDTO{
		Msg:  "hello world",
		Data: "POST Request",
	})
}

// DemoWS websocket connection handle
// ws://localhost/test/demoWS/123
func DemoWS(ctx *easierweb.Context) {

	// print the uri parameter
	fmt.Println("id:", ctx.Path.Int64("id"))

	// handles websocket connection
	for {
		// read string message
		msg, err := ctx.ReceiveString()
		// read bytes message
		// msg, err := ctx.Receive()
		if err != nil {
			panic(err)
		}

		fmt.Println("read websocket msg:", msg)

		// send json message
		err = ctx.SendJSON(DemoResultDTO{
			Msg:  "hello world",
			Data: "Websocket Connect",
		})
		if err != nil {
			panic(err)
		}

		time.Sleep(3 * time.Second)

		// when the function return, the connection is automatically closed
		return
	}
}

// DemoUpload file upload handle
// http://localhost/test/demoUpload
// form parameters: file=demo.txt(file)
func DemoUpload(ctx *easierweb.Context) {

	// get the keys for all form files
	fmt.Println("file keys:", ctx.FileKeys())

	// get the file
	file, err := ctx.GetFile("file")
	if err != nil {
		panic(err)
	}
	defer func(file multipart.File) {
		err = file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	fmt.Println("file:", file)

	// write response
	ctx.WriteJSON(http.StatusOK, DemoResultDTO{
		Msg:  "hello world",
		Data: "Upload File",
	})
}

// DemoDownload file download handle
// http://localhost/test/demoDownload/README.md
// download the specified file in the directory where the service is running
func DemoDownload(ctx *easierweb.Context) {

	// get the local file and return file data
	// if the contentType parameter is not specified, application/octet-stream is the default
	// if the fileName parameter is not specified, the file name is a timestamp by default
	ctx.WriteLocalFile(ctx.Path.Get("fileName"), ctx.Path.Get("fileName"))

	// returns the byte data of the file directly
	// ctx.WriteFile("", ctx.Path.Get("fileName"), []byte{})
}

// DemoSSE server-sent events handle
// http://localhost/test/demoSSE/123
func DemoSSE(ctx *easierweb.Context) {

	// print the uri parameter
	fmt.Println("id:", ctx.Path.Int64("id"))

	for i := 0; i < 5; i++ {

		// push string message, split the message with '\n\n'
		err := ctx.Push(fmt.Sprintf("sse push %v", i), "\n\n")
		if err != nil {
			panic(err)
		}

		// push json message, split the message with '\n\n'
		err = ctx.PushJSON(DemoResultDTO{
			Msg:  "hello world",
			Data: "Server-sent Events",
		}, "\n\n")
		if err != nil {
			panic(err)
		}

		time.Sleep(1 * time.Second)
	}
}

// DemoMiddleware middleware handle
func DemoMiddleware() easierweb.Handle {

	return func(ctx *easierweb.Context) {

		// before processing, print the url
		fmt.Println("\nrequest url:", ctx.Request.URL.String())

		// next handle
		ctx.Next()

		// after processing, print the result
		fmt.Println("result:", string(ctx.Result))
	}
}

// DemoCommand write request
// if you want to use the bind data feature, you need to configure the tag to get the field mapping
// when parsing json body data, use json tag
// when parsing uri query parameters data, use schema tag
type DemoCommand struct {
	Name string `json:"name"`
	Id   int64  `json:"id"`
}

// DemoQuery query request
// if you want to use the bind data feature, you need to configure the tag to get the field mapping
// when parsing json body data, use json tag
// when parsing query/path/form parameters data, use mapstructure tag
type DemoQuery struct {
	Type  int     `mapstructure:"type"`
	Price float64 `mapstructure:"price"`
	Name  string  `mapstructure:"name"`
}

// DemoResultDTO result data
type DemoResultDTO struct {
	Msg  string `json:"msg"`
	Data string `json:"data"`
}
