package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/igordevopslabs/jiujitsu-posts/config"
	"github.com/igordevopslabs/jiujitsu-posts/models"
	"go.uber.org/zap"
)

func CreatePost(c *gin.Context) {

	config.LogInfo("Init journey to create a new post", zap.String("journey", "createPost"))
	//pegar os dados do corpo da requisição
	var body models.Post

	c.ShouldBindJSON(&body)

	//criar o post
	post := models.Post{Title: body.Title, Body: body.Body, Belt: body.Belt, Author: body.Author}

	result := config.DB.Create(&post)

	if result.Error != nil {
		c.Status(401)
		return
	}

	//exibir dados
	c.JSON(200, gin.H{
		"post": post,
	})
	config.LogInfo("new post created", zap.String("journey", "createPost"))
}

func GetAllPosts(c *gin.Context) {
	config.LogInfo("Init journey to get all post", zap.String("journey", "getAllPosts"))

	posts := []models.Post{} //para trazer todos os posts, precisar ser um array []

	result := config.DB.Find(&posts)

	if result.Error != nil {
		c.Status(401)
		return
	}

	//exibir dados
	c.JSON(200, gin.H{
		"post": posts,
	})
	config.LogInfo("retrieve all posts", zap.String("journey", "getAllPosts"))

}

func GetPostByID(c *gin.Context) {
	config.LogInfo("Init journey to get single post", zap.String("journey", "getPostByID"))

	id := c.Param("id")
	post := models.Post{}

	result := config.DB.First(&post, id)

	if result.Error != nil {
		c.Status(401)
		return
	}

	//exibir dados
	c.JSON(200, gin.H{
		"post": post,
	})
	config.LogInfo("retrieve a post", zap.String("journey", "getPostByID"))

}

func UpdatePost(c *gin.Context) {

	config.LogInfo("Init journey to update a post", zap.String("journey", "updatePost"))
	//pegar os dados do corpo da requisição

	id := c.Param("id")

	var body models.Post

	c.ShouldBindJSON(&body)

	var post models.Post

	config.DB.First(&post, id)

	result := config.DB.Model(&post).Updates(models.Post{
		Title:  body.Title,
		Body:   body.Body,
		Author: body.Author,
		Belt:   body.Belt,
	})

	if result.Error != nil {
		c.Status(401)
		return
	}

	//exibir dados
	c.JSON(200, gin.H{
		"post": post,
	})
	config.LogInfo("post updated", zap.String("journey", "updatePost"))
}

func DeletePost(c *gin.Context) {
	config.LogInfo("Init journey to delete single post", zap.String("journey", "deletePost"))

	id := c.Param("id")
	post := models.Post{}

	result := config.DB.Delete(&post, id)

	if result.Error != nil {
		c.Status(401)
		return
	}

	//exibir dados
	c.JSON(200, gin.H{
		"post": post,
	})
	config.LogInfo("post deleted", zap.String("journey", "deletePost"))

}
