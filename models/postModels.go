package models

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	Title  string `json:"title" binding:"required"`
	Body   string `json:"body" binding:"required"`
	Author string `json:"author" binding:"required"`
	Belt   string `json:"belt" binding:"required"`
}
