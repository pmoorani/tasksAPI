package models

import (
	"time"

	"github.com/pmoorani/booksAPI/database"
	u "github.com/pmoorani/booksAPI/utils"
	uuid "github.com/satori/go.uuid"
)

type (
	Task struct {
		BaseModel
		Title       string    `json:"title"`
		Description string    `json:"description"`
		Priority    Priority  `json:"priority" gorm:"association_autocreate:false;default:1"`
		Status      Status    `json:"status" gorm:"association_autocreate:false;default:1"`
		StartDate   time.Time `json:"start_date"`
		EndDate     time.Time `json:"end_date"`
		UserID      uuid.UUID `json:"user_id"`
		//Assignee	User		`json:"assignee"`
	}

	Status struct {
		BaseGormModel
		Status   string `json:"status"`
		StatusDE string `json:"status_de"`
	}

	Priority struct {
		BaseGormModel
		Priority   string `json:"priority"`
		PriorityDE string `json:"priority_de"`
	}
)

func (p *Priority) Scan(v interface{}) error {
	var priority Priority
	err := database.DB.Where("id = ?", v).Find(&priority).Error

	if err != nil {
		return err
	}

	*p = priority
	return nil
}

func (s *Status) Scan(v interface{}) error {
	var status Status
	err := database.DB.Where("id = ?", v).Find(&status).Error

	if err != nil {
		return err
	}

	*s = status
	return nil
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

	if err != nil {
		return Task{}, err
	}

	//err = database.DB.Select("id, username, first_name, last_name, email").Where("id = ?", &task.UserID).Find(&task.Assignee).Error
	return task, nil
}

func AllStatuses() ([]Status, error) {
	var statuses []Status
	err := database.DB.Find(&statuses).Error

	if err != nil {
		return nil, err
	}
	return statuses, nil
}

func AllPriorities() ([]Priority, error) {
	var priorities []Priority
	err := database.DB.Find(&priorities).Error

	if err != nil {
		return nil, err
	}
	return priorities, nil
}
