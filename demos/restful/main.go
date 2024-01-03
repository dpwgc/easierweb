package main

import (
	"fmt"
	"github.com/dpwgc/easierweb"
	"github.com/dpwgc/easierweb/demos/restful/app"
	"github.com/dpwgc/easierweb/middlewares"
	"log"
	"time"
)

var memberController = app.MemberController{}

// restful demo
// a more complete example
func main() {

	// create a router, root path is /api
	router := easierweb.New(easierweb.RouterOptions{
		RootPath: "/api",
	})

	// use router middleware
	router.Use(middlewares.Logger)

	// create a group, group path is /v2, set group middleware
	v2Group := router.Group("/v2", timeCost)
	{
		// set methods
		v2Group.EasyPOST("/member", memberController.Add)
		v2Group.EasyDELETE("/member/:id", memberController.Del)
		v2Group.EasyPUT("/member/:id", memberController.Edit)
		v2Group.EasyGET("/member/:id", memberController.Get)
		v2Group.EasyGET("/members", memberController.List)
	}

	// started on port 80
	log.Fatal(router.Run(":80"))
}

// middleware
func timeCost(ctx *easierweb.Context) {
	start := time.Now().UnixMilli()
	ctx.Next()
	end := time.Now().UnixMilli()
	fmt.Printf("%s -> time cost: %vms \n", ctx.Request.RequestURI, end-start)
}
