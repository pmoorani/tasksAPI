package controllers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/pmoorani/booksAPI/database"
	"github.com/pmoorani/booksAPI/models"

	"net/http"

	"gopkg.in/go-playground/validator.v9"
)

func GetAllTasks(c *gin.Context) {
	var tasks, err = models.AllTasks()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg":     err.Error(),
			"success": 0,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"tasks":   tasks,
		"success": 1,
	})
}

func GetTask(c *gin.Context) {
	fmt.Println("inside tasksController.GetTask()")
	id := c.Param("uuid")
	fmt.Println(id)
	var task, err = models.FindTaskByID(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg":     err.Error(),
			"success": 0,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"task":    task,
		"success": 1,
	})
}

func CreateTask(c *gin.Context) {
	var task models.Task
	validate := validator.New()

	claimsInterface, _ := c.Get("claims")
	claims := claimsInterface.(*models.Claims)

	// Get the JSON body and decode into Task
	err := c.ShouldBindJSON(&task)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP err
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Something went wrong!"})
		return
	}

	err = validate.Struct(task)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		c.JSON(http.StatusBadRequest, gin.H{"msg": validationErrors.Error()})
		return
	}

	if resp, ok := task.Validate(); !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": resp, "success": 0})
		return
	}

	task.UserID = claims.UserId
	if err = database.DB.Debug().Create(&task).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Some error occurred!", "success": 0})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"msg": "Task has been created!", "task": task, "success": 1})
}

func UpdateTask(c *gin.Context) {
	var task models.Task
	id := c.Param("uuid")
	validate := validator.New()

	taskFromDB, err := models.FindTaskByID(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"msg":     err.Error(),
			"success": 0,
		})
		return
	}

	// Get the JSON body and decode into Task
	err = c.ShouldBindJSON(&task)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP err
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Something went wrong!"})
		return
	}

	err = validate.Struct(task)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		c.JSON(http.StatusBadRequest, gin.H{"msg": validationErrors.Error()})
		return
	}

	if resp, ok := task.Validate(); !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": resp, "success": 0})
		return
	}

	taskFromDB.Title = task.Title
	taskFromDB.Description = task.Description
	// taskFromDB.Priority.ID = task.Priority.ID
	// taskFromDB.Status.ID = task.Status.ID
	taskFromDB.Start = task.Start
	taskFromDB.End = task.End

	if err = database.DB.Debug().Save(&taskFromDB).Error; err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Some error occurred!", "success": 0})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"task":    taskFromDB,
		"success": 1,
	})
}

func DeleteTask(c *gin.Context) {
	id := c.Param("uuid")

	task, err := models.FindTaskByID(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"msg":     err.Error(),
			"success": 0,
		})
		return
	}

	if err = database.DB.Debug().Delete(&task).Error; err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Some error occurred!", "success": 0})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{
		"msg":     "Successfully deleted!",
		"success": 1,
	})
}

func GetAllStatuses(c *gin.Context) {
	var statuses, err = models.AllStatuses()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg":     err.Error(),
			"success": 0,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"tasks":   statuses,
		"success": 1,
	})
}

func GetAllPriorities(c *gin.Context) {
	var priorities, err = models.AllPriorities()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg":     err.Error(),
			"success": 0,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"tasks":   priorities,
		"success": 1,
	})
}
