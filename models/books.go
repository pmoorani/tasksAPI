package models

import (
	"github.com/jinzhu/gorm"
)

type Author struct {
	gorm.Model
	Fname string `json:"fname"`
	Lname string `json:"lname"`
	Books []Book `json:"books"`
}

type Book struct {
	gorm.Model
	Isbn   string  `json:"isbn"`
	Title  string  `json:"title"`
	AuthorID uint `json:"author_id"`
}

