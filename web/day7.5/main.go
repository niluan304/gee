package main

import (
	"gee/web/day7.5/internal"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/user/:id", internal.User.Get)
	r.GET("/team/:id", internal.Team.Get)

	r.Run()
}
