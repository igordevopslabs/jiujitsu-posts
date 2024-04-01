package config

import (
	"github.com/igordevopslabs/jiujitsu-posts/models"
)

func SyncDatabase() {

	DB.AutoMigrate(&models.Post{})
	DB.AutoMigrate(&models.User{})

}
