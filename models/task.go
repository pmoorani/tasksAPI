package models

import (
	"fmt"

	"github.com/pmoorani/booksAPI/database"
	u "github.com/pmoorani/booksAPI/utils"
	uuid "github.com/satori/go.uuid"
)

type Task struct {
	BaseModel
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	Status    uint      `json:"status" gorm:"default:1"`
	UserID    uuid.UUID `json:"user_id"`
}

func (task *Task) Validate() (map[string]interface{}, bool) {
	return u.Message(true, "Requirement passed!"), true
}

var tasks []Task

func AllTasks() ([]Task, error) {
	err := database.DB.Find(&tasks).Error
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func FindTaskByID(id interface{}) (Task, error) {
	var task Task
	err := database.DB.Where("id = ?", id).Find(&task).Error
	fmt.Println(err)

	if err != nil {
		return Task{}, err
	}
	return task, nil
}
