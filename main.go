package main

import (
	"github.com/gin-gonic/gin"
	"github.com/igordevopslabs/jiujitsu-posts/config"
	"github.com/igordevopslabs/jiujitsu-posts/controller"
	docs "github.com/igordevopslabs/jiujitsu-posts/docs"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {

	config.LoadEnvs()
	config.ConnectToDB()
}

// @title API Jiu-Jitsu Posts
// @version 1.0
// @description API simples para retornar os maiores jargões do jiu-jitsu através de posts objetivos.
// @host localhost:5000
// @BasePath /

func main() {
	r := gin.Default()

	docs.SwaggerInfo.BasePath = "/"

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.POST("/post", controller.CreatePost)
	r.GET("/posts", controller.GetAllPosts)
	r.GET("/posts/:id", controller.GetPostByID)
	r.PUT("/post/:id", controller.UpdatePost)
	r.DELETE("/post/:id", controller.DeletePost)

	r.Run()
}
