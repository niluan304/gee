package main

import (
	"gee/web/day10/handle"
	"gee/web/day10/internal/controller"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/user/:id", handle.Handle(controller.User.Get))
	r.GET("/user/:id/team", handle.Handle(controller.User.GetWithTeam))

	r.GET("/team/:id", handle.Handle(controller.Team.Get))
	r.GET("/team/:id/users", handle.Handle(handle.NewReqResFunc(controller.Team.GetUsers).DecodeFunc()))

	r.Run()
}
