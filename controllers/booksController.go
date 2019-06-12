package controllers

import (
	"fmt"
	"github.com/pmoorani/booksAPI/database"
	"github.com/pmoorani/booksAPI/models"
	"gopkg.in/go-playground/validator.v9"
	"net/http"

	"github.com/gin-gonic/gin"
)

var validate *validator.Validate
var books []models.Book
var users = map[string]string{
	"user1": "password1",
	"user2": "password2",
}

func GetAllBooks(c *gin.Context) {
	database.DB.Find(&books)
	c.JSON(http.StatusOK, gin.H{
		"result": books,
	})
}

func GetBook(c *gin.Context) {
	id := c.Param("id")
	database.DB.Find(&books, id)
	c.JSON(http.StatusOK, gin.H{
		"result": books,
	})
}

func CreateBook(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	c.JSON(http.StatusOK, gin.H{
		"token": authHeader,
	})
}

func UpdateBook(c *gin.Context) {
	fmt.Println(c.Param("id"))
}

func DeleteBook(c *gin.Context) {
	fmt.Println(c.Param("id"))
}
