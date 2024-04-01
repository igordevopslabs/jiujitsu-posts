package main

import (
	"github.com/gin-gonic/gin"
	"github.com/igordevopslabs/jiujitsu-posts/config"
	"github.com/igordevopslabs/jiujitsu-posts/controller"
	docs "github.com/igordevopslabs/jiujitsu-posts/docs"
	"github.com/igordevopslabs/jiujitsu-posts/middleware"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {

	config.LoadEnvs()
	config.ConnectToDB()
	config.SyncDatabase() //migrate DB
}

// @title API Jiu-Jitsu Posts
// @version 1.0
// @description API simples para retornar os maiores jargões do jiu-jitsu através de posts objetivos.
// @host localhost:5000
// @BasePath /
// @SecurityDefinitions BearerAuth
// @in header
// @name Authorization

func main() {
	r := gin.Default()

	docs.SwaggerInfo.BasePath = "/"

	//Docs
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	//Entities Routes
	r.POST("/singup", controller.Singup)
	r.POST("/login", controller.Login)
	r.GET("/posts/:id", controller.GetPostByID)
	r.GET("/posts", controller.GetAllPosts)
	r.POST("/post", middleware.RequireAuth, controller.CreatePost)
	r.PUT("/post/:id", middleware.RequireAuth, controller.UpdatePost)
	r.DELETE("/post/:id", middleware.RequireAuth, controller.DeletePost)

	//Probe Routes
	r.GET("/healthy", config.Health)
	r.GET("/ready", config.Ready)

	r.Run()
}
