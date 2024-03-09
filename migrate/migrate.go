package main

import (
	"github.com/igordevopslabs/jiujitsu-posts/config"
	"github.com/igordevopslabs/jiujitsu-posts/models"
)

func init() {
	config.LoadEnvs()
	config.ConnectToDB()
}

func main() {

	config.DB.AutoMigrate(&models.Post{})

}
