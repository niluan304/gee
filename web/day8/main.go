package main

import (
	"gee/web/day8/internal/controller"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/user/:id", controller.User.Get)
	r.GET("/user/:id/team", controller.User.GetWithTeam)
	r.GET("/team/:id", controller.Team.Get)

	r.Run()
}
