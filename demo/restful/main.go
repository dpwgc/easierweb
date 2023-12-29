package main

import (
	"fmt"
	"github.com/dpwgc/easierweb"
	"github.com/dpwgc/easierweb/demo/restful/app"
	"log"
	"time"
)

// RESTful API Demo
func main() {

	// create a route, root path is /api/v2
	router := easierweb.New().SetContextPath("/api/v2")

	// use middleware
	router.Use(timeCost)

	// set methods
	router.EasyPOST("/member", app.AddMember)
	router.EasyDELETE("/member/:id", app.DelMember)
	router.EasyPUT("/member/:id", app.EditMember)
	router.EasyGET("/member/:id", app.GetMember)
	router.EasyGET("/members", app.ListMember)

	// started on port 80
	log.Fatal(router.Run(":80"))
}

// middleware
func timeCost(ctx *easierweb.Context) {
	start := time.Now().UnixMilli()
	ctx.Next()
	end := time.Now().UnixMilli()
	fmt.Printf("[%s] time cost: %vms \n", ctx.Request.RequestURI, end-start)
}
