package controllers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/pmoorani/booksAPI/database"
	"github.com/pmoorani/booksAPI/models"

	"net/http"

	"gopkg.in/go-playground/validator.v9"
)

func GetAllUsers(c *gin.Context) {
	var users, err = models.AllUsers()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg":     err.Error(),
			"success": 0,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"users":   users,
		"success": 1,
	})
}

func GetUser(c *gin.Context) {
	fmt.Println("inside tasksController.GetTask()")
	id := c.Param("uuid")
	fmt.Println(id)
	var user, err = models.FindUserByID(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg":     err.Error(),
			"success": 0,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user":    user,
		"success": 1,
	})
}

func UpdateUser(c *gin.Context) {
	var user models.User
	id := c.Param("uuid")
	validate := validator.New()

	userFromDB, err := models.FindUserByID(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"msg":     err.Error(),
			"success": 0,
		})
		return
	}

	// Get the JSON body and decode into Task
	err = c.ShouldBindJSON(&user)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP err
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Something went wrong!"})
		return
	}

	err = validate.Struct(user)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		c.JSON(http.StatusBadRequest, gin.H{"msg": validationErrors.Error()})
		return
	}

	if resp, ok := user.Validate(); !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": resp, "success": 0})
		return
	}

	userFromDB.FirstName = user.FirstName
	userFromDB.LastName = user.LastName

	if err = database.DB.Debug().Save(&userFromDB).Error; err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Some error occurred!", "success": 0})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"task":    userFromDB,
		"success": 1,
	})
}

func DeleteUser(c *gin.Context) {
	id := c.Param("uuid")

	user, err := models.FindUserByID(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"msg":     err.Error(),
			"success": 0,
		})
		return
	}

	if err = database.DB.Debug().Delete(&user).Error; err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Some error occurred!", "success": 0})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{
		"msg":     "Successfully deleted!",
		"success": 1,
	})
}
