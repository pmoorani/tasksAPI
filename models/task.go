package models

import (
	"github.com/pmoorani/booksAPI/database"
	u "github.com/pmoorani/booksAPI/utils"
)

type Task struct {
	BaseModel
	Title string `json:"title"`
	Completed bool `json:"completed"`
}

func (task *Task) Validate() (map[string] interface{}, bool) {
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

func FindTaskByID(id interface{}) ([]Task, error) {
	err := database.DB.Where("id = ?", id).Find(&tasks).Error
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

