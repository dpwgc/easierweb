package main

import (
	"fmt"
	"github.com/dpwgc/easierweb"
	"github.com/dpwgc/easierweb/demo/restful/app"
	"log"
	"time"
)

// restful demo
// a more complete example
func main() {

	// create a router, root path is /api/v2
	router := easierweb.New(easierweb.RouterOptions{
		RootPath: "/api/v2",
	})

	// use middleware
	router.Use(timeCost)

	memberController := app.MemberController{}

	// set methods
	router.EasyPOST("/member", memberController.Add)
	router.EasyDELETE("/member/:id", memberController.Del)
	router.EasyPUT("/member/:id", memberController.Edit)
	router.EasyGET("/member/:id", memberController.Get)
	router.EasyGET("/members", memberController.List)

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
