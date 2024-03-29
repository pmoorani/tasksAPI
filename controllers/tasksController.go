package controllers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/pmoorani/tasksAPI/database"
	"github.com/pmoorani/tasksAPI/models"

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
	id := c.Param("uuid")
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
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "Something went wrong!", "error": err.Error()})
		return
	}

	err = validate.Struct(task)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": validationErrors.Error()})
		return
	}

	if resp, ok := task.Validate(); !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": resp, "success": 0})
		return
	}

	task.UserID = claims.UserId
	if err = database.DB.Debug().Create(&task).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "Some error occurred!", "success": 0})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"msg": "Task has been created!", "task": task, "success": 1})
}

func UpdateTask(c *gin.Context) {
	var task models.Task

	id := c.Param("uuid")
	validate := validator.New()

	transformedTask, err := models.FindTaskByID(id)
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
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": validationErrors.Error()})
		return
	}

	if resp, ok := task.Validate(); !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": resp, "success": 0})
		return
	}
	fmt.Println("task ===", &task)

	_task := models.Task{
		BaseModel:   transformedTask.BaseModel,
		Title:       transformedTask.Title,
		Description: transformedTask.Description,
		PriorityID:  transformedTask.Priority.ID,
		StatusID:    transformedTask.Status.ID,
		Start:       transformedTask.Start,
		End:         transformedTask.End,
		UserID:      transformedTask.User.ID,
	}

	if err = database.DB.Debug().Model(&_task).Update(&task).Error; err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "Some error occurred!", "success": 0})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"task":    _task,
		"success": 1,
	})
}

func DeleteTask(c *gin.Context) {
	id := c.Param("uuid")
	err := models.DeleteTaskByID(id)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"msg":     err.Error(),
			"success": 0,
		})
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
