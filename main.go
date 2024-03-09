package main

import (
	"github.com/gin-gonic/gin"
	"github.com/igordevopslabs/jiujitsu-posts/config"
	"github.com/igordevopslabs/jiujitsu-posts/controller"
)

func init() {

	config.LoadEnvs()
	config.ConnectToDB()
}

func main() {
	r := gin.Default()

	r.POST("/post", controller.CreatePost)
	r.GET("/posts", controller.GetAllPosts)
	r.GET("/posts/:id", controller.GetPostByID)
	r.PUT("/post/:id", controller.UpdatePost)
	r.DELETE("/post/:id", controller.DeletePost)

	r.Run()
}
