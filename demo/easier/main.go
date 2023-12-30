package main

import (
	"fmt"
	"github.com/dpwgc/easierweb"
	"log"
)

// easier demo
// easier to write api code
// don't need to write logic for binding data and writing response data. the framework will help you do this
func main() {

	// create a router
	router := easierweb.New()

	// the framework default use json format to process request and response data
	// if you want to change the format, you can use the plugin, framework comes with multiple plug-ins
	// example: use xml format to process request and response data
	// router.SetEasyHandlePlugins(plugins.XMLRequestHandle, plugins.XMLResponseHandle)

	// set handles
	router.EasyPOST("/submit", Submit)
	router.EasyGET("/getById/:id", GetById)
	router.EasyGET("/searchByParams", SearchByParams)

	// started on port 80
	log.Fatal(router.Run(":80"))
}

// Submit submit data and do not return result
// [POST] http://localhost/submit
/*
body:
{
  "name": "hello",
  "mobile": "12345678"
}
*/
func Submit(ctx *easierweb.Context, request Request) {

	// print the request data
	fmt.Printf("post request data (json body) -> name: %s, mobile: %s \n", request.Name, request.Mobile)

	// no return value
}

// GetById get a piece of data based on the id
// [GET] http://localhost/getById/1
func GetById(ctx *easierweb.Context) *Response {

	// print the path parameter
	fmt.Printf("path parameter -> id: %v", ctx.Path.Int64("id"))

	// return result
	return &Response{
		ID:     ctx.Path.Int64("id"),
		Name:   "hello",
		Mobile: "12345678",
	}
}

// SearchByParams search multiple data based on multiple parameters
// [GET] http://localhost/searchByParams?name=hello&mobile=12345678
func SearchByParams(ctx *easierweb.Context, request Request) *[]Response {

	// print the get request data
	fmt.Printf("get request data (uri query parameters) -> name: %s, mobile: %s \n", request.Name, request.Mobile)

	// build a list of results
	var list []Response
	list = append(list, Response{
		ID:     1,
		Name:   "hello",
		Mobile: "12345678",
	}, Response{
		ID:     2,
		Name:   "world",
		Mobile: "87654321",
	})

	// return the list of results
	return &list
}

// Request if you want to use the bind data feature, you need to configure the tag to get the field mapping
// when parsing json body data, use json tag
// when parsing uri query parameters data, use schema tag
type Request struct {
	Name   string `json:"name" schema:"name"`
	Mobile string `json:"mobile" schema:"mobile"`
}

// Response result data
type Response struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Mobile string `json:"mobile"`
}
